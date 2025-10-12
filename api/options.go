package api

import (
	"fmt"
	"net/http"
)

// ClientOption is a function that modifies a Client.
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

// RequestOption is a function that modifies an HTTP request.
type RequestOption func(*http.Request)

// WithBearerToken sets the bearer token to use for API requests.
func WithBearerToken(token string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// SetQueryParameter sets a query parameter on the request, replacing any existing values.
func SetQueryParameter[T any](key string, value T) RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Set(key, fmt.Sprint(value))
		r.URL.RawQuery = q.Encode()
	}
}

// AddQueryParameter adds a query parameter to the request without replacing any existing values.
func AddQueryParameter[T any](key string, value T) RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Add(key, fmt.Sprint(value))
		r.URL.RawQuery = q.Encode()
	}
}

// SetHeader sets a header on the request, replacing any existing values.
func SetHeader[T any](key string, value T) RequestOption {
	return func(r *http.Request) {
		r.Header.Set(key, fmt.Sprint(value))
	}
}

// AddHeader adds a header to the request without replacing any existing values.
func AddHeader[T any](key string, value T) RequestOption {
	return func(r *http.Request) {
		r.Header.Add(key, fmt.Sprint(value))
	}
}
