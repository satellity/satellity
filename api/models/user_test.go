package models

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/godiscourse/godiscourse/api/session"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer session.Database(ctx).Close()
	defer teardownTestContext(ctx)

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	public, err := x509.MarshalPKIXPublicKey(priv.Public())
	assert.Nil(err)

	userCases := []struct {
		email         string
		username      string
		nickname      string
		biography     string
		password      string
		sessionSecret string
		role          string
		count         int
		valid         bool
	}{
		{"im.yuqlee@gmailabcefgh.com", "username", "nickname", "", "password", hex.EncodeToString(public), "member", 0, false},
		{"im.yuqlee@gmail.com", "username", "nickname", "", "pass", hex.EncodeToString(public), "member", 0, false},
		{"im.yuqlee@gmail.com", "username", "nickname", "", "     pass     ", hex.EncodeToString(public), "member", 1, true},
	}

	for _, tc := range userCases {
		t.Run(fmt.Sprintf("user username %s", tc.username), func(t *testing.T) {
			if !tc.valid {
				user, err := CreateUser(ctx, tc.email, tc.username, tc.nickname, tc.biography, tc.password, tc.sessionSecret)
				assert.NotNil(err)
				assert.Nil(user)
				return
			}

			user, err := CreateUser(ctx, tc.email, tc.username, tc.nickname, tc.biography, tc.password, tc.sessionSecret)
			assert.Nil(err)
			assert.NotNil(user)

			new, err := ReadUser(ctx, user.UserID)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(user.Username, new.Username)
			assert.Equal(user.Nickname, new.Nickname)
			err = bcrypt.CompareHashAndPassword([]byte(new.EncryptedPassword.String), []byte(tc.password))
			assert.Nil(err)
			new, err = ReadUserByUsernameOrEmail(ctx, "None")
			assert.Nil(err)
			assert.Nil(new)
			new, err = ReadUserByUsernameOrEmail(ctx, tc.email)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = ReadUserByUsernameOrEmail(ctx, tc.username)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = ReadUserByUsernameOrEmail(ctx, strings.ToUpper(tc.email))
			assert.Nil(err)
			assert.NotNil(new)
			new, err = CreateSession(ctx, tc.email, tc.password, hex.EncodeToString(public))
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(tc.username, user.Username)
			assert.Equal(tc.role, user.Role())

			sess, err := readTestSession(ctx, new.UserID, new.SessionID)
			assert.Nil(err)
			assert.NotNil(sess)
			sess, err = readTestSession(ctx, uuid.Must(uuid.NewV4()).String(), new.SessionID)
			assert.Nil(err)
			assert.Nil(sess)

			claims := &jwt.MapClaims{
				"uid": new.UserID,
				"sid": new.SessionID,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
			ss, err := token.SignedString(priv)
			assert.Nil(err)
			new, err = AuthenticateUser(ctx, ss)
			assert.Nil(err)
			assert.NotNil(new)
			err = new.UpdateProfile(ctx, "Jason", "")
			assert.Nil(err)
			assert.Equal("Jason", new.Name())
			new, err = ReadUserByUsernameOrEmail(ctx, tc.username)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal("Jason", new.Name())
			users, err := ReadUsers(ctx, time.Time{})
			assert.Nil(err)
			assert.Len(users, tc.count)
		})
	}
}

func createTestUser(ctx context.Context, email, username, password string) *User {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	public, _ := x509.MarshalPKIXPublicKey(priv.Public())
	user, _ := CreateUser(ctx, email, username, "nickname", "", password, hex.EncodeToString(public))
	return user
}

func readTestSession(ctx context.Context, uid, sid string) (*Session, error) {
	var s *Session
	err := session.Database(ctx).RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		s, err = readSession(ctx, tx, uid, sid)
		return err
	})
	return s, err
}
