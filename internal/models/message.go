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

const messagesDDL = `
	message_id           VARCHAR(36) PRIMARY KEY,
	body                 TEXT NOT NULL,
	group_id             VARCHAR(36) NOT NULL REFERENCES groups ON DELETE CASCADE,
	user_id              VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
`

type Message struct {
	MessageID string
	Body      string
	GroupID   string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
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

func (u *User) CreateMessage(mctx *Context, groupID, body string) (*Group, error) {
	ctx := mctx.context
	body = strings.TrimSpace(body)
	if len(body) < 1 {
		return nil, session.BadDataError(ctx)
	}
	var group *Group
	err := mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		p, err := findParticipant(ctx, tx, groupID, u.UserID)
		if err != nil {
			return err
		} else if p == nil {
			return session.ForbiddenError(ctx)
		}

		t := time.Now()
		m := Message{
			MessageID: uuid.Must(uuid.NewV4()).String(),
			Body:      body,
			GroupID:   groupID,
			UserID:    u.UserID,
			CreatedAt: t,
			UpdatedAt: t,
		}

		columns, params := durable.PrepareColumnsWithValues(messageColumns)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO messages (%s) VALUES (%s)", columns, params), m.values()...)
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
