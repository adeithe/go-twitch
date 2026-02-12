package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

func TestAuthorization_UserToken(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.EnableHTTP2(), apitest.WithTLS())
	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), oauth2.HTTPClient, mock.Client()))
	defer cancel()

	clientID, clientSecret, err := mock.RegisterApplication()
	require.NoError(t, err)
	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  "http://localhost:8080/callback",
	}

	t.Run("Valid", func(t *testing.T) {
		exchange, err := oauthConfig.Exchange(ctx, "gulfwdmys5lsm6qyz4xiz9q32l10", oauth2.AccessTypeOffline)
		require.NoError(t, err)

		auth := api.UserToken(oauthConfig, exchange)
		token, err := auth.Token(ctx)
		require.NoError(t, err)
		require.Equal(t, exchange.AccessToken, token.AccessToken)
	})

	t.Run("Expired", func(t *testing.T) {
		exchange, err := oauthConfig.Exchange(ctx, "gulfwdmys5lsm6qyz4xiz9q32l10", oauth2.AccessTypeOffline)
		require.NoError(t, err)

		exchange.Expiry = time.Now().Add(-time.Hour)
		auth := api.UserToken(oauthConfig, exchange)
		token, err := auth.Token(ctx)
		require.NoError(t, err)
		require.NotEqual(t, exchange.AccessToken, token.AccessToken)
	})
}
