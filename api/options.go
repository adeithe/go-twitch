package api

import (
	"fmt"
	"net/http"
)

// ClientOption is a function that modifies a Client.
type ClientOption func(*Client)

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

// WithDefaultAuthorization sets the default authorization to use for API requests, usually an [App Access Token].
//
// [App Access Token]: https://dev.twitch.tv/docs/authentication/#app-access-tokens
func WithDefaultAuthorization(authorization Authorization) ClientOption {
	return func(c *Client) {
		c.authorization = authorization
	}
}

// WithHTTPClient sets the HTTP client to use for API requests.
func WithHTTPClient(client HTTPClient) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// RequestOption is a function that modifies an HTTP request.
type RequestOption func(*http.Request) error

// WithAuthorization sets the authorization to use for the API request, usually a [User Access Token].
//
// [User Access Token]: https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#authorization-code-grant-flow
func WithAuthorization(authorization Authorization) RequestOption {
	return func(req *http.Request) error {
		token, err := authorization.TokenSource(req.Context()).Token()
		if err != nil {
			return err
		}
		token.SetAuthHeader(req)
		return nil
	}
}

// WithBearerToken sets the bearer token to use for API requests.
func WithBearerToken(token string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}
}

// SetQueryParameter sets a query parameter on the request, replacing any existing values.
func SetQueryParameter[T any](key string, value T) RequestOption {
	return func(r *http.Request) error {
		q := r.URL.Query()
		q.Set(key, fmt.Sprint(value))
		r.URL.RawQuery = q.Encode()
		return nil
	}
}

// AddQueryParameter adds a query parameter to the request without replacing any existing values.
func AddQueryParameter[T any](key string, value T) RequestOption {
	return func(r *http.Request) error {
		q := r.URL.Query()
		q.Add(key, fmt.Sprint(value))
		r.URL.RawQuery = q.Encode()
		return nil
	}
}

// SetHeader sets a header on the request, replacing any existing values.
func SetHeader[T any](key string, value T) RequestOption {
	return func(r *http.Request) error {
		r.Header.Set(key, fmt.Sprint(value))
		return nil
	}
}

// AddHeader adds a header to the request without replacing any existing values.
func AddHeader[T any](key string, value T) RequestOption {
	return func(r *http.Request) error {
		r.Header.Add(key, fmt.Sprint(value))
		return nil
	}
}
