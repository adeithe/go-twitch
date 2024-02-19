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

// WithDefaultBearerToken sets the bearer token to use for API requests.
//
// This can be considered dangerous as if a token is not provided per request, this will become the default token.
// Some developers may prefer to default to the App Access Token using this method. However, it is recommended to still
// use the WithBearerToken option for requests that require a token as this method always will fail if the App Access Token has expired.
func WithDefaultBearerToken(token string) ClientOption {
	return func(c *Client) {
		c.bearerToken = token
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

// SetQueryParameter sets a query parameter on the request, replacing any existing values.
func SetQueryParameter(key, value string) RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Set(key, value)
		r.URL.RawQuery = q.Encode()
	}
}

// AddQueryParameter adds a query parameter to the request without replacing any existing values.
func AddQueryParameter(key, value string) RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Add(key, value)
		r.URL.RawQuery = q.Encode()
	}
}
