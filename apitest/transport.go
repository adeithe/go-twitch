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
	hostname := strings.ToLower(req.URL.Hostname())
	apiDomain := strings.EqualFold(hostname, "api.twitch.tv")
	oauthDomain := strings.EqualFold(hostname, "id.twitch.tv")
	if t.URL != nil && (apiDomain || oauthDomain) {
		u := *req.URL
		u.Scheme, u.Host = t.URL.Scheme, t.URL.Host
		u.Path = "/" + strings.ReplaceAll(hostname, ".", "/") + "/" + strings.TrimPrefix(u.Path, "/")
		req.URL = &u
	}
	return t.base.RoundTrip(req)
}
