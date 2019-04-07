package user

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"strings"

	"github.com/gofrs/uuid"
)

var userColumns = []string{"user_id", "email", "username", "nickname", "biography", "encrypted_password", "github_id", "created_at", "updated_at"}

func (u *Data) values() []interface{} {
	return []interface{}{u.UserID, u.Email, u.Username, u.Nickname, u.Biography, u.EncryptedPassword, u.GithubID, u.CreatedAt, u.UpdatedAt}
}

func userFromRows(row durable.Row) (*Data, error) {
	var u Data
	err := row.Scan(&u.UserID, &u.Email, &u.Username, &u.Nickname, &u.Biography, &u.EncryptedPassword, &u.GithubID, &u.CreatedAt, &u.UpdatedAt)
	return &u, err
}

func findUserByID(ctx context.Context, tx *sql.Tx, id string) (*Data, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE user_id=$1", strings.Join(userColumns, ",")), id)
	u, err := userFromRows(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}
