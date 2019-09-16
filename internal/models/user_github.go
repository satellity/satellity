package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"satellity/internal/configs"
	"satellity/internal/durable"
	"satellity/internal/external"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// GithubUser is the response body of github oauth.
type GithubUser struct {
	Login  string `json:"login"`
	NodeID string `json:"node_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// CreateGithubUser create a github user. TODO should use createUser
func CreateGithubUser(mctx *Context, code, sessionSecret string) (*User, error) {
	ctx := mctx.context
	token, err := fetchAccessToken(ctx, code)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	data, err := fetchOauthUser(ctx, token)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	var user *User
	err = mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		user, err = findUserByGithubID(ctx, tx, data.NodeID)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if user == nil {
		t := time.Now()
		user = &User{
			UserID:    uuid.Must(uuid.NewV4()).String(),
			Username:  fmt.Sprintf("%s_GH", data.Login),
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

	err = mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
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
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	go upsertStatistic(mctx, "users")
	return user, nil
}

func fetchAccessToken(ctx context.Context, code string) (string, error) {
	config := configs.AppConfig
	client := external.HTTPClient()
	data, err := json.Marshal(map[string]interface{}{
		"client_id":     config.Github.ClientID,
		"client_secret": config.Github.ClientSecret,
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
	client := external.HTTPClient()
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
	client := external.HTTPClient()
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

func findUserByGithubID(ctx context.Context, tx *sql.Tx, id string) (*User, error) {
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE github_id=$1", strings.Join(userColumns, ",")), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return userFromRows(rows)
}
