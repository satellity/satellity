package external

import (
	"net/http"
	"time"
)

// HttpClient is a client with Timeout (5 seconds).
func HttpClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}
