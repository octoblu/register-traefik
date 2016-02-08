package healthchecker

import (
	"net/http"
	"time"
)

// Healthy returns true if the HTTP response to
// a GET request to the URI is a 200, false
// otherwise
func Healthy(uri string) bool {
	client := &http.Client{Timeout: time.Second * 1}
	response, err := client.Get(uri)

	if err != nil {
		return false
	}

	return response.StatusCode == 200
}
