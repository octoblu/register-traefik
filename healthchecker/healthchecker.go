package healthchecker

import "net/http"

// Healthy returns true if the HTTP response to
// a GET request to the URI is a 200, false
// otherwise
func Healthy(uri string) bool {
	response, err := http.Get(uri)

	if err != nil {
		return false
	}

	return response.StatusCode == 200
}
