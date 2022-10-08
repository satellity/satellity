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

	link := "https://techcrunch.com/tag/bitcoin/feed/"
	resp, _ := client.Get(link)

	var feed feeds.Common
	d := xml.NewDecoder(resp.Body)
	d.Strict = false
	err := d.Decode(&feed)
	assert.Nil(err)
	log.Println(feed.Date())
	for _, entry := range feed.Channel.Entries {
		log.Printf("%#v", entry)
		break
	}
}
