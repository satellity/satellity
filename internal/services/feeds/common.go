package feeds

import (
	"context"
	"encoding/xml"
	"fmt"
	"satellity/internal/models"
	"sort"
	"time"
)

type EntryCommon struct {
	Id      string `xml:"guid"`
	Updated string `xml:"pubDate"`
	Link    string `xml:"link"`
	Title   string `xml:"title"`
	Content string `xml:"description"`
	Author  string `xml:"creator"`
}

type Common struct {
	Channel struct {
		Updated string         `xml:"pubDate"`
		Entries []*EntryCommon `xml:"item"`
	} `xml:"channel"`
}

func FetchCommon(ctx context.Context, s *models.Source) error {
	now := time.Now()
	resp, err := client.Get(s.Link)
	if err != nil {
		return fmt.Errorf("fetch error: %w", err)
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
	var common Common
	err = xml.NewDecoder(resp.Body).Decode(&common)
	if err != nil {
		return fmt.Errorf("xml decode error: %w", err)
	}

	feed := common.Channel
	updated, err := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", feed.Updated)
	if err != nil {
		return fmt.Errorf("time parse error: %w", err)
	}
	published := time.Time{}
	if updated.After(s.UpdatedAt) {
		entries := feed.Entries
		sort.Slice(entries, func(i, j int) bool {
			ati, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", entries[i].Updated)
			atj, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", entries[j].Updated)
			return ati.Before(atj)
		})
		for _, entry := range entries {
			at, err := time.Parse("Mon, 02 Jan 2006 15:04:05 +0000", entry.Updated)
			if err != nil {
				continue
			}
			if published.Before(at) {
				published = at
			}
			if at.Before(s.UpdatedAt) {
				continue
			}
			_, err = models.CreateGist(ctx, entry.Id, entry.Author, entry.Title, models.GIST_GENRE_DEFAULT, true, entry.Link, entry.Content, at, s)
			if err != nil {
				return fmt.Errorf("CreateGist error: %w", err)

			}
		}
	}
	return s.Update(ctx, "", "", "", 0, published, now)
}
