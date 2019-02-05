package models

import (
	"context"
	"strings"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/godiscourse/godiscourse/api/session"
	"github.com/satori/go.uuid"
)

const (
	minCommentBodySize = 6
)

const commentsDDL = `
CREATE TABLE IF NOT EXISTS comments (
	comment_id            VARCHAR(36) PRIMARY KEY,
	body                  TEXT NOT NULL,
  topic_id              VARCHAR(36) NOT NULL REFERENCES topics ON DELETE CASCADE,
	user_id               VARCHAR(36) NOT NULL REFERENCES users ON DELETE CASCADE,
	score                 INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON comments (topic_id, created_at);
CREATE INDEX ON comments (user_id, created_at);
CREATE INDEX ON comments (score DESC, created_at);
`

var commentColumns = []string{"comment_id", "body", "topic_id", "user_id", "score", "created_at", "updated_at"}

// Comment is struct for comment of topic
type Comment struct {
	CommentID string    `sql:"comment_id,pk"`
	Body      string    `sql:"body"`
	TopicID   string    `sql:"topic_id"`
	UserID    string    `sql:"user_id"`
	Score     int       `sql:"score,notnull"`
	CreatedAt time.Time `sql:"created_at"`
	UpdatedAt time.Time `sql:"updated_at"`

	User *User
}

// CreateComment create a new comment
func (user *User) CreateComment(ctx context.Context, topicID, body string) (*Comment, error) {
	body = strings.TrimSpace(body)
	if len(body) < minCommentBodySize {
		return nil, session.BadDataError(ctx)
	}
	c := &Comment{
		CommentID: uuid.Must(uuid.NewV4()).String(),
		Body:      body,
		UserID:    user.UserID,
	}
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
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
		c.TopicID = topic.TopicID
		tx.Update(topic)
		return tx.Insert(c)
	})

	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	c.User = user
	return c, nil
}

// UpdateComment update the comment by id
func (user *User) UpdateComment(ctx context.Context, id, body string) (*Comment, error) {
	body = strings.TrimSpace(body)
	if len(body) < minCommentBodySize {
		return nil, session.BadDataError(ctx)
	}
	var comment *Comment
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		var err error
		comment, err = findComment(ctx, tx, id)
		if err != nil {
			return err
		} else if comment == nil {
			return session.NotFoundError(ctx)
		} else if comment.UserID != user.UserID && !user.isAdmin() {
			return session.AuthorizationError(ctx)
		}
		comment.Body = body
		return tx.Update(comment)
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return comment, nil
}

// ReadComments read comments by topicID, parameters: offset
func (topic *Topic) ReadComments(ctx context.Context, offset time.Time) ([]*Comment, error) {
	var comments []*Comment
	if err := session.Database(ctx).Model(&comments).Relation("User").Where("comment.topic_id=? AND comment.created_at>?", topic.TopicID, offset).Order("comment.created_at").Limit(50).Select(); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comments, nil
}

// ReadComments read comments by userID, parameters: offset
func (user *User) ReadComments(ctx context.Context, offset time.Time) ([]*Comment, error) {
	if offset.IsZero() {
		offset = time.Now()
	}
	var comments []*Comment
	if _, err := session.Database(ctx).Query(&comments, "SELECT * FROM comments WHERE user_id=? AND created_at<? ORDER BY created_at DESC LIMIT 50", user.UserID, offset); err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return comments, nil
}

func (user *User) DeleteComment(ctx context.Context, id string) error {
	err := session.Database(ctx).RunInTransaction(func(tx *pg.Tx) error {
		comment, err := findComment(ctx, tx, id)
		if err != nil {
			return err
		} else if comment == nil {
			return nil
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
		tx.Update(topic)
		return tx.Delete(comment)
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return err
		}
		return session.TransactionError(ctx, err)
	}
	return nil
}

func findComment(ctx context.Context, tx *pg.Tx, id string) (*Comment, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	comment := &Comment{CommentID: id}
	if err := tx.Model(comment).Column(commentColumns...).WherePK().Select(); err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return comment, nil
}

func commentsCountByTopic(ctx context.Context, tx *pg.Tx, id string) (int, error) {
	return tx.Model(&Comment{}).Where("topic_id=?", id).Count()
}

// BeforeInsert hook insert
func (c *Comment) BeforeInsert(db orm.DB) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = c.CreatedAt
	return nil
}

// BeforeUpdate hook update
func (c *Comment) BeforeUpdate(db orm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}
