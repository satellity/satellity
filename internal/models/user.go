package models

import (
	"context"
	"crypto/ed25519"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"satellity/internal/clouds"
	"satellity/internal/configs"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// Constants for user
const (
	UserRoleAdmin  = "admin"
	UserRoleMember = "member"
)

// User contains info of a register user
type User struct {
	UserID            string
	PublicKey         sql.NullString
	Email             sql.NullString
	Username          sql.NullString
	Nickname          string
	AvatarURL         string
	Biography         string
	EncryptedPassword sql.NullString
	GithubID          sql.NullString
	Role              string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	SessionID string
	isNew     bool
}

var userColumns = []string{"user_id", "public_key", "email", "username", "nickname", "avatar_url", "biography", "encrypted_password", "github_id", "role", "created_at", "updated_at"}

func (u *User) values() []interface{} {
	return []interface{}{u.UserID, u.PublicKey, u.Email, u.Username, u.Nickname, u.AvatarURL, u.Biography, u.EncryptedPassword, u.GithubID, u.Role, u.CreatedAt, u.UpdatedAt}
}

func userFromRow(row durable.Row) (*User, error) {
	var u User
	err := row.Scan(&u.UserID, &u.PublicKey, &u.Email, &u.Username, &u.Nickname, &u.AvatarURL, &u.Biography, &u.EncryptedPassword, &u.GithubID, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// CreateUser create a new user
func CreateUser(ctx context.Context, email, username, nickname, biography, password string, sessionPub string) (*User, error) {
	email = strings.TrimSpace(email)
	if err := validateEmailFormat(ctx, email); err != nil {
		return nil, err
	}
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return nil, session.BadDataError(ctx)
	}
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		nickname = username
	}
	password, err := validateAndEncryptPassword(ctx, password)
	if err != nil {
		return nil, err
	}

	var user *User
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		user, err = createUser(ctx, tx, "", email, username, nickname, password, sessionPub, "", nil)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	UpsertStatistic(ctx, StatisticTypeUsers)
	return user, nil
}

func CreateWeb3User(ctx context.Context, nickname, publicKey, sessionPub, sig string) (*User, error) {
	sigBuf, _ := hex.DecodeString(sig)
	data := fmt.Sprintf("Satellity:%s:%s:%s", nickname, publicKey, sessionPub)
	data = "0x" + hex.EncodeToString(crypto.Keccak256Hash([]byte(data)).Bytes())
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	hash := crypto.Keccak256Hash([]byte(msg))
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), sigBuf)
	if err != nil {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "signature", "invalid", sig)
	}
	pubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "signature", "invalid", sig)
	}
	address := crypto.PubkeyToAddress(*pubKey)
	if publicKey != address.Hex() {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "signature", "invalid", sig)
	}

	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		nickname = publicKey
	}
	var user *User
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		user, err = createUser(ctx, tx, publicKey, "", "", nickname, "", sessionPub, "", nil)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	UpsertStatistic(ctx, StatisticTypeUsers)
	return user, nil
}

// UpdateProfile update user's profile
func (u *User) UpdateProfile(ctx context.Context, nickname, biography string, avatar string) error {
	nickname, biography = strings.TrimSpace(nickname), strings.TrimSpace(biography)
	if len(nickname) == 0 && len(biography) == 0 {
		return nil
	}
	if nickname != "" {
		u.Nickname = nickname
	}
	if biography != "" {
		u.Biography = biography
	}
	u.UpdatedAt = time.Now()
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		cols, posits := durable.PrepareColumnsWithParams([]string{"nickname", "biography", "updated_at"})
		_, err := tx.Exec(ctx, fmt.Sprintf("UPDATE users SET (%s)=(%s) WHERE user_id='%s'", cols, posits, u.UserID), u.Nickname, u.Biography, u.UpdatedAt)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	if len(avatar) > 1024 {
		url, err := clouds.UploadImage(ctx, fmt.Sprintf("/users/%s/cover", u.UserID), avatar)
		if err != nil {
			return session.ServerError(ctx, err)
		}
		_, err = session.Database(ctx).Exec(ctx, "UPDATE users SET avatar_url=$1 WHERE user_id=$2", url, u.UserID)
		if err != nil {
			return session.TransactionError(ctx, err)
		}
		u.AvatarURL = url
	}
	return nil
}

