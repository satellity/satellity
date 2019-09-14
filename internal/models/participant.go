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
)

const participantsDDL = `
CREATE TABLE IF NOT EXISTS participants (
  group_id               VARCHAR(36) NOT NULL REFERENCES groups ON DELETE CASCADE,
  user_id                VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
  role                   VARCHAR(128) NOT NULL,
  source                 VARCHAR(128) NOT NULL,
  expired_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (group_id, user_id)
);

CREATE INDEX IF NOT EXISTS participant_createdx ON participants (created_at);
CREATE INDEX IF NOT EXISTS participant_user_createdx ON participants (user_id,created_at);
CREATE INDEX IF NOT EXISTS participant_group_createdx ON participants (group_id,created_at);
`

// Roles of the participant
const (
	ParticipantRoleOwner  = "OWNER"
	ParticipantRoleAdmin  = "ADMIN"
	ParticipantRoleVIP    = "VIP"
	ParticipantRoleMember = "MEMBER"
	ParticipantRoleGuest  = "GUEST"

	ParticipantSourceInvitation = "INVITATION"
	ParticipantSourcePayment    = "PAYMENT"
)

// Participant represents the struct of a group member
type Participant struct {
	GroupID   string
	UserID    string
	Role      string
	Source    string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdateAt  time.Time
}

var participantColumns = []string{"group_id", "user_id", "role", "source", "expired_at", "created_at", "updated_at"}

func (p *Participant) values() []interface{} {
	return []interface{}{p.GroupID, p.UserID, p.Role, p.Source, p.ExpiredAt, p.CreatedAt, p.UpdateAt}
}

func participantFromRow(row durable.Row) (*Participant, error) {
	var p Participant
	err := row.Scan(&p.GroupID, &p.UserID, &p.Role, &p.Source, &p.ExpiredAt, &p.CreatedAt, &p.UpdateAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func createParticipant(ctx context.Context, tx *sql.Tx, group *Group, userID, source string) (*Participant, error) {
	t := time.Now()
	p := &Participant{
		GroupID:   group.GroupID,
		UserID:    userID,
		Role:      group.Role,
		Source:    source,
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
func (user *User) JoinGroup(mctx *Context, groupID, role string) (*Group, error) {
	ctx := mctx.context
	switch role {
	case ParticipantRoleAdmin,
		ParticipantRoleVIP,
		ParticipantRoleMember:
	default:
		return nil, session.BadDataError(ctx)
	}
	var group *Group
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		group, err = findGroup(ctx, tx, groupID)
		if err != nil || group == nil {
			return err
		}
		p, err := findParticipant(ctx, tx, groupID, user.UserID)
		if err != nil {
			return err
		} else if p != nil {
			return nil
		}
		owner, err := findUserByID(ctx, tx, group.UserID)
		if err != nil {
			return err
		}
		group.User = owner

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
		group.Role = role
		_, err = createParticipant(ctx, tx, group, user.UserID, ParticipantSourcePayment)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return group, nil
}

// ExitGroup exit the group by id
func (user *User) ExitGroup(mctx *Context, groupID string) (*Group, error) {
	ctx := mctx.context
	var group *Group
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		group, err = findGroup(ctx, tx, groupID)
		if err != nil {
			return err
		} else if group == nil {
			return nil
		}
		p, err := findParticipant(ctx, tx, groupID, user.UserID)
		if err != nil {
			return err
		} else if p == nil {
			return nil
		} else if p.Role == ParticipantRoleOwner {
			return nil
		}
		owner, err := findUserByID(ctx, tx, group.UserID)
		if err != nil {
			return err
		}
		group.User = owner

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
		group.Role = ParticipantRoleGuest
		_, err = tx.ExecContext(ctx, "DELETE FROM participants WHERE group_id=$1 AND user_id=$2", group.GroupID, user.UserID)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
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
	if g.GetRole() == ParticipantRoleGuest {
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

func (g *Group) UpdateParticipant(mctx *Context, current *User, id, role string) error {
	ctx := mctx.context
	switch role {
	case ParticipantRoleAdmin,
		ParticipantRoleVIP,
		ParticipantRoleMember:
	default:
		return session.BadDataError(ctx)
	}
	if current.UserID == id {
		return session.BadDataError(ctx)
	}

	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		p, err := findParticipant(ctx, tx, g.GroupID, current.UserID)
		if err != nil || p == nil {
			return err
		}
		switch p.Role {
		case ParticipantRoleOwner:
		case ParticipantRoleAdmin:
			if role == ParticipantRoleAdmin {
				return nil
			}
		default:
			return nil
		}
		p, err = findParticipant(ctx, tx, g.GroupID, id)
		if err != nil || p == nil || p.Role == role {
			return err
		}
		p.Role = role
		_, err = tx.ExecContext(ctx, "UPDATE participants SET role=$1 WHERE group_id=$2 AND user_id=$3", role, p.GroupID, p.UserID)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}
