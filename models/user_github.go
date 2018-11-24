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

// GithubUser is the response body of github oauth.
type GithubUser struct {
	Login  string `json:"login"`
	NodeID string `json:"node_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// CreateGithubUser create a github user.
func CreateGithubUser(ctx context.Context, code, sessionSecret string) (*User, error) {
	token, err := fetchAccessToken(ctx, code)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	data, err := fetchOauthUser(ctx, token)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	user, err := findUserByGithubId(ctx, data.NodeID)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if user == nil {
		t := time.Now()
		user = &User{
			UserID:    uuid.NewV4().String(),
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

	err = session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		if user.isNew {
			if err := tx.Insert(user); err != nil {
				return err
			}
		}
		sess, err := user.addSession(ctx, tx, sessionSecret)
		if err != nil {
			return err
		}
		user.SessionID = sess.SessionID
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func fetchAccessToken(ctx context.Context, code string) (string, error) {
	client := external.HttpClient()
	data, err := json.Marshal(map[string]interface{}{
		"client_id":     config.GithubClientID,
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
	if body["error"] != "" {
		return "", session.ServerError(ctx, fmt.Errorf("%v", body))
	}
	return body["access_token"], nil
}

func fetchOauthUser(ctx context.Context, accessToken string) (*GithubUser, error) {
	client := external.HttpClient()
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
	email, err := featchUserEmail(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	user.Email = email
	return &user, nil
}

func featchUserEmail(ctx context.Context, accessToken string) (string, error) {
	client := external.HttpClient()
	req, err := http.NewRequest("GET", "https://api.github.com/user/public_emails", nil)
	if err != nil {
		return "", err
	}
	req.Close = true
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}
	if len(emails) == 0 {
		return "", nil
	}
	return emails[0].Email, nil
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
