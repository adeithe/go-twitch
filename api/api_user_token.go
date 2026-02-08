package api

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type userToken struct {
	config *oauth2.Config
	token  *oauth2.Token
	source oauth2.TokenSource
}

// UserToken returns an [Authorization] that maintains a Twitch API user access token.
func UserToken(config *oauth2.Config, token *oauth2.Token) Authorization {
	return &userToken{
		config: config,
		token:  token,
	}
}

// Client returns an HTTP client using the provided token.
func (a *userToken) Client(ctx context.Context) *http.Client {
	return oauth2.NewClient(ctx, a.TokenSource(ctx))
}

// TokenSource returns an oauth2.TokenSource that maintains a Twitch API app access token.
func (a *userToken) TokenSource(ctx context.Context) oauth2.TokenSource {
	if a.source != nil {
		return a.source
	}
	a.source = a.config.TokenSource(ctx, a.token)
	return a.source
}

// Token returns a valid Twitch API app access token, renewing it if necessary.
func (a *userToken) Token(ctx context.Context) (*oauth2.Token, error) {
	return a.TokenSource(ctx).Token()
}
