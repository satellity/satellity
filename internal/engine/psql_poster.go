package engine

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"godiscourse/internal/durable"
	"godiscourse/internal/model"
	"godiscourse/internal/session"

	"github.com/gofrs/uuid"
)

type Psql struct {
	db *durable.Database
}

func NewPsql(db *durable.Database) *Psql {
	return &Psql{db: db}
}

func (p *Psql) UpdateUser(ctx context.Context, current *model.User, new *model.UserInfo) error {
	nickname, biography := strings.TrimSpace(new.Nickname), strings.TrimSpace(new.Biography)
	if len(nickname) == 0 && len(biography) == 0 {
		return nil
	}
	if nickname != "" {
		current.Nickname = nickname
	}
	if biography != "" {
		current.Biography = biography
	}
	current.UpdatedAt = time.Now()
	cols, params := durable.PrepareColumnsWithValues([]string{"nickname", "biography", "updated_at"})
	_, err := p.db.ExecContext(ctx, fmt.Sprintf("UPDATE users SET (%s)=(%s) WHERE user_id='%s'", cols, params, current.UserID), current.Nickname, current.Biography, current.UpdatedAt)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func (p *Psql) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user *model.User
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		user, err = model.FindUserByID(ctx, tx, id)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return user, nil
}

func (p *Psql) GetCategoryByID(ctx context.Context, id string) (*model.Category, error) {
	var category *model.Category
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		category, err = model.FindCategory(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return category, nil
}

func (p *Psql) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		categories, err = model.ReadCategories(ctx, tx)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return categories, nil
}

