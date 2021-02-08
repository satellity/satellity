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
	"github.com/jackc/pgx/v4"
	hashids "github.com/speps/go-hashids"
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
	ShortID        string
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

var topicColumns = []string{"topic_id", "short_id", "title", "body", "topic_type", "comments_count", "bookmarks_count", "likes_count", "views_count", "category_id", "user_id", "score", "draft", "created_at", "updated_at"}

func (t *Topic) values() []interface{} {
	return []interface{}{t.TopicID, t.ShortID, t.Title, t.Body, t.TopicType, t.CommentsCount, t.BookmarksCount, t.LikesCount, t.ViewsCount, t.CategoryID, t.UserID, t.Score, t.Draft, t.CreatedAt, t.UpdatedAt}
}

func topicFromRows(row durable.Row) (*Topic, error) {
	var t Topic
	err := row.Scan(&t.TopicID, &t.ShortID, &t.Title, &t.Body, &t.TopicType, &t.CommentsCount, &t.BookmarksCount, &t.LikesCount, &t.ViewsCount, &t.CategoryID, &t.UserID, &t.Score, &t.Draft, &t.CreatedAt, &t.UpdatedAt)
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

	if !(typ == TopicTypePost || typ == TopicTypeLink) {
		return nil, session.BadDataError(ctx)
	}

	if typ == TopicTypeLink {
		if !strings.HasPrefix(body, "http") {
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
	var err error
	topic.ShortID, err = generateShortID("topics", t)
	if err != nil {
		return nil, session.ServerError(ctx, err)
	}

	err = session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		category, err := findCategory(ctx, tx, categoryID)
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
		if !topic.Draft {
			_, err = upsertStatistic(ctx, tx, "topics")
			if err != nil {
				return err
			}
		}
		category.TopicsCount, category.UpdatedAt = count+1, time.Now()
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
		// TODO
		emitToCategory(ctx, topic.CategoryID)
	}
	return topic, nil
}

// UpdateTopic update a Topic by ID
func (user *User) UpdateTopic(ctx context.Context, id, title, body, typ, categoryID string, draft bool) (*Topic, error) {
	title, body = strings.TrimSpace(title), strings.TrimSpace(body)
	if title != "" && len(title) < titleSizeLimit {
		return nil, session.BadDataError(ctx)
	}

	if !(typ == TopicTypePost || typ == TopicTypeLink) {
		return nil, session.BadDataError(ctx)
	}

	if typ == TopicTypeLink {
		if !strings.HasPrefix(body, "http") {
			return nil, session.BadDataError(ctx)
		}
	}

	var topic *Topic
	var prevCategoryID string
	var prevDraft bool
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		topic, err = findTopic(ctx, tx, id)
		if err != nil || topic == nil {
			return err
		} else if topic.UserID != user.UserID && !user.isAdmin() {
			return session.AuthorizationError(ctx)
		}
		topic.User = user
		if topic.UserID != user.UserID {
			u, err := findUserByID(ctx, tx, topic.UserID)
			if err != nil {
				return err
			}
			topic.User = u
		}
		prevDraft = topic.Draft
		if !topic.Draft && draft {
			return session.BadDataError(ctx)
		}
		topic.Draft = draft
		if title != "" {
			topic.Title = title
		}
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
			topic.Category = category
		}
		if typ != "" {
			topic.TopicType = typ
		}
		if prevDraft && !topic.Draft {
			_, err = upsertStatistic(ctx, tx, "topics")
			if err != nil {
				return err
			}
		}
		cols, params := durable.PrepareColumnsWithParams([]string{"title", "body", "category_id", "draft"})
		values := []interface{}{topic.Title, topic.Body, topic.CategoryID, topic.Draft}
		_, err = tx.Exec(ctx, fmt.Sprintf("UPDATE topics SET (%s)=(%s) WHERE topic_id='%s'", cols, params, topic.TopicID), values...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if topic == nil {
		return nil, session.NotFoundError(ctx)
	}
	err = fillTopicWithAction(ctx, topic, user)
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	if prevDraft && !topic.Draft {
		// TODO
		emitToCategory(ctx, prevCategoryID)
	} else if prevCategoryID != "" {
		emitToCategory(ctx, prevCategoryID)
		emitToCategory(ctx, topic.CategoryID)
	}
	return topic, nil
}

//ReadTopic read a topic by ID
func ReadTopic(ctx context.Context, id string) (*Topic, error) {
	var topic *Topic
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		topic, err = findTopic(ctx, tx, id)
		if err != nil {
			return err
		}
		if topic == nil {
			subs := strings.Split(id, "-")
			if len(subs) < 1 || len(subs[0]) <= 5 {
				return nil
			}
			id = subs[0]
			topic, err = findTopicByShortID(ctx, tx, id)
			if topic == nil || err != nil {
				return err
			}
		}
		user, err := findUserByID(ctx, tx, topic.UserID)
		if err != nil {
			return err
		}
		category, err := findCategory(ctx, tx, topic.CategoryID)
		if err != nil {
			return err
		}
		topic.User = user
		topic.Category = category
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func (user *User) DeleteTopic(ctx context.Context, id string) error {
	if !user.isAdmin() {
		return session.ForbiddenError(ctx)
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		topic, err := findTopic(ctx, tx, id)
		if err != nil || topic == nil {
			return err
		}

		_, err = tx.Exec(ctx, "DELETE FROM topics WHERE topic_id=$1", topic.TopicID)
		return err
	})
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

//DraftTopic read the draft topic
func (user *User) DraftTopic(ctx context.Context) (*Topic, error) {
	var topic *Topic
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		query := fmt.Sprintf("SELECT %s FROM topics WHERE user_id=$1 AND draft=true LIMIT 1", strings.Join(topicColumns, ","))
		row := tx.QueryRow(ctx, query, user.UserID)
		topic, err = topicFromRows(row)
		if err == pgx.ErrNoRows {
			topic = nil
			return nil
		}
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

//ReadTopicWithRelation read a topic with user's status like and bookmark
func ReadTopicWithRelation(ctx context.Context, id string, user *User) (*Topic, error) {
	topic, err := ReadTopic(ctx, id)
	if err != nil || topic == nil {
		return topic, err
	}
	err = fillTopicWithAction(ctx, topic, user)
	if err != nil {
		return topic, session.TransactionError(ctx, err)
	}
	if err := topic.incrTopicViewsCount(ctx); err != nil {
		session.ServerError(ctx, err)
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

// ReadTopicByShortID read a topic by Short ID
func ReadTopicByShortID(ctx context.Context, id string) (*Topic, error) {
	subs := strings.Split(id, "-")
	if len(subs) < 1 || len(subs[0]) <= 5 {
		return nil, nil
	}
	id = subs[0]
	var topic *Topic
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		topic, err = findTopicByShortID(ctx, tx, id)
		if topic == nil || err != nil {
			return err
		}
		user, err := findUserByID(ctx, tx, topic.UserID)
		if err != nil {
			return err
		}
		category, err := findCategory(ctx, tx, topic.CategoryID)
		if err != nil {
			return err
		}
		topic.User = user
		topic.Category = category
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return topic, nil
}

func findTopicByShortID(ctx context.Context, tx pgx.Tx, id string) (*Topic, error) {
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM topics WHERE short_id=$1", strings.Join(topicColumns, ",")), id)
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

func topicsCountByCategory(ctx context.Context, tx pgx.Tx, id string) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, "SELECT count(*) FROM topics WHERE category_id=$1 AND draft=false", id).Scan(&count)
	return count, err
}

func (topic *Topic) incrTopicViewsCount(ctx context.Context) error {
	topic.ViewsCount += 1
	_, err := session.Database(ctx).Exec(ctx, "UPDATE topics SET views_count=$1 WHERE topic_id=$2", topic.ViewsCount, topic.TopicID)
	if err != nil {
		return session.TransactionError(ctx, err)
	}
	return nil
}

func topicsCount(ctx context.Context, tx pgx.Tx) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, "SELECT count(*) FROM topics WHERE draft=false").Scan(&count)
	return count, err
}

func generateShortID(table string, t time.Time) (string, error) {
	hd := hashids.NewData()
	hd.MinLength = 5
	h, _ := hashids.NewWithData(hd)
	return h.EncodeInt64([]int64{t.UnixNano()})
}
