package feeds

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"satellity/internal/models"
	"sort"
	"time"
)

type Entry struct {
	Id      string    `xml:"id"`
	Updated time.Time `xml:"updated"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Title   string `xml:"title"`
	Content string `xml:"content"`
}

type Github struct {
	Updated time.Time `xml:"updated"`
	Entries []*Entry  `xml:"entry"`
}

func FetchGithub(ctx context.Context, s *models.Source) error {
	now := time.Now()
	resp, err := client.Get(s.Link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return fmt.Errorf("%s too many requests %d", s.Link, resp.StatusCode)
	}
	if resp.StatusCode == 403 {
		return fmt.Errorf("%s forbidden %d", s.Link, resp.StatusCode)
	}
	if resp.StatusCode == 404 {
		return s.Delete(ctx)
	}
	var feed Github
	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return err
	}

	published := time.Time{}
	if feed.Updated.After(s.UpdatedAt) {
		entries := feed.Entries
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Updated.Before(entries[j].Updated)
		})
		for _, entry := range entries {
			if published.Before(entry.Updated) {
				published = entry.Updated
			}
			if entry.Updated.Before(s.UpdatedAt) {
				continue
			}
			_, err = models.CreateGist(ctx, entry.Id, "", entry.Title, models.GIST_GENRE_RELEASE, false, entry.Link.Href, entry.Content, entry.Updated, s)
			if err != nil {
				return err
			}
		}
	}
	return s.Update(ctx, "", "", "", 0, published, now)
}

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
}
