package feeds

import (
	"encoding/xml"
	"fmt"
	"io"
	"satellity/internal/models"
	"time"
)

type GithubEntry struct {
	ID      string    `xml:"id"`
	Updated time.Time `xml:"updated"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Title   string `xml:"title"`
	Content string `xml:"content"`
}

type Github struct {
	Updated time.Time      `xml:"updated"`
	Entries []*GithubEntry `xml:"entry"`
}

func parseGithub(r io.Reader) (*Channel, error) {
	var common Github
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
			UpdatedAt: e.Updated,
			Genre:     models.GIST_GENRE_RELEASE,
			Cardinal:  false,
		}
		channel.Gists = append(channel.Gists, gist)
	}
	return channel, nil
}
