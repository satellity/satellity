package models

import (
	"context"
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"strings"
	"time"
)

const participantsDDL = `
CREATE TABLE IF NOT EXISTS participants (
	group_id               VARCHAR(36) NOT NULL REFERENCES groups ON DELETE CASCADE,
	user_id                VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	role                   VARCHAR(128) NOT NULL,
	created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	PRIMARY KEY (group_id, user_id)
);

CREATE INDEX IF NOT EXISTS participant_createdx ON participants (created_at);
CREATE INDEX IF NOT EXISTS participant_group_createdx ON participants (group_id,created_at);
`

// Roles of the participant
const (
	ParticipantRoleOwner  = "OWNER"
	ParticipantRoleAdmin  = "ADMIN"
	ParticipantRoleVIP    = "VIP"
	ParticipantRoleMember = "MEMBER"
)

// Participant represents the struct of a group member
type Participant struct {
	GroupID   string
	UserID    string
	Role      string
	CreatedAt time.Time
	UpdateAt  time.Time
}

var participantColumns = []string{"group_id", "user_id", "role", "created_at", "updated_at"}

func (p *Participant) values() []interface{} {
	return []interface{}{p.GroupID, p.UserID, p.Role, p.CreatedAt, p.UpdateAt}
}

func participantFromRow(row durable.Row) (*Participant, error) {
	var p Participant
	err := row.Scan(&p.GroupID, &p.UserID, &p.Role, &p.CreatedAt, &p.UpdateAt)
	return &p, err
}

func createParticipant(ctx context.Context, tx *sql.Tx, groupID, userID, role string) (*Participant, error) {
	t := time.Now()
	p := &Participant{
		GroupID:   groupID,
		UserID:    userID,
		Role:      role,
		CreatedAt: t,
		UpdateAt:  t,
	}

	columns, params := durable.PrepareColumnsWithValues(participantColumns)
	query := fmt.Sprintf("INSERT INTO participants(%s) VALUES (%s)", columns, params)
	_, err := tx.ExecContext(ctx, query, p.values()...)
	return p, err
}

func findParticipant(ctx context.Context, tx *sql.Tx, groupID, userID string) (*Participant, error) {
	query := fmt.Sprintf("SELECT %s FROM participants WHERE group_id=$1 AND user_id=$2", strings.Join(participantColumns, ","))
	p, err := participantFromRow(tx.QueryRowContext(ctx, query, groupID, userID))
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return p, nil
}

// JoinGroup join the group by id
func (user *User) JoinGroup(mctx *Context, groupID, role string) error {
	ctx := mctx.context
	switch role {
	case ParticipantRoleAdmin,
		ParticipantRoleVIP,
		ParticipantRoleMember:
	default:
		return session.BadDataError(ctx)
	}
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		group, err := findGroup(ctx, tx, groupID)
		if err != nil {
			return err
		} else if group == nil {
			return session.NotFoundError(ctx)
		}
		p, err := findParticipant(ctx, tx, groupID, user.UserID)
		if err != nil {
			return err
		} else if p != nil {
			return nil
		}

		var count int64
		err = tx.QueryRowContext(ctx, "SELECT count(*) FROM participants WHERE group_id=$1", groupID).Scan(&count)
		if err != nil {
			return err
		}
		group.UsersCount = count + 1
		_, err = tx.ExecContext(ctx, "UPDATE groups SET users_count=$1 WHERE group_id=$2", group.UsersCount, group.GroupID)
		if err != nil {
			return err
		}
		_, err = createParticipant(ctx, tx, groupID, user.UserID, role)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return err
		}
		return session.TransactionError(ctx, err)
	}
	return nil
}

// ExitGroup exit the group by id
func (user *User) ExitGroup(mctx *Context, groupID string) error {
	ctx := mctx.context
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		group, err := findGroup(ctx, tx, groupID)
		if err != nil {
			return err
		} else if group == nil {
			return session.NotFoundError(ctx)
		}
		p, err := findParticipant(ctx, tx, groupID, user.UserID)
		if err != nil {
			return err
		} else if p == nil {
			return nil
		} else if p.Role == ParticipantRoleAdmin {
			return nil
		}

		var count int64
		err = tx.QueryRowContext(ctx, "SELECT count(*) FROM participants WHERE group_id=$1", groupID).Scan(&count)
		if err != nil {
			return err
		}
		group.UsersCount = count - 1
		_, err = tx.ExecContext(ctx, "UPDATE groups SET users_count=$1 WHERE group_id=$2", group.UsersCount, group.GroupID)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "DELETE FROM participants WHERE group_id=$1 AND user_id=$2", group.GroupID, user.UserID)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return err
		}
		return session.TransactionError(ctx, err)
	}
	return nil
}
