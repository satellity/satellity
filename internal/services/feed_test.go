package services

import (
	"encoding/xml"
	"log"
	"net/http"
	"satellity/internal/services/feeds"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	assert := assert.New(t)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// link := "https://www.coindesk.com/arc/outboundfeeds/rss/?outputType=xml"
	link := "https://dailycoin.com/feed/"
	req, err := http.NewRequest("GET", link, nil)
	assert.Nil(err)
	req.Header.Set("user-agent", "Mozilla/5.0 AppleWebKit/537.36 Chrome/105.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)

	var feed feeds.Common
	d := xml.NewDecoder(resp.Body)
	d.Strict = false
	err = d.Decode(&feed)
	assert.Nil(err)
	log.Println(feed.Date())
	log.Println(len(feed.Channel.Entries))
	for _, entry := range feed.Channel.Entries {
		log.Println("title>>", entry.Title)
		log.Println("link>>", entry.Link)
		log.Println("identity>>", entry.ID)
		log.Println("author>>", entry.Author)
		log.Println(entry.Date())
		break
	}
}
