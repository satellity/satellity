package feeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"satellity/internal/models"
	"time"
)

type EntryMirror struct {
	ID      string    `xml:"id"`
	Updated time.Time `xml:"updated"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Title   string `xml:"title"`
	Content string `xml:"content"`
	Author  struct {
		Name string `xml:"name"`
	} `xml:"author"`
}

type Mirror struct {
	Updated time.Time      `xml:"updated"`
	Entries []*EntryMirror `xml:"entry"`
}

func parseMirror(r io.Reader) (*Channel, error) {
	var common Mirror
	d := xml.NewDecoder(r)
	d.Strict = false

	err := d.Decode(&common)
	if err != nil {
		return nil, fmt.Errorf("xml decode error: %w", err)
	}

	channel := &Channel{
		UpdatedAt: common.Updated,
	}
	for _, e := range common.Entries {
		gist := &Gist{
			Title:     e.Title,
			ID:        e.ID,
			Link:      e.Link.Href,
			Content:   e.Content,
			Author:    e.Author.Name,
			UpdatedAt: e.Updated,
			Genre:     models.GIST_GENRE_DEFAULT,
			Cardinal:  true,
		}
		channel.Gists = append(channel.Gists, gist)
	}
	return channel, nil
}
