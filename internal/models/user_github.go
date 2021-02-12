package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"satellity/internal/configs"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

// GithubUser is the response body of github oauth.
type GithubUser struct {
	Login  string `json:"login"`
	NodeID string `json:"node_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// CreateGithubUser create a github user. TODO should use createUser
func CreateGithubUser(ctx context.Context, code, sessionSecret string) (*User, error) {
	token, err := fetchAccessToken(ctx, code)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	data, err := fetchOauthUser(ctx, token)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}
	var user *User
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		existing, err := findUserByGithubID(ctx, tx, data.NodeID)
		if err != nil {
			return err
		}
		user, err = createUser(ctx, tx, data.Email, fmt.Sprintf("%s_GH", data.Login), data.Name, "", sessionSecret, data.NodeID, existing)
		if err != nil {
			return nil
		}
		_, err = upsertStatistic(ctx, tx, "users")
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func fetchAccessToken(ctx context.Context, code string) (string, error) {
	config := configs.AppConfig
	client := HTTPClient()
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
	client := HTTPClient()
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
	email, err := fetchUserEmail(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	user.Email = email
	return &user, nil
}

func fetchUserEmail(ctx context.Context, accessToken string) (string, error) {
	client := HTTPClient()
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

func findUserByGithubID(ctx context.Context, tx pgx.Tx, id string) (*User, error) {
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM users WHERE github_id=$1", strings.Join(userColumns, ",")), id)
	return userFromRow(row)
}

// HTTPClient is a client with Timeout (5 seconds).
func HTTPClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}
