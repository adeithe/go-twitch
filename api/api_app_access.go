package api

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

// Authorization provides methods to authenticate with the Twitch API.
type Authorization interface {
	Client(ctx context.Context) *http.Client
	TokenSource(ctx context.Context) oauth2.TokenSource
	Token(ctx context.Context) (*oauth2.Token, error)
}

type appAccess struct {
	config *clientcredentials.Config
	source oauth2.TokenSource
}

// AppAccess returns an Authorization that maintains a Twitch API app access token.
func AppAccess(clientID, clientSecret string) Authorization {
	return &appAccess{
		config: &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     twitch.Endpoint.TokenURL,
			AuthStyle:    oauth2.AuthStyleInParams,
		},
	}
}

// Client returns an HTTP client using the provided token.
func (a *appAccess) Client(ctx context.Context) *http.Client {
	return oauth2.NewClient(ctx, a.TokenSource(ctx))
}

// TokenSource returns an oauth2.TokenSource that maintains a Twitch API app access token.
func (a *appAccess) TokenSource(ctx context.Context) oauth2.TokenSource {
	if a.source != nil {
		return a.source
	}
	a.source = a.config.TokenSource(ctx)
	return a.source
}

// Token returns a valid Twitch API app access token, renewing it if necessary.
func (a *appAccess) Token(ctx context.Context) (*oauth2.Token, error) {
	return a.TokenSource(ctx).Token()
}
