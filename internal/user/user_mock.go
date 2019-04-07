package user

import (
	"context"
	"database/sql"
	"godiscourse/internal/session"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type UserMock struct {
	users    []*Data
	sessions map[string]*Data
	mock.Mock
}

func NewMock() *UserMock {
	return &UserMock{
		sessions: make(map[string]*Data),
	}
}

// Create mock user.Create, does not affect SessionSecret
func (u *UserMock) Create(ctx context.Context, p *Params) (*Data, error) {
	encryptedPass, err := validateAndEncryptPassword(ctx, p.Password)
	if err != nil {
		return nil, err
	}

	t := time.Now()

	user := &Data{
		Email:             sql.NullString{String: p.Email, Valid: true},
		Username:          p.Username,
		Nickname:          p.Nickname,
		Biography:         p.Biography,
		EncryptedPassword: sql.NullString{String: encryptedPass, Valid: true},
		CreatedAt:         t,
		UpdatedAt:         t,
	}

	// todo: ensure uniqueness
	u.users = append(u.users, user)
	return user, nil
}

func (u *UserMock) Update(ctx context.Context, userData *Data, p *Params) error {
	for _, user := range u.users {
		if user.UserID == userData.UserID {
			user.Username = p.Username
			user.Biography = p.Biography

			userData.Username = p.Username
			userData.Biography = p.Biography
			return nil
		}
	}
	return nil
}

// todo: get user by tokenString
func (u *UserMock) Authenticate(ctx context.Context, tokenString string) (*Data, error) {
	return &Data{
		SessionID: uuid.Must(uuid.NewV4()).String(),
	}, nil
}

func (u *UserMock) GetByOffset(ctx context.Context, offset time.Time) ([]*Data, error) {
	var result []*Data
	for _, user := range u.users {
		if user.CreatedAt.Before(offset) {
			result = append(result, user)
		}
	}
	return result, nil
}

func (u *UserMock) GetByID(ctx context.Context, id string) (*Data, error) {
	var result *Data
	for _, user := range u.users {
		if user.UserID == id {
			result = user
			break
		}
	}
	return result, nil
}

func (u *UserMock) GetByUsernameOrEmail(ctx context.Context, identity string) (*Data, error) {
	var result *Data
	for _, user := range u.users {
		if user.Username == identity || user.Email.String == identity {
			result = user
			break
		}
	}
	return result, nil
}

func (u *UserMock) CreateSession(ctx context.Context, sp *SessionParams) (*Data, error) {
	err := checkSecret(ctx, sp.Secret)
	if err != nil {
		return nil, err
	}

	user, err := u.GetByUsernameOrEmail(ctx, sp.Identity)
	if user == nil {
		return nil, session.IdentityNonExistError(ctx)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword.String), []byte(sp.Password)); err != nil {
		return nil, session.InvalidPasswordError(ctx)
	}

	user.SessionID = uuid.Must(uuid.NewV4()).String()
	u.sessions[user.SessionID] = user
	return user, nil
}
