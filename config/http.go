package config

import (
	"net/http"
)

// HTTPClient defines API for sending http requests.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
