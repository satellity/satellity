package models

import (
	"context"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	commentBodySizeLimit = 6
)

// Comment is struct for comment of topic
type Comment struct {
	CommentID string
	Body      string
	TopicID   string
	UserID    string
	Score     int
	CreatedAt time.Time
	UpdatedAt time.Time

	User *User
}

var commentColumns = []string{"comment_id", "body", "topic_id", "user_id", "score", "created_at", "updated_at"}

func (c *Comment) values() []interface{} {
	return []interface{}{c.CommentID, c.Body, c.TopicID, c.UserID, c.Score, c.CreatedAt, c.UpdatedAt}
}

func commentFromRows(row durable.Row) (*Comment, error) {
	var c Comment
	err := row.Scan(&c.CommentID, &c.Body, &c.TopicID, &c.UserID, &c.Score, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

// CreateComment create a new comment
func (user *User) CreateComment(ctx context.Context, topicID, body string) (*Comment, error) {
	body = strings.TrimSpace(body)
	if len(body) < commentBodySizeLimit {
		return nil, session.BadDataError(ctx)
	}
	t := time.Now()
	c := &Comment{
		CommentID: uuid.Must(uuid.NewV4()).String(),
		Body:      body,
		UserID:    user.UserID,
		CreatedAt: t,
		UpdatedAt: t,
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		topic, err := findTopic(ctx, tx, topicID)
		if err != nil {
			return err
		} else if topic == nil {
			return session.NotFoundError(ctx)
		}
		count, err := commentsCountByTopic(ctx, tx, topicID)
		if err != nil {
			return err
		}
		topic.CommentsCount = count + 1
		if topic.UpdatedAt.Add(time.Hour * 24 * 30).After(time.Now()) {
			topic.UpdatedAt = t
		}
		cols, posits := durable.PrepareColumnsWithParams([]string{"comments_count", "updated_at"})
		query := fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id='%s'", cols, posits, topic.TopicID)
		_, err = tx.Exec(ctx, query, topic.CommentsCount, topic.UpdatedAt)
		if err != nil {
			return err
		}
		c.TopicID = topic.TopicID
		_, err = upsertStatistic(ctx, tx, "comments")
		if err != nil {
			return err
		}
		rows := [][]interface{}{c.values()}
		_, err = tx.CopyFrom(ctx, pgx.Identifier{"comments"}, commentColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	c.User = user
	return c, nil
}

// UpdateComment update the comment by id
func (user *User) UpdateComment(ctx context.Context, id, body string) (*Comment, error) {
	body = strings.TrimSpace(body)
	if len(body) < commentBodySizeLimit {
		return nil, session.BadDataError(ctx)
	}
	var comment *Comment
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		comment, err = findComment(ctx, tx, id)
		if err != nil {
			return err
		} else if comment == nil {
			return session.NotFoundError(ctx)
		} else if comment.UserID != user.UserID && !user.isAdmin() {
			return session.ForbiddenError(ctx)
		}
		comment.Body = body
		comment.UpdatedAt = time.Now()
		cols, posits := durable.PrepareColumnsWithParams([]string{"body", "updated_at"})
		_, err = tx.Exec(ctx, fmt.Sprintf("UPDATE comments SET (%s)=(%s) WHERE comment_id='%s'", cols, posits, comment.CommentID), comment.Body, comment.UpdatedAt)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comment, nil
}

// ReadComments read comments by topicID, parameters: offset
func (topic *Topic) ReadComments(ctx context.Context, offset time.Time) ([]*Comment, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var comments []*Comment
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM comments WHERE topic_id=$1 AND created_at<$2 ORDER BY created_at LIMIT $3", strings.Join(commentColumns, ","))
		rows, err := tx.Query(ctx, query, topic.TopicID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		userIds := make([]string, 0)
		for rows.Next() {
			c, err := commentFromRows(rows)
			if err != nil {
				return err
			}
			userIds = append(userIds, c.UserID)
			comments = append(comments, c)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		userSet, err := readUserSet(ctx, tx, userIds)
		if err != nil {
			return err
		}
		for i, c := range comments {
			comments[i].User = userSet[c.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comments, nil
}

// ReadComments read comments by userID, parameters: offset
func (user *User) ReadComments(ctx context.Context, offset time.Time) ([]*Comment, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	rows, err := session.Database(ctx).Query(ctx, fmt.Sprintf("SELECT %s FROM comments WHERE user_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(commentColumns, ",")), user.UserID, offset, LIMIT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		c, err := commentFromRows(rows)
		if err != nil {
			return nil, err
		}
		c.User = user
		comments = append(comments, c)
	}
	if err := rows.Err(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comments, nil
}

// ReadComments read all comments, parameters: offset default time.Now()
func ReadComments(ctx context.Context, offset time.Time) ([]*Comment, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var comments []*Comment
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM comments WHERE updated_at<$1 ORDER BY updated_at DESC LIMIT $2", strings.Join(commentColumns, ","))
		rows, err := tx.Query(ctx, query, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		userIds := []string{}
		for rows.Next() {
			comment, err := commentFromRows(rows)
			if err != nil {
				return err
			}
			userIds = append(userIds, comment.UserID)
			comments = append(comments, comment)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		userSet, err := readUserSet(ctx, tx, userIds)
		if err != nil {
			return err
		}
		for i, comment := range comments {
			comments[i].User = userSet[comment.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comments, nil
}

// DeleteComment delete a comment by ID
func (user *User) DeleteComment(ctx context.Context, id string) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		comment, err := findComment(ctx, tx, id)
		if err != nil || comment == nil {
			return err
		}
		if !user.isAdmin() && user.UserID != comment.UserID {
			return session.ForbiddenError(ctx)
		}
		topic, err := findTopic(ctx, tx, comment.TopicID)
		if err != nil {
			return err
		} else if topic == nil {
			return session.BadDataError(ctx)
		}
		count, err := commentsCountByTopic(ctx, tx, comment.TopicID)
		if err != nil {
			return err
		}
		topic.CommentsCount = count - 1
		cols, posits := durable.PrepareColumnsWithParams([]string{"comments_count", "updated_at"})
		query := fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id='%s'", cols, posits, topic.TopicID)
		_, err = tx.Exec(ctx, query, topic.CommentsCount, topic.UpdatedAt)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "DELETE FROM comments WHERE comment_id=$1", comment.CommentID)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func findComment(ctx context.Context, tx pgx.Tx, id string) (*Comment, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM comments WHERE comment_id=$1", strings.Join(commentColumns, ",")), id)
	c, err := commentFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func commentsCountByTopic(ctx context.Context, tx pgx.Tx, id string) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, "SELECT count(*) FROM comments WHERE topic_id=$1", id).Scan(&count)
	return count, err
}

func commentsCount(ctx context.Context, tx pgx.Tx) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, "SELECT count(*) FROM comments").Scan(&count)
	return count, err
}
