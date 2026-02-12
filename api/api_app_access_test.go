package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestAuthorization_AppAccess(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.EnableHTTP2(), apitest.WithTLS())
	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), oauth2.HTTPClient, mock.Client()))
	defer cancel()

	clientID, clientSecret, err := mock.RegisterApplication()
	require.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		auth := api.AppAccess(clientID, clientSecret)
		token1, err := auth.Token(ctx)
		require.NoError(t, err)
		require.Empty(t, token1.RefreshToken)
		require.NotEmpty(t, token1.AccessToken)

		token2, err := auth.Token(ctx)
		require.NoError(t, err)
		require.Empty(t, token1.RefreshToken)
		require.NotEmpty(t, token2.AccessToken)
		require.Equal(t, token1.AccessToken, token2.AccessToken)
	})

	t.Run("Expired", func(t *testing.T) {
		auth := api.AppAccess(clientID, clientSecret)
		token1, err := auth.Token(ctx)
		require.NoError(t, err)
		require.Empty(t, token1.RefreshToken)
		require.NotEmpty(t, token1.AccessToken)
		token1.Expiry = time.Now().Add(-time.Hour)

		token2, err := auth.Token(ctx)
		require.NoError(t, err)
		require.Empty(t, token1.RefreshToken)
		require.NotEmpty(t, token2.AccessToken)
		require.NotEqual(t, token1.AccessToken, token2.AccessToken)
	})
}
