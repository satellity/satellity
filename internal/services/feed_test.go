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

	link := "https://messari.io/rss"
	resp, _ := client.Get(link)

	var feed feeds.Common
	err := xml.NewDecoder(resp.Body).Decode(&feed)
	assert.Nil(err)
	log.Println(feed.Channel.Updated)
	for _, entry := range feed.Channel.Entries {
		log.Printf("%#v", entry)
		break
	}
}
