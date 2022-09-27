package feeds

import (
	"context"
	"encoding/xml"
	"fmt"
	"satellity/internal/models"
	"sort"
	"time"
)

type EntryMirror struct {
	Id      string    `xml:"id"`
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

func FetchMirror(ctx context.Context, s *models.Source) error {
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
	var feed Mirror
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
			_, err = models.CreateGist(ctx, entry.Id, entry.Author.Name, entry.Title, models.GIST_GENRE_DEFAULT, true, entry.Link.Href, entry.Content, entry.Updated, s)
			if err != nil {
				return err
			}
		}
	}
	return s.Update(ctx, "", "", "", 0, published, now)
}
