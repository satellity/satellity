package models

import (
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

const groupsDDL = `
CREATE TABLE IF NOT EXISTS groups (
	group_id               VARCHAR(36) PRIMARY KEY,
	name                   VARCHAR(512) NOT NULL,
	description            TEXT NOT NULL,
	user_id                VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	users_count            BIGINT NOT NULL DEFAULT 0,
	created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS groups_userx ON groups (user_id);
`

type Group struct {
	GroupID     string
	Name        string
	Description string
	UserID      string
	UsersCount  int64
	CreatedAt   time.Time
	UpdateAt    time.Time
}

var groupColumns = []string{"group_id", "name", "description", "user_id", "users_count", "created_at", "updated_at"}

func (g *Group) values() []interface{} {
	return []interface{}{g.GroupID, g.Name, g.Description, g.UserID, g.UsersCount, g.CreatedAt, g.UpdateAt}
}

func groupFromRow(row durable.Row) (*Group, error) {
	var g Group
	err := row.Scan(&g.GroupID, &g.Name, &g.Description, &g.UserID, &g.UsersCount, &g.CreatedAt, &g.UpdateAt)
	return &g, err
}

func (user *User) CreateGroup(mctx *Context, name, description string) (*Group, error) {
	ctx := mctx.context

	if len(name) < 3 {
		return nil, session.BadDataError(ctx)
	}
	t := time.Now()
	group := &Group{
		GroupID:     uuid.Must(uuid.NewV4()).String(),
		Name:        name,
		Description: description,
		UserID:      user.UserID,
		UsersCount:  1,
		CreatedAt:   t,
		UpdateAt:    t,
	}

	participant := &Participant{
		GroupID:   group.GroupID,
		UserID:    user.UserID,
		Role:      ParticipantRoleOwner,
		CreatedAt: t,
		UpdateAt:  t,
	}

	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		gcols, gparams := durable.PrepareColumnsWithValues(groupColumns)
		query := fmt.Sprintf("INSERT INTO groups(%s) VALUES (%s)", gcols, gparams)
		_, err := tx.ExecContext(ctx, query, group.values()...)
		if err != nil {
			return err
		}
		pcols, pparams := durable.PrepareColumnsWithValues(participantColumns)
		pquery := fmt.Sprintf("INSERT INTO participants(%s) VALUES (%s)", pcols, pparams)
		_, err = tx.ExecContext(ctx, pquery, participant.values()...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return group, nil
}

func ReadGroup(mctx *Context, id string) (*Group, error) {
	ctx := mctx.context
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row, err := mctx.database.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM groups WHERE group_id=$1", strings.Join(groupColumns, ",")))
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	group, err := groupFromRow(row)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return group, nil
}
