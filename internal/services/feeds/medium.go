package feeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"satellity/internal/models"
	"time"
)

type EntryMedium struct {
	ID      string    `xml:"guid"`
	Updated time.Time `xml:"updated"`
	Link    string    `xml:"link"`
	Title   string    `xml:"title"`
	Content string    `xml:"encoded"`
	Author  string    `xml:"creator"`
}

type Medium struct {
	Channel struct {
		Updated string         `xml:"lastBuildDate"`
		Entries []*EntryMedium `xml:"item"`
	} `xml:"channel"`
}

func (m *Medium) Date() (time.Time, error) {
	updated := m.Channel.Updated
	return time.Parse(time.RFC1123, updated)
}

func parseMedium(r io.Reader) (*Channel, error) {
	var common Medium
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
			Title:     e.Title,
			ID:        e.ID,
			Link:      e.Link,
			Content:   e.Content,
			Author:    e.Author,
			UpdatedAt: e.Updated,
			Genre:     models.GIST_GENRE_DEFAULT,
			Cardinal:  true,
		}
		channel.Gists = append(channel.Gists, gist)
	}
	return channel, nil
}
