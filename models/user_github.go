package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/config"
	"github.com/godiscourse/godiscourse/external"
	"github.com/godiscourse/godiscourse/session"
	"github.com/godiscourse/godiscourse/uuid"
)

type GithubUser struct {
	Login  string `json:"login"`
	NodeId string `json:"node_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func CreateGithubUser(ctx context.Context, code, sessionSecret string) (*User, error) {
	token, err := fetchAccessToken(ctx, code)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	data, err := fetchOauthUser(ctx, token)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	user, err := findUserByGithubId(ctx, data.NodeId)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if user == nil {
		t := time.Now()
		user = &User{
			UserId:    uuid.NewV4().String(),
			Email:     sql.NullString{data.Email, true},
			Username:  fmt.Sprintf("GH_%s", data.Login),
			Nickname:  data.Name,
			GithubId:  sql.NullString{data.NodeId, true},
			CreatedAt: t,
			UpdatedAt: t,
			IsNew:     true,
		}
	}

	err = session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		if user.IsNew {
			if err := tx.Insert(user); err != nil {
				return err
			}
		}
		sess, err := user.addSession(ctx, tx, sessionSecret)
		if err != nil {
			return err
		}
		user.SessionId = sess.SessionId
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func fetchAccessToken(ctx context.Context, code string) (string, error) {
	client := external.HttpClient()
	client = &http.Client{Timeout: 5 * time.Second}
	data, err := json.Marshal(map[string]interface{}{
		"client_id":     config.GithubClientId,
		"client_secret": config.GithubClientSecret,
		"code":          code,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	return body["access_token"], nil
}

func fetchOauthUser(ctx context.Context, accessToken string) (*GithubUser, error) {
	client := external.HttpClient()
	client = &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GithubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func findUserByGithubId(ctx context.Context, id string) (*User, error) {
	user := &User{}
	if err := session.Database(ctx).Model(user).Column(userCols...).Where("github_id = ?", id).Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
