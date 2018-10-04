package external

import (
	"net/http"
	"time"
)

func HttpClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}
