package models

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"satellity/internal/session"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)

	public, _, err := ed25519.GenerateKey(rand.Reader)
	assert.Nil(err)

	userCases := []struct {
		email         string
		nickname      string
		biography     string
		password      string
		sessionSecret string
		role          string
		count         int
		valid         bool
	}{
		{"im.yuqlee@gmail.com", "nickname", "", "pass", hex.EncodeToString(public), "member", 0, false},
		{"im.yuqlee@gmail.com", "nickname", "", "     pass     ", hex.EncodeToString(public), "member", 1, true},
	}

	for _, tc := range userCases {
		t.Run(fmt.Sprintf("user username %s", tc.nickname), func(t *testing.T) {
			ctx := setupTestContext()
			defer teardownTestContext(ctx)

			if !tc.valid {
				user, err := CreateUser(ctx, tc.email, tc.nickname, tc.biography, tc.password, tc.sessionSecret)
				assert.NotNil(err)
				assert.Nil(user)
				return
			}

			user, err := CreateUser(ctx, tc.email, tc.nickname, tc.biography, tc.password, tc.sessionSecret)
			assert.Nil(err)
			assert.NotNil(user)

			existing, err := ReadUser(ctx, user.UserID)
			assert.Nil(err)
			assert.NotNil(existing)
			assert.Equal(user.Nickname, existing.Nickname)
			err = bcrypt.CompareHashAndPassword([]byte(existing.EncryptedPassword.String), []byte(tc.password))
			assert.Nil(err)
			existing, err = ReadUser(ctx, uuid.Must(uuid.NewV4()).String())
			assert.Nil(err)
			assert.Nil(existing)
			existing, err = ReadUserByEmail(ctx, "None")
			assert.Nil(err)
			assert.Nil(existing)
			existing, err = ReadUserByEmail(ctx, tc.email)
			assert.Nil(err)
			assert.NotNil(existing)
			existing, err = ReadUserByEmail(ctx, strings.ToUpper(tc.email))
			assert.Nil(err)
			assert.NotNil(existing)
			publicNew, privNew, err := ed25519.GenerateKey(rand.Reader)
			assert.Nil(err)
			existing, err = CreateSession(ctx, tc.email, tc.password, hex.EncodeToString(publicNew))
			assert.Nil(err)
			assert.NotNil(existing)
			assert.Equal(tc.role, user.GetRole())

			sess, err := readTestSession(ctx, existing.UserID, existing.SessionID)
			assert.Nil(err)
			assert.NotNil(sess)
			sess, err = readTestSession(ctx, uuid.Must(uuid.NewV4()).String(), existing.SessionID)
			assert.Nil(err)
			assert.Nil(sess)

			claims := &jwt.MapClaims{
				"uid": existing.UserID,
				"sid": existing.SessionID,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
			ss, err := token.SignedString(privNew)
			assert.Nil(err)
			existing, err = AuthenticateUser(ctx, ss)
			assert.Nil(err)
			assert.NotNil(existing)
			err = existing.UpdateProfile(ctx, "Jason", "", "")
			assert.Nil(err)
			assert.Equal("Jason", existing.Name())
			users, err := ReadUsers(ctx, time.Time{})
			assert.Nil(err)
			assert.Len(users, tc.count)
		})
	}
}

func TestWeb3UserCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := setupTestContext()
	defer teardownTestContext(ctx)

	public, _, err := ed25519.GenerateKey(rand.Reader)
	assert.Nil(err)

	privateKey, err := crypto.HexToECDSA("0123456789012345678901234567890123456789012345678901234567890123")
	assert.Nil(err)

	nickname := "abc"
	publicKey := "0x14791697260E4c9A71f18484C9f997B308e59325"
	data := fmt.Sprintf("Satellite:%s:%s:%s", nickname, publicKey, hex.EncodeToString(public))
	data = "0x" + hex.EncodeToString(crypto.Keccak256Hash([]byte(data)).Bytes())
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	hash := crypto.Keccak256Hash([]byte(msg))
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	assert.Nil(err)

	user, err := CreateWeb3User(ctx, nickname, publicKey, hex.EncodeToString(public), hex.EncodeToString(signature))
	assert.Nil(err)
	assert.NotNil(user)
	assert.Equal(nickname, user.Nickname)

	old, err := CreateWeb3User(ctx, nickname, publicKey, hex.EncodeToString(public), hex.EncodeToString(signature))
	assert.Nil(err)
	assert.NotNil(user)
	assert.Equal(user.UserID, old.UserID)

	nickname = "abcd"
	data = fmt.Sprintf("Satellite:%s:%s:%s", nickname, publicKey, hex.EncodeToString(public))
	data = "0x" + hex.EncodeToString(crypto.Keccak256Hash([]byte(data)).Bytes())
	msg = fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	hash = crypto.Keccak256Hash([]byte(msg))
	signature, err = crypto.Sign(hash.Bytes(), privateKey)
	assert.Nil(err)

	old, err = CreateWeb3User(ctx, nickname, publicKey, hex.EncodeToString(public), hex.EncodeToString(signature))
	assert.Nil(err)
	assert.NotNil(user)
	assert.Equal(user.UserID, old.UserID)
	assert.Equal("abc", old.Nickname)
}

func createTestUser(ctx context.Context, email, password string) *User {
	public, _, _ := ed25519.GenerateKey(rand.Reader)
	user, _ := CreateUser(ctx, email, "nickname", "", password, hex.EncodeToString(public))
	return user
}

func readTestSession(ctx context.Context, uid, sid string) (*Session, error) {
	var s *Session
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		s, err = readSession(ctx, tx, uid, sid)
		return err
	})
	return s, err
}