func (p *Psql) CreateTopic(ctx context.Context, userID string, t *model.TopicInfo) (*model.Topic, error) {
	title, body := strings.TrimSpace(t.Title), strings.TrimSpace(t.Body)
	if len(title) < minTitleSize {
		return nil, session.BadDataError(ctx)
	}

	now := time.Now()
	topic := &model.Topic{
		TopicID:   uuid.Must(uuid.NewV4()).String(),
		Title:     title,
		Body:      body,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	var err error
	topic.ShortID, err = model.GenerateShortID("topics", now)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}

	err = p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		category, err := p.GetCategoryByID(ctx, t.CategoryID)
		if err != nil {
			return err
		}
		if category == nil {
			return session.BadDataError(ctx)
		}
		topic.CategoryID = category.CategoryID
		category.LastTopicID = sql.NullString{String: topic.TopicID, Valid: true}
		count, err := topicsCountByCategory(ctx, tx, category.CategoryID)
		if err != nil {
			return err
		}
		category.TopicsCount, category.UpdatedAt = count+1, time.Now()
		cols, params := durable.PrepareColumnsWithValues(model.TopicColumns)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO topics(%s) VALUES (%s)", cols, params), topic.Values()...)
		if err != nil {
			return err
		}
		ccols, cparams := durable.PrepareColumnsWithValues([]string{"last_topic_id", "topics_count", "updated_at"})
		cvals := []interface{}{category.LastTopicID, category.TopicsCount, category.UpdatedAt}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", ccols, cparams, category.CategoryID), cvals...)
		if err != nil {
			return err
		}
		// _, err = upsertStatistic(ctx, tx, "topics")
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

// dispersalCategory update category's info, e.g.: LastTopicID, TopicsCount
func (p *Psql) dispersalCategory(ctx context.Context, id string) (*model.Category, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	var result *model.Category
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		result, err = model.FindCategory(ctx, tx, id)
		if err != nil {
			return err
		} else if result == nil {
			return session.NotFoundError(ctx)
		}
		topic, err := model.LastTopic(ctx, result.CategoryID, tx)
		if err != nil {
			return err
		}
		var lastTopicID = sql.NullString{String: "", Valid: false}
		if topic != nil {
			lastTopicID = sql.NullString{String: topic.TopicID, Valid: true}
		}
		if result.LastTopicID.String != lastTopicID.String {
			result.LastTopicID = lastTopicID
		}
		result.TopicsCount = 0
		if result.LastTopicID.Valid {
			count, err := topicsCountByCategory(ctx, tx, result.CategoryID)
			if err != nil {
				return err
			}
			result.TopicsCount = count
		}
		result.UpdatedAt = time.Now()
		cols, params := durable.PrepareColumnsWithValues([]string{"last_topic_id", "topics_count", "updated_at"})
		vals := []interface{}{result.LastTopicID, result.TopicsCount, result.UpdatedAt}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE categories SET (%s)=(%s) WHERE category_id='%s'", cols, params, result.CategoryID), vals...)
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

func (p *Psql) UpdateTopic(ctx context.Context, id string, t *model.TopicInfo) (*model.Topic, error) {
	title, body := strings.TrimSpace(t.Title), strings.TrimSpace(t.Body)
	if title != "" && len(title) < minTitleSize {
		return nil, session.BadDataError(ctx)
	}

	var topic *model.Topic
	var prevCategoryID string
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		topic, err = model.FindTopic(ctx, tx, id)
		if err != nil {
			return err
		} else if topic == nil {
			return nil
		}
		// todo: move to level up
		// } else if topic.UserID != user.UserID && !user.IsAdmin() {
		// 	return session.AuthorizationError(ctx)
		// }
		if title != "" {
			topic.Title = title
		}
		topic.Body = body
		if t.CategoryID != "" && topic.CategoryID != t.CategoryID {
			prevCategoryID = topic.CategoryID
			category, err := model.FindCategory(ctx, tx, t.CategoryID)
			if err != nil {
				return err
			} else if category == nil {
				return session.BadDataError(ctx)
			}
			topic.CategoryID = category.CategoryID
		}
		cols, params := durable.PrepareColumnsWithValues([]string{"title", "body", "category_id"})
		vals := []interface{}{topic.Title, topic.Body, topic.CategoryID}
		_, err = tx.ExecContext(ctx, fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id='%s'", cols, params, topic.TopicID), vals...)
		return err
	})
	if err != nil {
		if _, ok := err.(session.Error); ok {
			return nil, err
		}
		return nil, session.TransactionError(ctx, err)
	}
	if topic == nil {
		return nil, session.NotFoundError(ctx)
	}
	if prevCategoryID != "" {
		// go t.dispersalCategory(ctx, prevCategoryID)
		// go t.dispersalCategory(ctx, topic.CategoryID)
	}
	return topic, nil
}

// todo: rewrite with join
func (p *Psql) GetTopicByID(ctx context.Context, id string) (*model.Topic, error) {
	var topic *model.Topic
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		topic, err = model.FindTopic(ctx, tx, id)
		if err != nil {
			return err
		}
		if topic == nil {
			subs := strings.Split(id, "-")
			if len(subs) < 1 || len(subs[0]) <= 5 {
				return nil
			}
			id = subs[0]
			topic, err = model.FindTopicByShortID(ctx, tx, id)
			if topic == nil || err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func (p *Psql) GetTopicByUserID(ctx context.Context, userID string, offset time.Time) ([]*model.Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*model.Topic
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		// todo: join query
		query := fmt.Sprintf("SELECT %s FROM topics WHERE user_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(model.TopicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, userID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			topic, err := model.TopicFromRows(rows)
			if err != nil {
				return err
			}
			topics = append(topics, topic)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

func (p *Psql) GetTopicsByCategoryID(ctx context.Context, categoryID string, offset time.Time) ([]*model.Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*model.Topic
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		// todo: join query
		query := fmt.Sprintf("SELECT %s FROM topics WHERE category_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(model.TopicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, categoryID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

func (p *Psql) GetTopicsByOffset(ctx context.Context, offset time.Time) ([]*model.Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*model.Topic
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		// todo: join query
		query := fmt.Sprintf("SELECT %s FROM topics WHERE created_at<$1 ORDER BY created_at DESC LIMIT $2", strings.Join(model.TopicColumns, ","))
		rows, err := tx.QueryContext(ctx, query, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

func (p *Psql) CreateComment(ctx context.Context, c *model.CommentInfo) (*model.Comment, error) {
	body := strings.TrimSpace(c.Body)
	if len(body) < minCommentBodySize {
		return nil, session.BadDataError(ctx)
	}
	now := time.Now()
	result := &model.Comment{
		CommentID: uuid.Must(uuid.NewV4()).String(),
		Body:      body,
		UserID:    c.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t, err := p.GetTopicByID(ctx, c.TopicID)
	if err != nil {
		return nil, err
	} else if t == nil {
		return nil, session.NotFoundError(ctx)
	}

	err = p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		count, err := commentsCountByTopic(ctx, tx, c.TopicID)
		if err != nil {
			return err
		}
		t.CommentsCount = count + 1
		t.UpdatedAt = now
		result.TopicID = t.TopicID
		cols, params := durable.PrepareColumnsWithValues(model.CommentColumns)
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO comments (%s) VALUES (%s)", cols, params), result.Values()...)
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

func (p *Psql) UpdateComment(ctx context.Context, c *model.CommentInfo) (*model.Comment, error) {
	body := strings.TrimSpace(c.Body)
	if len(body) < minCommentBodySize {
		return nil, session.BadDataError(ctx)
	}
	var result *model.Comment
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		var err error
		result, err = model.FindComment(ctx, tx, c.CommentID)
		if err != nil {
			return err
		} else if result == nil {
			return session.NotFoundError(ctx)
		} else if result.UserID != c.UserID /*&& !user.isAdmin()*/ { // todo: move to level up
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

func (p *Psql) DeleteComment(ctx context.Context, id, uid string) error {
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		comment, err := model.FindComment(ctx, tx, id)
		if err != nil || comment == nil {
			return err
		}
		if /*!user.isAdmin() &&*/ uid != comment.UserID {
			return session.ForbiddenError(ctx)
		}
		topic, err := p.GetTopicByID(ctx, comment.TopicID)
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

func (p *Psql) GetCommentsByTopicID(ctx context.Context, topicID string, offset time.Time) ([]*model.Comment, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var result []*model.Comment
	err := p.db.RunInTransaction(ctx, func(tx *sql.Tx) error {
		// todo: join with user, category in query
		query := fmt.Sprintf("SELECT %s FROM comments WHERE topic_id=$1 AND created_at<$2 ORDER BY created_at DESC LIMIT $3", strings.Join(model.CommentColumns, ","))
		rows, err := tx.QueryContext(ctx, query, topicID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return result, nil
}
