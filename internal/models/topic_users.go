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

CREATE UNIQUE INDEX IF NOT EXISTS topic_users_reversex ON topic_users(user_id, topic_id);
CREATE INDEX IF NOT EXISTS topic_users_likedx ON topic_users(topic_id, liked);
CREATE INDEX IF NOT EXISTS topic_users_bookmarkedx ON topic_users(topic_id, bookmarked);
`

//
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

	isNew bool
}

var topicUserColumns = []string{"topic_id", "user_id", "liked", "bookmarked", "created_at", "updated_at"}

func (tu *TopicUser) values() []interface{} {
	return []interface{}{tu.TopicID, tu.UserID, tu.Liked, tu.Bookmarked, tu.CreatedAt, tu.UpdatedAt}
}

// ActiondBy execute user action, like or bookmark a topic
func (topic *Topic) ActiondBy(mctx *Context, user *User, action string, state bool) (*Topic, error) {
	ctx := mctx.context
	if action != TopicUserActionLiked &&
		action != TopicUserActionBookmarked {
		return topic, session.BadDataError(ctx)
	}
	tu, err := readTopicUser(mctx, topic.TopicID, user.UserID)
	if err != nil {
		return topic, session.TransactionError(ctx, err)
	}
	if tu == nil {
		t := time.Now()
		tu = &TopicUser{
			TopicID:   topic.TopicID,
			UserID:    user.UserID,
			CreatedAt: t,
			UpdatedAt: t,

			isNew: true,
		}
	}
	err = mctx.database.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var lcount, bcount int64
		if action == TopicUserActionLiked {
			tu.Liked = state
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM topic_users WHERE topic_id=$1 AND liked=true", topic.TopicID).Scan(&lcount); err != nil {
				return err
			}
			if lcount > 0 {
				topic.LikesCount = lcount - 1
			}
			if state {
				topic.LikesCount = lcount + 1
			}
		}
		if action == TopicUserActionBookmarked {
			tu.Bookmarked = state
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM topic_users WHERE topic_id=$1 AND bookmarked=true", topic.TopicID).Scan(&bcount); err != nil {
				return err
			}
			if bcount > 0 {
				topic.BookmarksCount = bcount - 1
			}
			if state {
				topic.BookmarksCount = bcount + 1
			}
		}
		topic.IsLikedBy = tu.Liked
		topic.IsBookmarkedBy = tu.Bookmarked
		if _, err := tx.ExecContext(ctx, "UPDATE topics SET (likes_count,bookmarks_count)=($1,$2) WHERE topic_id=$3", topic.LikesCount, topic.BookmarksCount, topic.TopicID); err != nil {
			return err
		}
		if tu.isNew {
			params, positions := durable.PrepareColumnsWithValues(topicUserColumns)
			query := fmt.Sprintf("INSERT INTO topic_users (%s) VALUES (%s)", params, positions)
			if _, err := tx.ExecContext(ctx, query, tu.values()...); err != nil {
				return err
			}
			return nil
		}
		query := fmt.Sprintf("UPDATE topic_users SET %s=$1 WHERE topic_id=$2 AND user_id=$3", action)
		if _, err := mctx.database.ExecContext(ctx, query, state, tu.TopicID, tu.UserID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return topic, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func readTopicUser(mctx *Context, topicID, userID string) (*TopicUser, error) {
	ctx := mctx.context
	query := fmt.Sprintf("SELECT %s FROM topic_users WHERE topic_id=$1 AND user_id=$2", strings.Join(topicUserColumns, ","))
	row, err := mctx.database.QueryRowContext(ctx, query, topicID, userID)
	if err != nil {
		return nil, err
	}
	tu, err := topicUserFromRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tu, err
}

func fillTopicWithAction(mctx *Context, topic *Topic, user *User) error {
	if user == nil {
		return nil
	}
	tu, err := readTopicUser(mctx, topic.TopicID, user.UserID)
	if err != nil || tu == nil {
		return err
	}
	topic.IsLikedBy, topic.IsBookmarkedBy = tu.Liked, tu.Bookmarked
	return nil
}

func topicUserFromRow(row durable.Row) (*TopicUser, error) {
	var tu TopicUser
	err := row.Scan(&tu.TopicID, &tu.UserID, &tu.Liked, &tu.Bookmarked, &tu.CreatedAt, &tu.UpdatedAt)
	return &tu, err
}
