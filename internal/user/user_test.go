package user

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCRUD(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	userMock := NewMock()

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
				user, err := userMock.Create(ctx, &Params{
					Email:         tc.email,
					Username:      tc.username,
					Nickname:      tc.nickname,
					Biography:     tc.biography,
					Password:      tc.password,
					SessionSecret: tc.sessionSecret,
				})
				assert.NotNil(err)
				assert.Nil(user)
				return
			}

			user, err := userMock.Create(ctx, &Params{
				Email:         tc.email,
				Username:      tc.username,
				Nickname:      tc.nickname,
				Biography:     tc.biography,
				Password:      tc.password,
				SessionSecret: tc.sessionSecret,
			})
			assert.Nil(err)
			assert.NotNil(user)

			new, err := userMock.GetByID(ctx, user.UserID)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(user.Username, new.Username)
			assert.Equal(user.Nickname, new.Nickname)
			err = bcrypt.CompareHashAndPassword([]byte(new.EncryptedPassword.String), []byte(tc.password))
			assert.Nil(err)
			new, err = userMock.GetByID(ctx, uuid.Must(uuid.NewV4()).String())
			assert.Nil(err)
			assert.Nil(new)
			new, err = userMock.GetByUsernameOrEmail(ctx, "None")
			assert.Nil(err)
			assert.Nil(new)
			new, err = userMock.GetByUsernameOrEmail(ctx, tc.email)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = userMock.GetByUsernameOrEmail(ctx, tc.username)
			assert.Nil(err)
			assert.NotNil(new)
			new, err = userMock.GetByUsernameOrEmail(ctx, strings.ToUpper(tc.email))
			assert.Nil(err)
			assert.NotNil(new)
			new, err = userMock.CreateSession(ctx, &SessionParams{
				Identity: tc.email,
				Password: tc.password,
				Secret:   hex.EncodeToString(public),
			})
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal(tc.username, user.Username)
			assert.Equal(tc.role, user.Role())

			// sess, err := readTestSession(ctx, new.UserID, new.SessionID)
			// assert.Nil(err)
			// assert.NotNil(sess)
			// sess, err = readTestSession(ctx, uuid.Must(uuid.NewV4()).String(), new.SessionID)
			// assert.Nil(err)
			// assert.Nil(sess)

			claims := &jwt.MapClaims{
				"uid": new.UserID,
				"sid": new.SessionID,
			}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
			ss, err := token.SignedString(priv)
			assert.Nil(err)
			new, err = userMock.Authenticate(ctx, ss)
			assert.Nil(err)
			assert.NotNil(new)
			err = userMock.Update(ctx, new, &Params{
				Username:  "Jason",
				Biography: "",
			})
			assert.Nil(err)
			assert.Equal("Jason", new.Name())
			new, err = userMock.GetByUsernameOrEmail(ctx, tc.username)
			assert.Nil(err)
			assert.NotNil(new)
			assert.Equal("Jason", new.Name())
			users, err := userMock.GetByOffset(ctx, time.Time{})
			assert.Nil(err)
			assert.Len(users, tc.count)
		})
	}
}

// func createTestUser(mctx *Context, email, username, password string) *Data {
// 	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	public, _ := x509.MarshalPKIXPublicKey(priv.Public())
// 	user, _ := CreateUser(mctx, email, username, "nickname", "", password, hex.EncodeToString(public))
// 	return user
// }
