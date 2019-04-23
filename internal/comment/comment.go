package comment

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"godiscourse/internal/durable"
	"godiscourse/internal/session"
	"godiscourse/internal/topic"
	"godiscourse/internal/user"

	"github.com/gofrs/uuid"
)

const (
	minCommentBodySize = 6
	LIMIT              = 50
)

type Model struct {
	CommentID string    `sql:"comment_id,pk"`
	Body      string    `sql:"body"`
	TopicID   string    `sql:"topic_id"`
	UserID    string    `sql:"user_id"`
	Score     int       `sql:"score,notnull"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`
}

type Params struct {
	CommentID string
	TopicID   string
	UserID    string
	Body      string
}

type CommentDatastore interface {
	Create(ctx context.Context, p *Params) (*Model, error)
	Update(ctx context.Context, p *Params) (*Model, error)
	Delete(ctx context.Context, id, uid string) error
	GetByTopicID(ctx context.Context, tid string, offset time.Time) ([]*Model, error)
	GetByUserID(ctx context.Context, uid string, offset time.Time) ([]*Model, error)
}

type Comment struct {
	db         *durable.Database
	userStore  *user.User
	topicStore *topic.Topic
}

func New(db *durable.Database, u *user.User, t *topic.Topic) *Comment {
	return &Comment{
		db:         db,
		userStore:  u,
		topicStore: t,
	}
}

func (c *Comment) Create(ctx context.Context, p *Params) (*Model, error) {
	body := strings.TrimSpace(p.Body)
	if len(body) < minCommentBodySize {
		return nil, session.BadDataError(ctx)
	}
	now := time.Now()
	result := &Model{
		CommentID: uuid.Must(uuid.NewV4()).String(),
		Body:      body,
		UserID:    p.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t, err := c.topicStore.GetByID(ctx, p.TopicID)
	if err != nil {
		return nil, err
	} else if t == nil {
		return nil, session.NotFoundError(ctx)
	}

	err = c.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		count, err := commentsCountByTopic(ctx, tx, p.TopicID)
		if err != nil {
			return err
		}
		t.CommentsCount = count + 1
		t.UpdatedAt = now
		result.TopicID = t.TopicID
		cols, params := durable.PrepareColumnsWithValues(commentColumns)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO comments (%s) VALUES (%s)", cols, params), result.values()...)
		if err != nil {
			return err
		}
		tcols, tparams := durable.PrepareColumnsWithValues([]string{"comments_count", "updated_at"})
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id='%s'", tcols, tparams, t.TopicID), t.CommentsCount, t.UpdatedAt)
		if err != nil {
			return err
		}
		// _, err = upsertStatistic(ctx, tx, "comments")
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return result, nil
}

func (c *Comment) Update(ctx context.Context, p *Params) (*Model, error) {
	body := strings.TrimSpace(p.Body)
	if len(body) < minCommentBodySize {
		return nil, session.BadDataError(ctx)
	}
	var result *Model
	err := c.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		result, err = findComment(ctx, tx, p.CommentID)
		if err != nil {
			return err
		} else if result == nil {
			return session.NotFoundError(ctx)
		} else if result.UserID != p.UserID /*&& !user.isAdmin()*/ {
			return session.ForbiddenError(ctx)
		}
		result.Body = body
		result.UpdatedAt = time.Now()
		cols, params := durable.PrepareColumnsWithValues([]string{"body", "updated_at"})
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE comments SET (%s)=(%s) WHERE comment_id='%s'", cols, params, result.CommentID), result.Body, result.UpdatedAt)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return result, nil
}

func (c *Comment) Delete(ctx context.Context, id, uid string) error {
	err := c.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		comment, err := findComment(ctx, tx, id)
		if err != nil || comment == nil {
			return err
		}
		if /*!user.isAdmin() &&*/ uid != comment.UserID {
			return session.ForbiddenError(ctx)
		}
		topic, err := c.topicStore.GetByID(ctx, comment.TopicID)
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
		topic.UpdatedAt = time.Now()
		cols, params := durable.PrepareColumnsWithValues([]string{"comments_count", "updated_at"})
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id='%s'", cols, params, topic.TopicID), topic.CommentsCount, topic.UpdatedAt)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "DELETE FROM comments WHERE comment_id=$1", comment.CommentID)
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

func (c *Comment) GetByTopicID(ctx context.Context, tid string, offset time.Time) ([]*Model, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var result []*Model
	err := c.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM comments WHERE topic_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(commentColumns, ","))
		rows, err := tx.QueryContext(ctx, query, tid, offset, LIMIT)
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
			result = append(result, comment)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		// ?
		// userSet, err := c.userStore.GetUserSet(ctx, tx, userIds)
		// if err != nil {
		// 	return err
		// }
		// for i, c := range comments {
		// 	comments[i].User = userSet[c.UserID]
		// }
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return result, nil
}

func (c *Comment) GetByUserID(ctx context.Context, uid string, offset time.Time) ([]*Model, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	rows, err := c.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM comments WHERE user_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(commentColumns, ",")), uid, offset, LIMIT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Model
	for rows.Next() {
		comment, err := commentFromRows(rows)
		if err != nil {
			return nil, err
		}
		comment.UserID = uid
		result = append(result, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return result, nil
}
