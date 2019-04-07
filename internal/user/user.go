package user

import (
	"context"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"godiscourse/internal/durable"
	"godiscourse/internal/session"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Data contains info of a register user
type Data struct {
	UserID            string
	Email             sql.NullString
	Username          string
	Nickname          string
	Biography         string
	EncryptedPassword sql.NullString
	GithubID          sql.NullString
	CreatedAt         time.Time
	UpdatedAt         time.Time

	SessionID string
	isNew     bool
}

type Params struct {
	Email         string
	Username      string
	Nickname      string
	Biography     string
	Password      string
	SessionSecret string
}

type SessionParams struct {
	Identity string
	Password string
	Secret   string
}

type UserDatastore interface {
	Create(context.Context, *Params) (*Data, error)
	CreateGithubUser(context.Context, string, string) (*Data, error)
	Update(context.Context, *Data, *Params) error
	Authenticate(context.Context, string) (*Data, error)
	GetByOffset(context.Context, time.Time) ([]*Data, error)
	GetByID(context.Context, string) (*Data, error) // equal ReadUser
	GetByUsernameOrEmail(context.Context, string) (*Data, error)
	CreateSession(context.Context, *SessionParams) (*Data, error)
}

type User struct {
	db *durable.Database
}

func New(db *durable.Database) *User {
	return &User{db: db}
}

func (u *User) Create(ctx context.Context, p *Params) (*Data, error) {
	err := checkSecret(ctx, p.SessionSecret)
	if err != nil {
		return nil, err
	}

	p.Email = strings.TrimSpace(p.Email)
	if err := validateEmailFormat(ctx, p.Email); err != nil {
		return nil, err
	}
	p.Username = strings.TrimSpace(p.Username)
	if len(p.Username) < 3 {
		return nil, session.BadDataError(ctx)
	}
	p.Nickname = strings.TrimSpace(p.Nickname)
	if p.Nickname == "" {
		p.Nickname = p.Username
	}
	p.Password, err = validateAndEncryptPassword(ctx, p.Password)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	user := &Data{
		UserID:            uuid.Must(uuid.NewV4()).String(),
		Email:             sql.NullString{String: p.Email, Valid: true},
		Username:          p.Username,
		Nickname:          p.Nickname,
		Biography:         p.Biography,
		EncryptedPassword: sql.NullString{String: p.Password, Valid: true},
		CreatedAt:         t,
		UpdatedAt:         t,
	}

	err = u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		cols, params := durable.PrepareColumnsWithValues(userColumns)
		_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO users(%s) VALUES (%s)", cols, params), user.values()...)
		if err != nil {
			return err
		}
		s, err := user.addSession(ctx, tx, p.SessionSecret)
		if err != nil {
			return err
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
	return user, nil
}

// CreateGithubUser create a github user.
func (u *User) CreateGithubUser(ctx context.Context, code, sessionSecret string) (*Data, error) {
	token, err := fetchAccessToken(ctx, code)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	data, err := fetchOauthUser(ctx, token)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	var user *Data
	err = u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		user, err = findUserByGithubID(ctx, tx, data.NodeID)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if user == nil {
		t := time.Now()
		user = &Data{
			UserID:    uuid.Must(uuid.NewV4()).String(),
			Username:  fmt.Sprintf("GH_%s", data.Login),
			Nickname:  data.Name,
			GithubID:  sql.NullString{String: data.NodeID, Valid: true},
			CreatedAt: t,
			UpdatedAt: t,
			isNew:     true,
		}
		if data.Email != "" {
			user.Email = sql.NullString{String: data.Email, Valid: true}
		}
	}

	err = u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		if user.isNew {
			cols, params := durable.PrepareColumnsWithValues(userColumns)
			_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO users(%s) VALUES (%s)", cols, params), user.values()...)
			if err != nil {
				return err
			}
		}
		s, err := user.addSession(ctx, tx, sessionSecret)
		if err != nil {
			return err
		}
		user.SessionID = s.SessionID
		// todo: fix it
		// _, err = upsertStatistic(ctx, tx, "users")
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

// Update update user's profile
func (u *User) Update(ctx context.Context, current *Data, new *Params) error {
	nickname, biography := strings.TrimSpace(new.Nickname), strings.TrimSpace(new.Biography)
	if len(nickname) == 0 && len(biography) == 0 {
		return nil
	}
	if nickname != "" {
		current.Nickname = nickname
	}
	if biography != "" {
		current.Biography = biography
	}
	current.UpdatedAt = time.Now()
	cols, params := durable.PrepareColumnsWithValues([]string{"nickname", "biography", "updated_at"})
	_, err := u.db.ExecContext(ctx, fmt.Sprintf("UPDATE users SET (%s)=(%s) WHERE user_id='%s'", cols, params, current.UserID), current.Nickname, current.Biography, current.UpdatedAt)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

// Authenticate read a user by tokenString. tokenString is a jwt token, more
// about jwt: https://github.com/dgrijalva/jwt-go
func (u *User) Authenticate(ctx context.Context, tokenString string) (*Data, error) {
	var user *Data
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
		err := u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
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

func (u *User) GetByID(ctx context.Context, id string) (*Data, error) {
	var user *Data
	err := u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
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

// GetByUsernameOrEmail read user by identity, which is an email or username.
func (u *User) GetByUsernameOrEmail(ctx context.Context, identity string) (*Data, error) {
	identity = strings.ToLower(strings.TrimSpace(identity))
	if len(identity) < 3 {
		return nil, nil
	}

	var user *Data
	err := u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE username=$1 OR email=$1", strings.Join(userColumns, ",")), identity)
		defer rows.Close()

		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return err
			}
			return nil
		}
		user, err = userFromRows(rows)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return user, session.TransactionError(ctx, err)
	}
	return user, nil
}

// CreateSession create a new user session
func (u *User) CreateSession(ctx context.Context, sp *SessionParams) (*Data, error) {
	err := checkSecret(ctx, sp.Secret)
	if err != nil {
		return nil, err
	}

	user, err := u.GetByUsernameOrEmail(ctx, sp.Identity)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, session.IdentityNonExistError(ctx)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword.String), []byte(sp.Password)); err != nil {
		return nil, session.InvalidPasswordError(ctx)
	}

	err = u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		s, err := user.addSession(ctx, tx, sp.Secret)
		if err != nil {
			return err
		}
		user.SessionID = s.SessionID
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}
