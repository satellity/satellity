package external

import (
	"net/http"
	"time"
)

// HTTPClient is a client with Timeout (5 seconds).
func HTTPClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}