// AuthenticateUser read a user by tokenString. tokenString is a jwt token, more
// about jwt: https://github.com/dgrijalva/jwt-go
func AuthenticateUser(ctx context.Context, tokenString string) (*User, error) {
	var user *User
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, nil
		}
		err := claims.Valid()
		if err != nil {
			return nil, err
		}
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, nil
		}
		uid, sid := fmt.Sprint(claims["uid"]), fmt.Sprint(claims["sid"])
		var s *Session
		err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
			u, err := findUserByID(ctx, tx, uid)
			if err != nil {
				return err
			} else if u == nil {
				return nil
			}
			user = u
			s, err = readSession(ctx, tx, uid, sid)
			if err != nil {
				return err
			} else if s == nil {
				return nil
			}
			user.SessionID = s.SessionID
			return nil
		})
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		if s == nil {
			return nil, nil
		}
		pub, err := hex.DecodeString(s.PublicKey)
		if err != nil {
			return nil, err
		}
		return ed25519.PublicKey(pub), nil
	})
	if err != nil || !token.Valid {
		log.Println("err:::", err, token.Valid)
		return nil, nil
	}
	return user, nil
}

// ReadUsers read users by offset
func ReadUsers(ctx context.Context, offset time.Time) ([]*User, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	rows, err := session.Database(ctx).Query(ctx, fmt.Sprintf("SELECT %s FROM users WHERE created_at<$1 ORDER BY created_at DESC LIMIT 100", strings.Join(userColumns, ",")), offset)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user, err := userFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return users, nil
}

func readUsersByIds(ctx context.Context, tx pgx.Tx, ids []string) ([]*User, error) {
	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT %s FROM users WHERE user_id IN ('%s') LIMIT 100", strings.Join(userColumns, ","), strings.Join(ids, "','")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user, err := userFromRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func readUserSet(ctx context.Context, tx pgx.Tx, ids []string) (map[string]*User, error) {
	users, err := readUsersByIds(ctx, tx, ids)
	if err != nil {
		return nil, err
	}
	set := make(map[string]*User, 0)
	for _, u := range users {
		set[u.UserID] = u
	}
	return set, nil
}

// ReadUser read user by id.
func ReadUser(ctx context.Context, id string) (*User, error) {
	var user *User
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		user, err = findUserByID(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

// ReadUserByUsernameOrEmail read user by identity, which is an email or username.
func ReadUserByUsernameOrEmail(ctx context.Context, identity string) (*User, error) {
	identity = strings.ToLower(strings.TrimSpace(identity))
	if len(identity) < 3 {
		return nil, nil
	}

	var user *User
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		user, err = findUserByIdentity(ctx, tx, identity)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func findUserByIdentity(ctx context.Context, tx pgx.Tx, identity string) (*User, error) {
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM users WHERE username=$1 OR email=$1 LIMIT 1", strings.Join(userColumns, ",")), identity)
	user, err := userFromRow(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAvatar return the avatar of the user
func (u *User) GetAvatar() string {
	if len(u.AvatarURL) > 0 {
		return u.AvatarURL
	}
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=180&d=wavatar", md5.Sum([]byte(strings.ToLower(u.Email.String))))
}

// Role of an user, contains admin and member for now.
func (u *User) GetRole() string {
	if configs.AppConfig.OperatorSet[u.Email.String] {
		return UserRoleAdmin
	}
	if u.Role != "" {
		return u.Role
	}
	return UserRoleMember
}

// Name is nickname or username
func (u *User) Name() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username.String
}

func (u *User) isAdmin() bool {
	return u.GetRole() == UserRoleAdmin
}

func findUserByID(ctx context.Context, tx pgx.Tx, id string) (*User, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM users WHERE user_id=$1", strings.Join(userColumns, ",")), id)
	u, err := userFromRow(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func usersCount(ctx context.Context, tx pgx.Tx) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, "SELECT count(*) FROM users").Scan(&count)
	return count, err
}

func validateAndEncryptPassword(ctx context.Context, password string) (string, error) {
	if len(password) < 8 {
		return password, session.PasswordTooSimpleError(ctx)
	}
	if len(password) > 64 {
		return password, session.BadDataError(ctx)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return password, session.ServerError(ctx, err)
	}
	return string(hashedPassword), nil
}

func isPermit(userID string, user *User) bool {
	return userID == user.UserID || user.isAdmin()
}
