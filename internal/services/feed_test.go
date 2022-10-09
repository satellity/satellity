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

	link := "https://zycrypto.com/category/news/feed/"
	resp, _ := client.Get(link)

	var feed feeds.Common
	d := xml.NewDecoder(resp.Body)
	d.Strict = false
	err := d.Decode(&feed)
	assert.Nil(err)
	log.Println(feed.Date())
	log.Println(len(feed.Channel.Entries))
	for _, entry := range feed.Channel.Entries {
		log.Println("title>>", entry.Title)
		log.Println("link>>", entry.Link)
		log.Println(entry.Date())
		break
	}
}
