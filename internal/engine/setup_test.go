package engine

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"

	"godiscourse/internal/configs"
	"godiscourse/internal/durable"
	"godiscourse/internal/models"
	u "godiscourse/internal/user"
	"log"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestContext() (*durable.Database, func()) {
	opts := configs.DefaultOptions()

	db := durable.OpenDatabaseClient(context.Background(), &durable.ConnectionInfo{
		User:     opts.DbUser,
		Password: opts.DbPassword,
		Host:     opts.DbHost,
		Port:     opts.DbPort,
		Name:     opts.DbName,
	})

	if _, err := db.Exec("DROP SCHEMA IF EXISTS godiscourse_test CASCADE"); err != nil {
		log.Panicln(err)
	}
	if _, err := db.Exec("CREATE SCHEMA godiscourse_test"); err != nil {
		log.Panicln(err)
	}

	if _, err := db.Exec("SET search_path='godiscourse_test'"); err != nil {
		log.Panicln(err)
	}

	tables := []string{
		models.UsersDDL,
		models.SessionsDDL,
		models.CategoriesDDL,
		models.TopicsDDL,
		models.CommentsDDL,
		models.StatisticsDDL,
	}
	for _, q := range tables {
		if _, err := db.Exec(q); err != nil {
			log.Panicln(err)
		}
	}

	teardown := func() {
		if _, err := db.Exec("DROP SCHEMA godiscourse_test CASCADE"); err != nil {
			log.Panicln(err)
		}
	}

	return durable.WrapDatabase(db), teardown
}

func seedUsers(user *u.User, t *testing.T) []string {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(t, err)
	public, err := x509.MarshalPKIXPublicKey(priv.Public())
	assert.Nil(t, err)

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
		{"im.yuqlee@gmail.com", "username", "nickname", "", "password", hex.EncodeToString(public), "member", 0, false},
		{"im.yuqlee1@gmail.com", "username1", "nickname1", "", "password", hex.EncodeToString(public), "member", 0, false},
		{"im.yuqlee2@gmail.com", "username2", "nickname2", "", "     pass     ", hex.EncodeToString(public), "member", 1, true},
	}

	var result []string

	for _, uc := range userCases {
		created, err := user.Create(context.Background(), &u.Params{
			Email:         uc.email,
			Username:      uc.username,
			Nickname:      uc.nickname,
			Biography:     uc.biography,
			Password:      uc.password,
			SessionSecret: uc.sessionSecret,
		})

		assert.Nil(t, err)
		result = append(result, created.UserID)
	}

	return result
}
