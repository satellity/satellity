package models

import (
	"context"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

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
func (topic *Topic) ActiondBy(ctx context.Context, user *User, action string, state bool) (*Topic, error) {
	if action != TopicUserActionLiked &&
		action != TopicUserActionBookmarked {
		return topic, session.BadDataError(ctx)
	}
	tu, err := readTopicUser(ctx, topic.TopicID, user.UserID)
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
	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var lcount, bcount int64
		if action == TopicUserActionLiked {
			tu.Liked = state
			if err := tx.QueryRow(ctx, "SELECT count(*) FROM topic_users WHERE topic_id=$1 AND liked=true", topic.TopicID).Scan(&lcount); err != nil {
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
			if err := tx.QueryRow(ctx, "SELECT count(*) FROM topic_users WHERE topic_id=$1 AND bookmarked=true", topic.TopicID).Scan(&bcount); err != nil {
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
		if _, err := tx.Exec(ctx, "UPDATE topics SET (likes_count,bookmarks_count)=($1,$2) WHERE topic_id=$3", topic.LikesCount, topic.BookmarksCount, topic.TopicID); err != nil {
			return err
		}
		if tu.isNew {
			rows := [][]interface{}{tu.values()}
			_, err := tx.CopyFrom(ctx, pgx.Identifier{"topic_users"}, topicUserColumns, pgx.CopyFromRows(rows))
			return err
		}
		query := fmt.Sprintf("UPDATE topic_users SET %s=$1 WHERE topic_id=$2 AND user_id=$3", action)
		if _, err := session.Database(ctx).Exec(ctx, query, state, tu.TopicID, tu.UserID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return topic, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func readTopicUser(ctx context.Context, topicID, userID string) (*TopicUser, error) {
	query := fmt.Sprintf("SELECT %s FROM topic_users WHERE topic_id=$1 AND user_id=$2", strings.Join(topicUserColumns, ","))
	row := session.Database(ctx).QueryRow(ctx, query, topicID, userID)
	tu, err := topicUserFromRow(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return tu, err
}

func fillTopicWithAction(ctx context.Context, topic *Topic, user *User) error {
	if user == nil {
		return nil
	}
	tu, err := readTopicUser(ctx, topic.TopicID, user.UserID)
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
