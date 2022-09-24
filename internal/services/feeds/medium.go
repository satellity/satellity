package feeds

import (
	"context"
	"encoding/xml"
	"fmt"
	"satellity/internal/models"
	"time"
)

type EntryMedium struct {
	Id      string    `xml:"guid"`
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

func FetchMedium(ctx context.Context, s *models.Source) error {
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
	var medium Medium
	err = xml.NewDecoder(resp.Body).Decode(&medium)
	if err != nil {
		return err
	}

	feed := medium.Channel
	updated, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", feed.Updated)
	if err != nil {
		return err
	}
	if updated.After(s.UpdatedAt) {
		for _, entry := range feed.Entries {
			if entry.Updated.Before(s.UpdatedAt) {
				break
			}
			_, err = models.CreateGist(ctx, entry.Id, entry.Author, entry.Title, models.GIST_GENRE_DEFAULT, true, entry.Link, entry.Content, entry.Updated, s)
			if err != nil {
				return err
			}
		}
	}
	return s.Update(ctx, "", "", "", now)
}