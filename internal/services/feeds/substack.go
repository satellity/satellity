package feeds

import (
	"context"
	"encoding/xml"
	"fmt"
	"satellity/internal/models"
	"sort"
	"time"
)

type EntrySubStack struct {
	Id      string `xml:"guid"`
	Updated string `xml:"pubDate"`
	Link    string `xml:"link"`
	Title   string `xml:"title"`
	Content string `xml:"description"`
	Author  string `xml:"creator"`
}

type SubStack struct {
	Channel struct {
		Updated string           `xml:"lastBuildDate"`
		Entries []*EntrySubStack `xml:"item"`
	} `xml:"channel"`
}

func FetchSubStack(ctx context.Context, s *models.Source) error {
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
	var substack SubStack
	err = xml.NewDecoder(resp.Body).Decode(&substack)
	if err != nil {
		return err
	}

	feed := substack.Channel
	updated, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", feed.Updated)
	if err != nil {
		return err
	}

	published := time.Time{}
	if updated.After(s.UpdatedAt) {
		entries := feed.Entries
		sort.Slice(entries, func(i, j int) bool {
			ati, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", entries[i].Updated)
			atj, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", entries[j].Updated)
			return ati.Before(atj)
		})
		for _, entry := range entries {
			at, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", entry.Updated)
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
				return err
			}
		}
	}
	return s.Update(ctx, "", "", "", 0, published, now)
}
