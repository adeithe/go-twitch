package apitest

import (
	"net/http"
	"net/url"
	"strings"
)

type mockTransport struct {
	base http.RoundTripper

	URL *url.URL
}

func newMockTransport(client *http.Client, url *url.URL) *mockTransport {
	return &mockTransport{base: client.Transport, URL: url}
}

// RoundTrip implements the http.RoundTripper interface for mockTransport.
func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.EqualFold(req.URL.Hostname(), "api.twitch.tv") && t.URL != nil {
		u := *req.URL
		u.Scheme, u.Host = t.URL.Scheme, t.URL.Host
		req.URL = &u
	}
	return t.base.RoundTrip(req)
}
