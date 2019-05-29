package models

import (
	"context"
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

// Group represent the struct of a group
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

// CreateGroup create a group by an user
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

	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		gcols, gparams := durable.PrepareColumnsWithValues(groupColumns)
		query := fmt.Sprintf("INSERT INTO groups(%s) VALUES (%s)", gcols, gparams)
		_, err := tx.ExecContext(ctx, query, group.values()...)
		if err != nil {
			return err
		}
		_, err = createParticipant(ctx, tx, group.GroupID, group.UserID, ParticipantRoleOwner)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return group, nil
}

// UpdateGroup update the group by id
func (user *User) UpdateGroup(mctx *Context, id, name, description string) (*Group, error) {
	ctx := mctx.context
	name, description = strings.TrimSpace(name), strings.TrimSpace(description)
	if name == "" && description == "" {
		return nil, session.BadDataError(ctx)
	}

	var group *Group
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		group, err = findGroup(ctx, tx, id)
		if group.UserID != user.UserID {
			return session.ForbiddenError(ctx)
		}
		if len(name) >= 3 {
			group.Name = name
		}
		if description != "" {
			group.Description = description
		}
		query := "UPDATE groups SET (name, description)=($1,$2) WHERE group_id=$3"
		_, err = tx.ExecContext(ctx, query, group.Name, group.Description, group.GroupID)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return group, nil
}

// ReadGroup read group by an id
func ReadGroup(mctx *Context, id string) (*Group, error) {
	ctx := mctx.context
	var group *Group
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		group, err = findGroup(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return group, nil
}

func findGroup(ctx context.Context, tx *sql.Tx, id string) (*Group, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	row := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT %s FROM groups WHERE group_id=$1", strings.Join(groupColumns, ",")), id)
	group, err := groupFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return group, nil
}

// Participants return members of a group
func (g *Group) Participants(mctx *Context) ([]*Participant, error) {
	ctx := mctx.context
	query := fmt.Sprintf("SELECT %s FROM participants WHERE group_id=$1", strings.Join(participantColumns, ","))
	rows, err := mctx.database.QueryContext(ctx, query, g.GroupID)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	defer rows.Close()

	participants := make([]*Participant, 0)
	for rows.Next() {
		p, err := participantFromRow(rows)
		if err != nil {
			return nil, session.TransactionError(ctx, err)
		}
		participants = append(participants, p)
	}
	return participants, nil
}
