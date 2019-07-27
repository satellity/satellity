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
	created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS messages_groupx ON messages (group_id);
`

// Message represent the struct of a message
type Message struct {
	MessageID string
	Body      string
	GroupID   string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User
}

var messageColumns = []string{"message_id", "body", "group_id", "user_id", "created_at", "updated_at"}

func (m *Message) values() []interface{} {
	return []interface{}{m.MessageID, m.Body, m.GroupID, m.UserID, m.CreatedAt, m.UpdatedAt}
}

func messageFromRow(row durable.Row) (*Message, error) {
	var m Message
	err := row.Scan(&m.MessageID, &m.Body, &m.GroupID, &m.UserID, &m.CreatedAt, &m.UpdatedAt)
	return &m, err
}

// CreateMessage create a message
func (u *User) CreateMessage(mctx *Context, groupID, body string) (*Message, error) {
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

		t := time.Now()
		message = &Message{
			MessageID: uuid.Must(uuid.NewV4()).String(),
			Body:      body,
			GroupID:   groupID,
			UserID:    u.UserID,
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
	limit := 64
	if !g.IsMember {
		offset = time.Now()
		limit = 8
	}

	var messages []*Message
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM messages WHERE created_at<$1 ORDER BY created_at DESC LIMIT $2", strings.Join(messageColumns, ","))
		rows, err := tx.QueryContext(ctx, query, offset, limit)
		if err != nil {
			return err
		}
		defer rows.Close()

		userIds := make([]string, 0)
		for rows.Next() {
			msg, err := messageFromRow(rows)
			if err != nil {
				return err
			}
			userIds = append(userIds, msg.UserID)
			messages = append(messages, msg)
		}
		if err := rows.Err(); err != nil {
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
			return session.ForbiddenError(ctx)
		}
		if message.UserID != user.UserID && !user.isAdmin() {
			return session.ForbiddenError(ctx)
		}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM messages WHERE message_id=$1"), id)
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
