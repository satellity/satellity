package models

import (
	"context"
	"database/sql"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

const messagesDDL = `
CREATE TABLE IF NOT EXISTS messages (
	message_id           VARCHAR(36) PRIMARY KEY,
	body                 TEXT NOT NULL,
	group_id             VARCHAR(36) NOT NULL REFERENCES groups ON DELETE CASCADE,
	user_id              VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	parent_id            VARCHAR(36) NOT NULL,
	created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS messages_group_created_parentx ON messages (group_id, created_at DESC, parent_id);
CREATE INDEX IF NOT EXISTS messages_parent_createdx ON messages (parent_id, created_at DESC);
`

// Message represent the struct of a message
type Message struct {
	MessageID string
	Body      string
	GroupID   string
	UserID    string
	ParentID  string
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User
}

var messageColumns = []string{"message_id", "body", "group_id", "user_id", "parent_id", "created_at", "updated_at"}

func (m *Message) values() []interface{} {
	return []interface{}{m.MessageID, m.Body, m.GroupID, m.UserID, m.ParentID, m.CreatedAt, m.UpdatedAt}
}

func messageFromRow(row durable.Row) (*Message, error) {
	var m Message
	err := row.Scan(&m.MessageID, &m.Body, &m.GroupID, &m.UserID, &m.ParentID, &m.CreatedAt, &m.UpdatedAt)
	return &m, err
}

// CreateMessage create a message
func (u *User) CreateMessage(mctx *Context, groupID, body, parentID string) (*Message, error) {
	ctx := mctx.context
	body = strings.TrimSpace(body)
	if len(body) < 1 {
		return nil, session.BadDataError(ctx)
	}
	var message *Message
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		p, err := findParticipant(ctx, tx, groupID, u.UserID)
		if err != nil {
			return err
		} else if p == nil {
			return session.ForbiddenError(ctx)
		}

		id := uuid.Must(uuid.NewV4()).String()
		parent, err := findMessageByID(ctx, tx, parentID)
		parentID = id
		if err != nil {
			return err
		} else if parent != nil {
			parentID = parent.MessageID
		}

		t := time.Now()
		message = &Message{
			MessageID: id,
			Body:      body,
			GroupID:   groupID,
			UserID:    u.UserID,
			ParentID:  parentID,
			CreatedAt: t,
			UpdatedAt: t,
		}

		columns, params := durable.PrepareColumnsWithValues(messageColumns)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO messages (%s) VALUES (%s)", columns, params), message.values()...)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	message.User = u
	return message, nil
}

// ReadMessage read a message by id
func ReadMessage(mctx *Context, id string) (*Message, error) {
	ctx := mctx.context

	var message *Message
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		message, err = findMessageByID(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return message, nil
}

func findMessageByID(ctx context.Context, tx *sql.Tx, id string) (*Message, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT %s FROM messages WHERE message_id=$1", strings.Join(messageColumns, ","))
	row := tx.QueryRowContext(ctx, query, id)
	m, err := messageFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return m, err
}

// ReadMessages read all messages, parameters: offset default time.Now()
func (g *Group) ReadMessages(mctx *Context, offset time.Time) ([]*Message, error) {
	ctx := mctx.context
	if offset.IsZero() {
		offset = time.Now()
	}
	limit := 16
	if g.Role == ParticipantRoleGuest {
		offset = time.Now()
		limit = 8
	}

	var messages []*Message
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		query := "SELECT message_id FROM messages WHERE group_id=$1 AND created_at<$2 AND message_id=parent_id ORDER BY group_id,created_at DESC,parent_id LIMIT $3"
		rows, err := tx.QueryContext(ctx, query, g.GroupID, offset, limit)
		if err != nil {
			return err
		}
		defer rows.Close()

		userIds := make([]string, 0)
		messageIds := make([]string, 0)
		for rows.Next() {
			var id string
			err := rows.Scan(&id)
			if err != nil {
				return err
			}
			messageIds = append(messageIds, id)
		}
		if err := rows.Err(); err != nil {
			return err
		}

		query = fmt.Sprintf("SELECT %s FROM messages WHERE parent_id IN ('%s') ORDER BY created_at DESC LIMIT $1", strings.Join(messageColumns, ","), strings.Join(messageIds, "','"))
		frows, err := tx.QueryContext(ctx, query, limit)
		if err != nil {
			return err
		}
		defer frows.Close()

		for frows.Next() {
			msg, err := messageFromRow(frows)
			if err != nil {
				return err
			}
			userIds = append(userIds, msg.UserID)
			messages = append(messages, msg)
		}
		if err := frows.Err(); err != nil {
			return err
		}

		userSet, err := readUserSet(ctx, tx, userIds)
		if err != nil {
			return err
		}
		for i, msg := range messages {
			messages[i].User = userSet[msg.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return messages, nil
}

// UpdateMessage update a message by id
func (user *User) UpdateMessage(mctx *Context, id, body string) (*Message, error) {
	ctx := mctx.context
	var message *Message
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		message, err = findMessageByID(ctx, tx, id)
		if err != nil {
			return err
		} else if message == nil {
			return session.ForbiddenError(ctx)
		}
		if message.UserID != user.UserID && !user.isAdmin() {
			return session.ForbiddenError(ctx)
		}

		body = strings.TrimSpace(body)
		if len(body) < 1 {
			return session.BadDataError(ctx)
		}
		message.Body = body
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE messages SET body=$1 WHERE message_id=$2"), body, id)
		return err
	})
	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return nil, sessionErr
		}
		return nil, session.TransactionError(ctx, err)
	}
	message.User = user
	return message, nil
}

// DeleteMessage delete a message by id
func (user *User) DeleteMessage(mctx *Context, id string) error {
	ctx := mctx.context
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		message, err := findMessageByID(ctx, tx, id)
		if err != nil {
			return err
		} else if message == nil {
			return nil
		}
		if message.UserID != user.UserID && !user.isAdmin() {
			return session.ForbiddenError(ctx)
		}
		query := "DELETE FROM messages WHERE message_id=$1"
		if message.ParentID == message.MessageID {
			query = "DELETE FROM messages WHERE parent_id=$1"
		}
		_, err = tx.ExecContext(ctx, query, id)
		return err
	})
	if err != nil {
		if sessionErr, ok := err.(session.Error); ok {
			return sessionErr
		}
		return session.TransactionError(ctx, err)
	}
	return nil
}
