package api

import (
	"fmt"
	"net/http"
)

type ClientOption func(*Client)

// WithClientSecret sets the client secret to use for API requests.
func WithClientSecret(secret string) ClientOption {
	return func(c *Client) {
		c.clientSecret = secret
	}
}

// WithHTTPClient sets the HTTP client to use for API requests.
func WithHTTPClient(client HTTPClient) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

type RequestOption func(*http.Request)

// WithBearerToken sets the bearer token to use for API requests.
func WithBearerToken(token string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}
