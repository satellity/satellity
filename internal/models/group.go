package models

import (
	"context"
	"database/sql"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strconv"
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

const (
	MaximumGroupCount = 3
)

// Group represent the struct of a group
type Group struct {
	GroupID     string
	Name        string
	Description string
	UserID      string
	UsersCount  int64
	CreatedAt   time.Time
	UpdateAt    time.Time

	IsMember bool
	User     *User
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
	if !validateGroupFields(name) {
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
		groups, err := findGroupsByUser(ctx, tx, user)
		if err != nil {
			return err
		}
		if len(groups) >= MaximumGroupCount {
			return session.TooManyGroupsError(ctx)
		}
		columns, values := durable.PrepareColumnsWithValues(groupColumns)
		query := fmt.Sprintf("INSERT INTO groups(%s) VALUES (%s)", columns, values)
		_, err = tx.ExecContext(ctx, query, group.values()...)
		if err != nil {
			return err
		}
		_, err = createParticipant(ctx, tx, group.GroupID, group.UserID, ParticipantRoleOwner)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	group.User = user
	return group, nil
}

// UpdateGroup update the group by id
func (user *User) UpdateGroup(mctx *Context, id, name, description string) (*Group, error) {
	ctx := mctx.context
	name, description = strings.TrimSpace(name), strings.TrimSpace(description)
	if len(name) < 3 && description == "" {
		return nil, session.BadDataError(ctx)
	}

	var group *Group
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		group, err = findGroup(ctx, tx, id)
		if group.UserID != user.UserID {
			return session.ForbiddenError(ctx)
		}
		if name != "" {
			group.Name = name
		}
		if description != "" {
			group.Description = description
		}
		query := "UPDATE groups SET (name,description,updated_at)=($1,$2,$3) WHERE group_id=$4"
		_, err = tx.ExecContext(ctx, query, group.Name, group.Description, time.Now(), group.GroupID)
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
func ReadGroup(mctx *Context, id string, current *User) (*Group, error) {
	ctx := mctx.context
	var group *Group
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		group, err = findGroup(ctx, tx, id)
		if err != nil {
			return err
		}
		if group == nil {
			return nil
		}
		if current != nil {
			p, err := findParticipant(ctx, tx, group.GroupID, current.UserID)
			if err != nil {
				return err
			}
			if p != nil {
				group.IsMember = true
			}
		}
		user, err := findUserByID(ctx, tx, group.UserID)
		group.User = user
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
		return nil, err
	}
	return group, nil
}

// Participants return members of a group TODO should support pagination
func (g *Group) Participants(mctx *Context, current *User, offset time.Time, limit string) ([]*User, error) {
	ctx := mctx.context

	if offset.IsZero() {
		offset = time.Now()
	}
	l, _ := strconv.ParseInt(limit, 10, 64)
	if l < 1 || l > 512 {
		l = 512
	}
	if !g.IsMember {
		offset = time.Now()
		l = 64
	}

	users := make([]*User, 0)
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM participants p INNER JOIN users u ON u.user_id=p.user_id WHERE p.group_id=$1 AND p.created_at<$2 ORDER BY p.created_at LIMIT %d", "u."+strings.Join(userColumns, ",u."), l)
		rows, err := mctx.database.QueryContext(ctx, query, g.GroupID, time.Now())
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			user, err := userFromRows(rows)
			if err != nil {
				return err
			}
			users = append(users, user)
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return users, nil
}

//ReadGroups read groups by offset (time) and limit
func ReadGroups(mctx *Context, offset time.Time, limit int64) ([]*Group, error) {
	ctx := mctx.context

	if offset.IsZero() {
		offset = time.Now()
	}
	if limit < 0 || limit > 64 {
		limit = 64
	}

	groups := make([]*Group, 0)
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM groups WHERE created_at<$1 ORDER BY created_at DESC LIMIT $2", strings.Join(groupColumns, ","))
		rows, err := mctx.database.QueryContext(ctx, query, offset, limit)
		if err != nil {
			return err
		}
		defer rows.Close()

		ids := make([]string, 0)
		for rows.Next() {
			group, err := groupFromRow(rows)
			if err != nil {
				return err
			}
			groups = append(groups, group)
			ids = append(ids, group.UserID)
		}
		set, err := readUserSet(ctx, tx, ids)
		if err != nil {
			return err
		}
		for i, group := range groups {
			groups[i].User = set[group.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return groups, nil
}

func ReadGroupsByUser(mctx *Context, userId string) ([]*Group, error) {
	ctx := mctx.context
	groups := make([]*Group, 0)
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		user, err := findUserByID(ctx, tx, userId)
		if err != nil {
			return err
		} else if user == nil {
			return session.NotFoundError(ctx)
		}
		groups, err = findGroupsByUser(ctx, tx, user)
		return err
	})
	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return nil, sessionErr
		}
		return nil, session.TransactionError(ctx, err)
	}
	return groups, nil
}

func (u *User) ReadGroups(mctx *Context) ([]*Group, error) {
	ctx := mctx.context
	groups := make([]*Group, 0)
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		groups, err = findGroupsByUser(ctx, tx, u)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return groups, nil
}

func (u *User) RelatedGroups(mctx *Context, limit int64) ([]*Group, error) {
	ctx := mctx.context

	if limit < 1 || limit > 90 {
		limit = 90
	}
	groups := make([]*Group, 0)
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM groups INNER JOIN participants ON participants.group_id=groups.group_id WHERE participants.user_id=$1 ORDER BY participants.user_id,participants.created_at LIMIT $2", "groups."+strings.Join(groupColumns, ",groups."))
		rows, err := tx.QueryContext(ctx, query, u.UserID, limit)
		if err != nil {
			return err
		}
		defer rows.Close()

		ids := make([]string, 0)
		for rows.Next() {
			group, err := groupFromRow(rows)
			if err != nil {
				return err
			}
			groups = append(groups, group)
			ids = append(ids, group.UserID)
		}
		set, err := readUserSet(ctx, tx, ids)
		if err != nil {
			return err
		}
		for i, group := range groups {
			groups[i].User = set[group.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return groups, nil
}

func findGroupsByUser(ctx context.Context, tx *sql.Tx, u *User) ([]*Group, error) {
	groups := make([]*Group, 0)
	query := fmt.Sprintf("SELECT %s FROM groups WHERE user_id=$1", strings.Join(groupColumns, ","))
	rows, err := tx.QueryContext(ctx, query, u.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		group, err := groupFromRow(rows)
		if err != nil {
			return nil, err
		}
		group.User = u
		groups = append(groups, group)
	}
	return groups, nil
}
