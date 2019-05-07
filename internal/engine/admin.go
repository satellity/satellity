package engine

import (
	"context"
	"fmt"
	"godiscourse/internal/models"
	"godiscourse/internal/session"
	"strings"
	"time"
)

func (s *Store) GetUsersByOffset(ctx context.Context, offset time.Time) ([]*models.User, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM users WHERE created_at<$1 ORDER BY created_at DESC LIMIT 100", strings.Join(models.UserColumns, ",")), offset)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := models.UserFromRows(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return users, nil
}
