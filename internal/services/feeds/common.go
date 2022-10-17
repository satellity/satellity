package feeds

import (
	"context"
	"fmt"
	"net/http"
	"satellity/internal/models"
	"sort"
	"time"
)

type EntryCommon struct {
	Title   string `xml:"title"`
	ID      string `xml:"guid"`
	Link    string `xml:"link"`
	Content string `xml:"description"`
	Author  string `xml:"creator"`
	Updated string `xml:"pubDate"`
}

type Common struct {
	Channel struct {
		Updated       string         `xml:"pubDate"`
		LastBuildDate string         `xml:"lastBuildDate"`
		Entries       []*EntryCommon `xml:"item"`
	} `xml:"channel"`
}

// time: "Mon, 02 Jan 2006 15:04:05 +0000"
func (c *Common) Date() (time.Time, error) {
	updated := c.Channel.Updated
	if updated == "" {
		updated = c.Channel.LastBuildDate
	}
	if updated == "" {
		return time.Now(), nil
	}
	t, err := time.Parse(time.RFC1123Z, updated)
	if err != nil {
		return time.Parse(time.RFC1123, updated)
	}
	return t, nil
}

func (e *EntryCommon) Date() (time.Time, error) {
	t, err := time.Parse(time.RFC1123Z, e.Updated)
	if err != nil {
		return time.Parse(time.RFC1123, e.Updated)
	}
	return t, nil
}

func FetchCommon(ctx context.Context, s *models.Source) error {
	now := time.Now()
	req, err := http.NewRequest("GET", s.Link, nil)
	if err != nil {
		return fmt.Errorf("new request error: %v", err)
	}
	req.Header.Set("user-agent", "Mozilla/5.0 AppleWebKit/537.36 Chrome/105.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch error: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 429:
		return fmt.Errorf("%s too many requests %d", s.Link, resp.StatusCode)
	case 403:
		return fmt.Errorf("%s forbidden %d", s.Link, resp.StatusCode)
	case 404:
		return s.Delete(ctx)
	}

	var channel *Channel
	switch s.Locality {
	case "github":
		channel, err = parseGithub(resp.Body)
	case "medium":
		channel, err = parseMedium(resp.Body)
	case "mirror":
		channel, err = parseMirror(resp.Body)
	default:
		channel, err = parseCommon(resp.Body)
	}

	if err != nil {
		return err
	}

	published := time.Time{}
	if channel.UpdatedAt.After(s.UpdatedAt) {
		gists := channel.Gists
		sort.Slice(gists, func(i, j int) bool {
			return gists[i].UpdatedAt.After(gists[j].UpdatedAt)
		})
		for _, gist := range gists {
			if published.Before(gist.UpdatedAt) {
				published = gist.UpdatedAt
			}
			if published.Before(s.UpdatedAt) {
				continue
			}
			_, err = models.CreateGist(ctx, gist.ID, gist.Author, gist.Title, gist.Genre, gist.Cardinal, gist.Link, gist.Content, gist.UpdatedAt, s)
			if err != nil {
				return fmt.Errorf("CreateGist error: %w", err)

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
