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
	"godiscourse/internal/models"
	"godiscourse/internal/session"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Data contains info of a register user
type Model struct {
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
	CreateGithubUser(context.Context, string, string) (*models.User, error)

	// todo: combine with session and move to standalone auth
	Authenticate(context.Context, string) (*models.User, error)
	GetByUsernameOrEmail(context.Context, string) (*models.User, error)
	Create(context.Context, *Params) (*models.User, error)
	CreateSession(context.Context, *SessionParams) (*models.User, error)
}

type User struct {
	db *durable.Database
}

func New(db *durable.Database) *User {
	return &User{db: db}
}

func (u *User) Create(ctx context.Context, p *Params) (*models.User, error) {
	err := models.CheckSecret(ctx, p.SessionSecret)
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
	user := &models.User{
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
		cols, params := durable.PrepareColumnsWithValues(models.UserColumns)
		_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO users(%s) VALUES (%s)", cols, params), user.Values()...)
		if err != nil {
			return err
		}
		s, err := user.AddSession(ctx, tx, p.SessionSecret)
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
func (u *User) CreateGithubUser(ctx context.Context, code, sessionSecret string) (*models.User, error) {
	token, err := fetchAccessToken(ctx, code)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	data, err := fetchOauthUser(ctx, token)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	var user *models.User
	err = u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		user, err = models.FindUserByGithubID(ctx, tx, data.NodeID)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if user == nil {
		t := time.Now()
		user = &models.User{
			UserID:    uuid.Must(uuid.NewV4()).String(),
			Username:  fmt.Sprintf("GH_%s", data.Login),
			Nickname:  data.Name,
			GithubID:  sql.NullString{String: data.NodeID, Valid: true},
			CreatedAt: t,
			UpdatedAt: t,
			IsNew:     true,
		}
		if data.Email != "" {
			user.Email = sql.NullString{String: data.Email, Valid: true}
		}
	}

	err = u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		if user.IsNew {
			cols, params := durable.PrepareColumnsWithValues(userColumns)
			_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO users(%s) VALUES (%s)", cols, params), user.Values()...)
			if err != nil {
				return err
			}
		}
		s, err := user.AddSession(ctx, tx, sessionSecret)
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

// Authenticate read a user by tokenString. tokenString is a jwt token, more
// about jwt: https://github.com/dgrijalva/jwt-go
func (u *User) Authenticate(ctx context.Context, tokenString string) (*models.User, error) {
	var user *models.User
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, nil
		}
		uid, sid := fmt.Sprint(claims["uid"]), fmt.Sprint(claims["sid"])
		var s *models.Session
		err := u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
			u, err := models.FindUserByID(ctx, tx, uid)
			if err != nil {
				return err
			} else if u == nil {
				return nil
			}
			user = u
			s, err = models.ReadSession(ctx, tx, uid, sid)
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

// GetByUsernameOrEmail read user by identity, which is an email or username.
func (u *User) GetByUsernameOrEmail(ctx context.Context, identity string) (*models.User, error) {
	identity = strings.ToLower(strings.TrimSpace(identity))
	if len(identity) < 3 {
		return nil, nil
	}

	var user *models.User
	err := u.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE username=$1 OR email=$1", strings.Join(models.UserColumns, ",")), identity)
		defer rows.Close()

		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return err
			}
			return nil
		}
		user, err = models.UserFromRows(rows)
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
func (u *User) CreateSession(ctx context.Context, sp *SessionParams) (*models.User, error) {
	err := models.CheckSecret(ctx, sp.Secret)
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
		s, err := user.AddSession(ctx, tx, sp.Secret)
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
