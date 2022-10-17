package feeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"satellity/internal/models"
	"time"
)

type Gist struct {
	Title     string
	ID        string
	Link      string
	Content   string
	Author    string
	UpdatedAt time.Time
	Genre     string
	Cardinal  bool
}

type Channel struct {
	UpdatedAt time.Time
	Gists     []*Gist
}

func parseCommon(r io.Reader) (*Channel, error) {
	var common Common
	d := xml.NewDecoder(r)
	d.Strict = false

	err := d.Decode(&common)
	if err != nil {
		return nil, fmt.Errorf("xml decode error: %w", err)
	}

	cha := common.Channel
	channel := &Channel{}
	channel.UpdatedAt, _ = common.Date()

	for _, e := range cha.Entries {
		gist := &Gist{
			Title:    e.Title,
			ID:       e.ID,
			Link:     e.Link,
			Content:  e.Content,
			Author:   e.Author,
			Genre:    models.GIST_GENRE_DEFAULT,
			Cardinal: true,
		}
		gist.UpdatedAt, _ = e.Date()
		channel.Gists = append(channel.Gists, gist)
	}

	return channel, nil
}
