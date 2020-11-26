package graphql

import (
	"net/http"

	"github.com/Adeithe/go-twitch/api"
)

type httpTransport struct {
	client  *Client
	tripper http.RoundTripper
}

// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
//
// RoundTrip should not attempt to interpret the response. In particular, RoundTrip must return err == nil if it obtained a response, regardless of the response's HTTP status code. A non-nil err should be reserved for failure to obtain a response. Similarly, RoundTrip should not attempt to handle higher-level protocol details such as redirects, authentication, or cookies.
//
// RoundTrip should not modify the request, except for consuming and closing the Request's Body. RoundTrip may read fields of the request in a separate goroutine. Callers should not mutate or reuse the request until the Response's Body has been closed.
//
// RoundTrip must always close the body, including on errors, but depending on the implementation may do so in a separate goroutine even after RoundTrip returns. This means that callers wanting to reuse the body for subsequent requests must arrange to wait for the Close call before doing so.
//
// The Request's URL and Header fields must be initialized.
func (t httpTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if len(t.client.ID) < 1 {
		t.client.ID = api.Official.ID
	}
	r.Header.Set("Client-ID", t.client.ID)
	if len(t.client.bearer) > 0 {
		r.Header.Set("Authorization", "OAuth "+t.client.bearer)
	}
	return t.tripper.RoundTrip(r)
}
