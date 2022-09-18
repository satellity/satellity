package feeds

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
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

type Feed struct {
	Updated string   `xml:"updated"`
	Entries []*Entry `xml:"entry"`
}

func Release(ctx context.Context, link string) error {
	resp, err := client.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return fmt.Errorf("%s too many requests %d", link, resp.StatusCode)
	}
	if resp.StatusCode == 403 {
		return fmt.Errorf("%s forbidden %d", link, resp.StatusCode)
	}
	var feed Feed
	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return err
	}
	log.Printf("Feed %#v", feed)
	return nil
}

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
}
