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

func (u *Model) values() []interface{} {
	return []interface{}{u.UserID, u.Email, u.Username, u.Nickname, u.Biography, u.EncryptedPassword, u.GithubID, u.CreatedAt, u.UpdatedAt}
}

func userFromRows(row durable.Row) (*Model, error) {
	var u Model
	err := row.Scan(&u.UserID, &u.Email, &u.Username, &u.Nickname, &u.Biography, &u.EncryptedPassword, &u.GithubID, &u.CreatedAt, &u.UpdatedAt)
	return &u, err
}

func findUserByID(ctx context.Context, tx *sql.Tx, id string) (*Model, error) {
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

func readUsersByIds(ctx context.Context, tx *sql.Tx, ids []string) ([]*Model, error) {
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE user_id IN ('%s') LIMIT 100", strings.Join(userColumns, ","), strings.Join(ids, "','")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*Model
	for rows.Next() {
		user, err := userFromRows(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
