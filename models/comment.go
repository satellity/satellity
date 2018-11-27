package models

import "time"

const commentsDDL = `
CREATE TABLE IF NOT EXISTS comments (
	comment_id            VARCHAR(36) PRIMARY KEY,
	body                  TEXT NOT NULL,
  topic_id              VARCHAR(36) NOT NULL REFERENCES topics ON DELETE CASCADE,
	user_id               VARCHAR(36) NOT NULL,
	score                 INTEGER NOT NULL DEFAULT 0,
	created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX ON comments (topic_id, created_at);
CREATE INDEX ON comments (user_id, created_at);
CREATE INDEX ON comments (score, created_at);
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
}
