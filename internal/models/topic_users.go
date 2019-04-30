package models

import (
	"database/sql"
	"fmt"
	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"strings"
	"time"
)

// learn from https://github.com/discourse/discourse/blob/master/app/models/topic_user.rb
const topicUsersDDL = `
CREATE TABLE IF NOT EXISTS topic_users (
	topic_id              VARCHAR(36) NOT NULL REFERENCES topics ON DELETE CASCADE,
	user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	liked                 BOOL NOT NULL DEFAULT false,
	bookmarked            BOOL NOT NULL DEFAULT false,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (topic_id, user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS user_topicx ON topic_users (user_id, topic_id);
`

const (
	TopicUserActionLiked      = "liked"
	TopicUserActionBookmarked = "bookmarked"
)

// TopicUser contains the relationships between topic and user
type TopicUser struct {
	TopicID    string
	UserID     string
	Liked      bool
	Bookmarked bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

var topicUserColumns = []string{"topic_id", "user_id", "liked", "bookmarked", "created_at", "updated_at"}

func (tu *TopicUser) values() []interface{} {
	return []interface{}{tu.TopicID, tu.UserID, tu.Liked, tu.Bookmarked, tu.CreatedAt, tu.UpdatedAt}
}

// ActiondBy execute user action, like or bookmark a topic
func (topic *Topic) ActiondBy(mctx *Context, user *User, action string, state bool) error {
	ctx := mctx.context
	if action != TopicUserActionLiked &&
		action != TopicUserActionBookmarked {
		return session.BadDataError(ctx)
	}
	tu, err := readTopicUser(mctx, topic.TopicID, user.UserID)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	query := ""
	if tu == nil {
		params, positions := durable.PrepareColumnsWithValues(topicUserColumns)
		t := time.Now()
		tu = &TopicUser{
			TopicID:   topic.TopicID,
			UserID:    user.UserID,
			CreatedAt: t,
			UpdatedAt: t,
		}
		if action == TopicUserActionLiked {
			tu.Liked = state
		}
		if action == TopicUserActionBookmarked {
			tu.Bookmarked = state
		}
		query = fmt.Sprintf("INSERT INTO topic_users (%s) VALUES (%s)", params, positions)
		_, err = mctx.database.ExecContext(ctx, query, tu.values()...)
		if err != nil {
			return session.TransactionError(ctx, err)
		}
		return nil
	}
	query = fmt.Sprintf("UPDATE topic_users SET %s=$1 WHERE topic_id=$2 AND user_id=$3", action)
	_, err = mctx.database.ExecContext(ctx, query, state, tu.TopicID, tu.UserID)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func readTopicUser(mctx *Context, topicID, userID string) (*TopicUser, error) {
	ctx := mctx.context
	query := fmt.Sprintf("SELECT %s FROM topic_users WHERE topic_id=$1 AND user_id=$2", strings.Join(topicUserColumns, ","))
	row, err := mctx.database.QueryRowContext(ctx, query, topicID, userID)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	tu, err := topicUserFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tu, err
}

func topicUserFromRow(row durable.Row) (*TopicUser, error) {
	var tu TopicUser
	err := row.Scan(&tu.TopicID, &tu.UserID, &tu.Liked, &tu.Bookmarked, &tu.CreatedAt, &tu.UpdatedAt)
	return &tu, err
}
