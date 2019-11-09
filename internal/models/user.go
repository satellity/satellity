package models

import (
	"context"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"fmt"
	"satellity/internal/clouds"
	"satellity/internal/configs"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
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
	Email             sql.NullString
	Username          string
	Nickname          string
	AvatarURL         string
	Biography         string
	EncryptedPassword sql.NullString
	GithubID          sql.NullString
	CreatedAt         time.Time
	UpdatedAt         time.Time

	SessionID string
	isNew     bool
}

var userColumns = []string{"user_id", "email", "username", "nickname", "avatar_url", "biography", "encrypted_password", "github_id", "created_at", "updated_at"}

func (u *User) values() []interface{} {
	return []interface{}{u.UserID, u.Email, u.Username, u.Nickname, u.AvatarURL, u.Biography, u.EncryptedPassword, u.GithubID, u.CreatedAt, u.UpdatedAt}
}

func userFromRows(row durable.Row) (*User, error) {
	var u User
	err := row.Scan(&u.UserID, &u.Email, &u.Username, &u.Nickname, &u.AvatarURL, &u.Biography, &u.EncryptedPassword, &u.GithubID, &u.CreatedAt, &u.UpdatedAt)
	return &u, err
}

// CreateUser create a new user
func CreateUser(mctx *Context, email, username, nickname, biography, password string, sessionSecret string) (*User, error) {
	ctx := mctx.context
	data, err := hex.DecodeString(sessionSecret)
	if err != nil {
		return nil, session.BadDataError(ctx)
	}
	public, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, session.BadDataError(ctx)
	}
	switch public.(type) {
	case *ecdsa.PublicKey:
	default:
		return nil, session.BadDataError(ctx)
	}

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
	password, err = validateAndEncryptPassword(ctx, password)
	if err != nil {
		return nil, err
	}

	var user *User
	err = mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		user, err = createUser(ctx, tx, email, username, username, password, sessionSecret, "", nil)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	go upsertStatistic(mctx, "users")
	return user, nil
}

// UpdateProfile update user's profile
func (u *User) UpdateProfile(mctx *Context, nickname, biography string, avatar string) error {
	ctx := mctx.context
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
	columns, params := durable.PrepareColumnsWithParams([]string{"nickname", "biography", "updated_at"})
	_, err := mctx.database.ExecContext(ctx, fmt.Sprintf("UPDATE users SET (%s)=(%s) WHERE user_id='%s'", columns, params, u.UserID), u.Nickname, u.Biography, u.UpdatedAt)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	if len(avatar) > 1024 {
		url, err := clouds.UploadImage(ctx, fmt.Sprintf("/users/%s/cover", u.UserID), avatar)
		if err != nil {
			return session.ServerError(ctx, err)
		}
		_, err = mctx.database.ExecContext(ctx, "UPDATE users SET avatar_url=$1 WHERE user_id=$2", url, u.UserID)
		if err != nil {
			return session.TransactionError(ctx, err)
		}
		u.AvatarURL = url
	}
	return nil
}

// AuthenticateUser read a user by tokenString. tokenString is a jwt token, more
// about jwt: https://github.com/dgrijalva/jwt-go
func AuthenticateUser(mctx *Context, tokenString string) (*User, error) {
	ctx := mctx.context
	var user *User
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, nil
		}
		uid, sid := fmt.Sprint(claims["uid"]), fmt.Sprint(claims["sid"])
		var s *Session
		err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
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
			if _, ok := err.(session.Error); ok {
				return nil, err
			}
			return nil, session.TransactionError(ctx, err)
		}
		pkix, err := hex.DecodeString(s.Secret)
		if err != nil {
			return nil, err
		}
		return x509.ParsePKIXPublicKey(pkix)
	})
	if err != nil || !token.Valid {
		return nil, nil
	}
	return user, nil
}

// ReadUsers read users by offset
func ReadUsers(mctx *Context, offset time.Time) ([]*User, error) {
	ctx := mctx.context
	if offset.IsZero() {
		offset = time.Now()
	}
	rows, err := mctx.database.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE created_at<$1 ORDER BY created_at DESC LIMIT 100", strings.Join(userColumns, ",")), offset)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user, err := userFromRows(rows)
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

func readUsersByIds(ctx context.Context, tx *sql.Tx, ids []string) ([]*User, error) {
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE user_id IN ('%s') LIMIT 100", strings.Join(userColumns, ","), strings.Join(ids, "','")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user, err := userFromRows(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func readUserSet(ctx context.Context, tx *sql.Tx, ids []string) (map[string]*User, error) {
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
func ReadUser(mctx *Context, id string) (*User, error) {
	ctx := mctx.context
	var user *User
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		user, err = findUserByID(ctx, tx, id)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

// ReadUserByUsernameOrEmail read user by identity, which is an email or username.
func ReadUserByUsernameOrEmail(mctx *Context, identity string) (*User, error) {
	ctx := mctx.context
	identity = strings.ToLower(strings.TrimSpace(identity))
	if len(identity) < 3 {
		return nil, nil
	}

	var user *User
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		user, err = findUserByIdentity(ctx, tx, identity)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func findUserByIdentity(ctx context.Context, tx *sql.Tx, identity string) (*User, error) {
	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE username=$1 OR email=$1 LIMIT 1", strings.Join(userColumns, ",")), identity)
	user, err := userFromRows(row)
	if err == sql.ErrNoRows {
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
func (u *User) Role() string {
	if configs.AppConfig.OperatorSet[u.Email.String] {
		return UserRoleAdmin
	}
	return UserRoleMember
}

// Name is nickname or username
func (u *User) Name() string {
	if u.Nickname != "" {
		return u.Nickname
	}
	return u.Username
}

func (u *User) isAdmin() bool {
	return u.Role() == UserRoleAdmin
}

func findUserByID(ctx context.Context, tx *sql.Tx, id string) (*User, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE user_id=$1", strings.Join(userColumns, ",")), id)
	u, err := userFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func usersCount(ctx context.Context, tx *sql.Tx) (int64, error) {
	var count int64
	err := tx.QueryRowContext(ctx, "SELECT count(*) FROM users").Scan(&count)
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

const usersDDL = `
CREATE TABLE IF NOT EXISTS users (
	user_id                VARCHAR(36) PRIMARY KEY,
	email                  VARCHAR(512),
	username               VARCHAR(64) NOT NULL CHECK (username ~* '^[a-z0-9][a-z0-9_]{3,63}$'),
	nickname               VARCHAR(64) NOT NULL DEFAULT '',
	avatar_url             VARCHAR(512) NOT NULL DEFAULT '',
	biography              VARCHAR(2048) NOT NULL DEFAULT '',
	encrypted_password     VARCHAR(1024),
	github_id              VARCHAR(1024) UNIQUE,
	created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_emailx ON users ((LOWER(email)));
CREATE UNIQUE INDEX IF NOT EXISTS users_usernamex ON users ((LOWER(username)));
CREATE INDEX IF NOT EXISTS users_createdx ON users (created_at);
`

const dropUsersDDL = `DROP TABLE IF EXISTS users;`
