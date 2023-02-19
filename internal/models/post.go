package models

import (
	"satellity/internal/durable"
	"time"
)

type Post struct {
	PostID    string
	Title     string
	Link      string
	Meta      string
	OriginID  string
	CreatedAt time.Time
}

var postColumns = []string{"post_id", "title", "link", "meta", "origin_id", "created_at"}

func (p *Post) values() []any {
	return []any{p.PostID, p.Title, p.Link, p.Meta, p.OriginID, p.CreatedAt}
}

func postFromRows(row durable.Row) (*Post, error) {
	var p Post
	err := row.Scan(&p.PostID, &p.Title, &p.Link, &p.Meta, &p.OriginID, &p.CreatedAt)
	return &p, err
}
