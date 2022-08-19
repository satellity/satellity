package models

import (
	"context"
	"fmt"
	"net/url"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

// Topic related CONST
const (
	titleSizeLimit = 3
	LIMIT          = 30

	TopicTypePost = "POST"
	TopicTypeLink = "LINK"
)

// Topic is what use talking about
type Topic struct {
	TopicID        string
	Title          string
	Body           string
	TopicType      string
	CommentsCount  int64
	BookmarksCount int64
	LikesCount     int64
	ViewsCount     int64
	CategoryID     string
	UserID         string
	Score          int
	Draft          bool
	CreatedAt      time.Time
	UpdatedAt      time.Time

	IsLikedBy      bool
	IsBookmarkedBy bool
	User           *User
	Category       *Category
}

var topicColumns = []string{"topic_id", "title", "body", "topic_type", "comments_count", "bookmarks_count", "likes_count", "views_count", "category_id", "user_id", "score", "draft", "created_at", "updated_at"}

func (t *Topic) values() []interface{} {
	return []interface{}{t.TopicID, t.Title, t.Body, t.TopicType, t.CommentsCount, t.BookmarksCount, t.LikesCount, t.ViewsCount, t.CategoryID, t.UserID, t.Score, t.Draft, t.CreatedAt, t.UpdatedAt}
}

func topicFromRows(row durable.Row) (*Topic, error) {
	var t Topic
	err := row.Scan(&t.TopicID, &t.Title, &t.Body, &t.TopicType, &t.CommentsCount, &t.BookmarksCount, &t.LikesCount, &t.ViewsCount, &t.CategoryID, &t.UserID, &t.Score, &t.Draft, &t.CreatedAt, &t.UpdatedAt)
	return &t, err
}

//CreateTopic create a new Topic
func (user *User) CreateTopic(ctx context.Context, title, body, typ, categoryID string, draft bool) (*Topic, error) {
	if draft {
		t, err := user.DraftTopic(ctx)
		if err != nil {
			return nil, err
		}
		if t != nil {
			return nil, session.BadDataError(ctx)
		}
	}

	title, body = strings.TrimSpace(title), strings.TrimSpace(body)
	if len(title) < titleSizeLimit {
		return nil, session.BadDataError(ctx)
	}

	switch typ {
	case TopicTypePost, TopicTypeLink:
	default:
		return nil, session.BadDataError(ctx)
	}

	if typ == TopicTypeLink {
		_, err := url.ParseRequestURI(body)
		if err != nil {
			return nil, session.BadDataError(ctx)
		}
	}

	t := time.Now()
	topic := &Topic{
		TopicID:   uuid.Must(uuid.NewV4()).String(),
		Title:     title,
		Body:      body,
		TopicType: typ,
		UserID:    user.UserID,
		Draft:     draft,
		CreatedAt: t,
		UpdatedAt: t,
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		category, err := findCategory(ctx, tx, categoryID)
		if err != nil {
			return err
		}
		if category == nil {
			return session.BadDataError(ctx)
		}
		topic.CategoryID = category.CategoryID
		rows := [][]interface{}{
			topic.values(),
		}
		_, err = tx.CopyFrom(ctx, pgx.Identifier{"topics"}, topicColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if !topic.Draft {
		UpsertStatistic(ctx, StatisticTypeTopics)
		emitToCategory(ctx, topic.CategoryID)
	}
	return topic, nil
}

// UpdateTopic update a Topic by ID
func (user *User) UpdateTopic(ctx context.Context, id, title, body, typ, categoryID string, draft bool) (*Topic, error) {
	title, body = strings.TrimSpace(title), strings.TrimSpace(body)
	if len(title) < titleSizeLimit {
		return nil, session.BadDataError(ctx)
	}

	switch typ {
	case TopicTypePost, TopicTypeLink:
	default:
		return nil, session.BadDataError(ctx)
	}

	if typ == TopicTypeLink {
		_, err := url.ParseRequestURI(body)
		if err != nil {
			return nil, session.BadDataError(ctx)
		}
	}

	var topic *Topic
	var prevCategoryID string
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		topic, err = findTopic(ctx, tx, id)
		if err != nil || topic == nil {
			return err
		}
		if !topic.isPermit(user) {
			return session.ForbiddenError(ctx)
		}
		if draft && !topic.Draft {
			return session.ForbiddenError(ctx)
		}
		topic.Draft = draft

		topic.Title = title
		topic.Body = body
		if categoryID != "" && topic.CategoryID != categoryID {
			prevCategoryID = topic.CategoryID
			category, err := findCategory(ctx, tx, categoryID)
			if err != nil {
				return err
			} else if category == nil {
				return session.BadDataError(ctx)
			}
			topic.CategoryID = category.CategoryID
		}
		topic.TopicType = typ
		topic.UpdatedAt = time.Now()
		cols, params := durable.PrepareColumnsAndExpressions([]string{"title", "body", "category_id", "draft", "updated_at"}, 1)
		values := []interface{}{topic.TopicID, topic.Title, topic.Body, topic.CategoryID, topic.Draft, topic.UpdatedAt}
		_, err = tx.Exec(ctx, fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id=$1", cols, params), values...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if topic == nil {
		return nil, nil
	}
	if !topic.Draft {
		UpsertStatistic(ctx, StatisticTypeTopics)
		emitToCategory(ctx, topic.CategoryID)
		if prevCategoryID != "" {
			emitToCategory(ctx, prevCategoryID)
		}
	}
	return topic, nil
}

//ReadTopic read a topic by ID
func ReadTopic(ctx context.Context, id string) (*Topic, error) {
	var topic *Topic
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		topic, err = findTopic(ctx, tx, id)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

//ReadTopicWithRelation read a topic with user's status like and bookmark
func ReadTopicFull(ctx context.Context, id string, user *User) (*Topic, error) {
	topic, err := ReadTopic(ctx, id)
	if err != nil || topic == nil {
		return topic, err
	}
	err = topic.FillOut(ctx, user)
	if err != nil {
		return nil, err
	}
	topic.IncrViewsCount(ctx)
	return topic, nil
}

func (topic *Topic) Delete(ctx context.Context, user *User) error {
	if user == nil {
		return session.ForbiddenError(ctx)
	}
	if !(user.isAdmin() || (topic.Draft && topic.UserID == user.UserID)) {
		return session.ForbiddenError(ctx)
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		topic, err := findTopic(ctx, tx, topic.TopicID)
		if err != nil || topic == nil {
			return err
		}

		_, err = tx.Exec(ctx, "DELETE FROM topics WHERE topic_id=$1", topic.TopicID)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	UpsertStatistic(ctx, StatisticTypeTopics)
	return nil
}

//DraftTopic read the draft topic
func (user *User) DraftTopic(ctx context.Context) (*Topic, error) {
	var topic *Topic
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM topics WHERE user_id=$1 AND draft=true LIMIT 1", strings.Join(topicColumns, ","))
		row := tx.QueryRow(ctx, query, user.UserID)
		exist, err := topicFromRows(row)
		if err == pgx.ErrNoRows {
			return nil
		}
		topic = exist
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func findTopic(ctx context.Context, tx pgx.Tx, id string) (*Topic, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, nil
	}
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE topic_id=$1", strings.Join(topicColumns, ",")), id)
	t, err := topicFromRows(row)
	if pgx.ErrNoRows == err {
		return nil, nil
	}
	return t, err
}

// ReadTopics read all topics, parameters: offset default time.Now()
func ReadTopics(ctx context.Context, offset time.Time) ([]*Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*Topic
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		set, err := readCategorySet(ctx, tx)
		if err != nil {
			return err
		}

		query := fmt.Sprintf("SELECT %s FROM topics WHERE draft=false AND updated_at<$1 ORDER BY draft,updated_at DESC LIMIT $2", strings.Join(topicColumns, ","))
		rows, err := tx.Query(ctx, query, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		userIDs := []string{}
		for rows.Next() {
			topic, err := topicFromRows(rows)
			if err != nil {
				return err
			}
			userIDs = append(userIDs, topic.UserID)
			topic.Category = set[topic.CategoryID]
			topics = append(topics, topic)
		}
		if rows.Err() != nil {
			return rows.Err()
		}
		userSet, err := readUserSet(ctx, tx, userIDs)
		if err != nil {
			return err
		}
		for i, topic := range topics {
			topics[i].User = userSet[topic.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

// ReadTopics read user's topics, parameters: offset default time.Now()
func (user *User) ReadTopics(ctx context.Context, offset time.Time) ([]*Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*Topic
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		set, err := readCategorySet(ctx, tx)
		if err != nil {
			return err
		}
		query := fmt.Sprintf("SELECT %s FROM topics WHERE user_id=$1 AND draft=false AND created_at<$2 ORDER BY user_id,draft,created_at DESC LIMIT $3", strings.Join(topicColumns, ","))
		rows, err := tx.Query(ctx, query, user.UserID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			topic, err := topicFromRows(rows)
			if err != nil {
				return err
			}
			topic.User = user
			topic.Category = set[topic.CategoryID]
			topics = append(topics, topic)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

// ReadTopics read topics by CategoryID order by created_at DESC
func (category *Category) ReadTopics(ctx context.Context, offset time.Time) ([]*Topic, error) {
	if offset.IsZero() {
		offset = time.Now()
	}

	var topics []*Topic
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM topics WHERE category_id=$1 AND draft=false AND updated_at<$2 ORDER BY category_id,draft,updated_at DESC LIMIT $3", strings.Join(topicColumns, ","))
		rows, err := tx.Query(ctx, query, category.CategoryID, offset, LIMIT)
		if err != nil {
			return err
		}
		defer rows.Close()

		userIDs := []string{}
		for rows.Next() {
			topic, err := topicFromRows(rows)
			if err != nil {
				return err
			}
			userIDs = append(userIDs, topic.UserID)
			topic.Category = category
			topics = append(topics, topic)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		userSet, err := readUserSet(ctx, tx, userIDs)
		if err != nil {
			return err
		}
		for i, topic := range topics {
			topics[i].User = userSet[topic.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topics, nil
}

func (category *Category) latestTopic(ctx context.Context, tx pgx.Tx) (*Topic, error) {
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE category_id=$1 AND draft=false ORDER BY category_id,draft,updated_at DESC LIMIT 1", strings.Join(topicColumns, ",")), category.CategoryID)
	t, err := topicFromRows(row)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return t, err
}

func (topic *Topic) FillOut(ctx context.Context, user *User) error {
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		user, err := findUserByID(ctx, tx, topic.UserID)
		if err != nil {
			return err
		}
		topic.User = user
		category, err := findCategory(ctx, tx, topic.CategoryID)
		if err != nil {
			return err
		}
		topic.Category = category
		if user != nil {
			tu, err := findTopicUser(ctx, tx, topic.TopicID, user.UserID)
			if err != nil || tu == nil {
				return err
			}
			topic.IsLikedBy, topic.IsBookmarkedBy = tu.LikedAt.Valid, tu.BookmarkedAt.Valid
		}
		return nil
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func (topic *Topic) IncrViewsCount(ctx context.Context) error {
	topic.ViewsCount += 1
	_, err := session.Database(ctx).Exec(ctx, "UPDATE topics SET views_count=$1 WHERE topic_id=$2", topic.ViewsCount, topic.TopicID)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func (topic *Topic) isPermit(user *User) bool {
	if user == nil {
		return false
	}
	if user.isAdmin() {
		return true
	}
	return topic.UserID == user.UserID
}

func topicsCount(ctx context.Context, tx pgx.Tx) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, "SELECT count(*) FROM topics WHERE draft=false").Scan(&count)
	return count, err
}

func topicsCountByCategory(ctx context.Context, tx pgx.Tx, id string) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, "SELECT count(*) FROM topics WHERE category_id=$1 AND draft=false", id).Scan(&count)
	return count, err
}
