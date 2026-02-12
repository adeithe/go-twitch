package apitest_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

func TestMockAPI(t *testing.T) {
	chatters := []api.UserInfo{
		{UserID: "3456", UserLogin: "testuser", UserName: "TestUser"},
	}

	tests := []struct {
		name         string
		method, path string
		opts         []apitest.MockTwitchAPIOption
		endpoint     func(*apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint
		fetch        func(*api.Client, ...api.RequestOption) (func(t *testing.T), error)
	}{
		{
			"No TLS", http.MethodGet, api.EndpointChatGetChatters,
			[]apitest.MockTwitchAPIOption{},
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, api.EndpointChatGetChatters, &api.ResponseData[api.UserInfo]{
					Total: len(chatters),
					Data:  chatters,
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, chatters, res.Total)
					require.ElementsMatch(t, res.Data, chatters)
				}, err
			},
		},
		{
			"No TLS With HTTP2", http.MethodGet, api.EndpointChatGetChatters,
			[]apitest.MockTwitchAPIOption{apitest.WithTLS(), apitest.EnableHTTP2()},
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, api.EndpointChatGetChatters, &api.ResponseData[api.UserInfo]{
					Total: len(chatters),
					Data:  chatters,
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, chatters, res.Total)
					require.ElementsMatch(t, res.Data, chatters)
				}, err
			},
		},
		{
			"With TLS", http.MethodGet, api.EndpointChatGetChatters,
			[]apitest.MockTwitchAPIOption{apitest.WithTLS()},
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, api.EndpointChatGetChatters, &api.ResponseData[api.UserInfo]{
					Total: len(chatters),
					Data:  chatters,
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, chatters, res.Total)
					require.ElementsMatch(t, res.Data, chatters)
				}, err
			},
		},
		{
			"With TLS And HTTP2", http.MethodGet, api.EndpointChatGetChatters,
			[]apitest.MockTwitchAPIOption{apitest.WithTLS(), apitest.EnableHTTP2()},
			func(mock *apitest.MockTwitchAPI) *apitest.MockTwitchAPIEndpoint {
				return apitest.SetMockResponse(mock, http.MethodGet, api.EndpointChatGetChatters, &api.ResponseData[api.UserInfo]{
					Total: len(chatters),
					Data:  chatters,
				}, apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"))
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, chatters, res.Total)
					require.ElementsMatch(t, res.Data, chatters)
				}, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := apitest.NewMockAPI(t, tt.opts...)
			endpoint := tt.endpoint(mock)

			clientID, secret, err := mock.RegisterApplication()
			require.NoError(t, err)
			require.NotEmpty(t, clientID)
			require.NotEmpty(t, secret)

			token, err := mock.NewBearerToken(clientID)
			require.NoError(t, err)
			require.NotEmpty(t, token)

			oauthEndpoint := mock.OAuthTokenEndpoint()
			client := api.New(clientID, api.WithHTTPClient(mock.Client()))
			check, err := tt.fetch(client, api.WithBearerToken(token.AccessToken))
			require.NoError(t, err)
			require.Exactly(t, 0, oauthEndpoint.TimesCalled)
			require.Exactly(t, 0, oauthEndpoint.Successes)
			require.Exactly(t, 0, oauthEndpoint.Failures)
			require.Exactly(t, 1, endpoint.TimesCalled)
			require.Exactly(t, 1, endpoint.Successes)
			require.Exactly(t, 0, endpoint.Failures)
			check(t)
		})
	}
}

func TestMockAPI_ValidationFailure(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.WithTLS())
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters, apitest.RequireQueryParam("first"))

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	token, err := mock.NewBearerToken(clientID)
	require.NoError(t, err)

	oauthEndpoint := mock.OAuthTokenEndpoint()
	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background(), api.WithBearerToken(token.AccessToken))
	require.Error(t, err)
	require.Exactly(t, 0, oauthEndpoint.TimesCalled)
	require.Exactly(t, 0, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_HasCertificate(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.WithTLS())
	require.NotNil(t, mock.Certificate())
}

func TestMockAPI_NoMockResponse(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	clientID, clientSecret, err := mock.RegisterApplication()
	require.NoError(t, err)

	oauthEndpoint := mock.OAuthTokenEndpoint()
	authorization := api.AppAccess(clientID, clientSecret)
	client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithDefaultAuthorization(authorization))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 1, oauthEndpoint.TimesCalled)
	require.Exactly(t, 1, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
}

func TestMockAPI_MissingClientID(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	clientID, clientSecret, err := mock.RegisterApplication()
	require.NoError(t, err)

	oauthEndpoint := mock.OAuthTokenEndpoint()
	authorization := api.AppAccess(clientID, clientSecret)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)
	client := api.New("", api.WithHTTPClient(mock.Client()), api.WithDefaultAuthorization(authorization))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 1, oauthEndpoint.TimesCalled)
	require.Exactly(t, 1, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_InvalidClientID(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	clientID, clientSecret, err := mock.RegisterApplication()
	require.NoError(t, err)

	oauthEndpoint := mock.OAuthTokenEndpoint()
	authorization := api.AppAccess(clientID, clientSecret)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)
	client := api.New("testclient", api.WithHTTPClient(mock.Client()), api.WithDefaultAuthorization(authorization))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 1, oauthEndpoint.TimesCalled)
	require.Exactly(t, 1, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_MissingToken(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	oauthEndpoint := mock.OAuthTokenEndpoint()
	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 0, oauthEndpoint.TimesCalled)
	require.Exactly(t, 0, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_InvalidToken(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	oauthEndpoint := mock.OAuthTokenEndpoint()
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)
	client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithDefaultAuthorization(nil))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 0, oauthEndpoint.TimesCalled)
	require.Exactly(t, 0, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_MalformedToken(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	oauthEndpoint := mock.OAuthTokenEndpoint()
	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background(), api.SetHeader("Authorization", "Bearer"))
	require.Error(t, err)
	require.Exactly(t, 0, oauthEndpoint.TimesCalled)
	require.Exactly(t, 0, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_TokenMismatch(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)

	clientID, clientSecret, err := mock.RegisterApplication()
	require.NoError(t, err)

	_, _, err = mock.RegisterApplication()
	require.NoError(t, err)

	authorization := api.AppAccess(clientID, clientSecret)
	oauthEndpoint := mock.OAuthTokenEndpoint()
	client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithDefaultAuthorization(authorization))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background(), api.WithBearerToken("invalidtoken"))
	require.Error(t, err)
	require.Exactly(t, 0, oauthEndpoint.TimesCalled)
	require.Exactly(t, 0, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_OAuth2_EndpointNotFound(t *testing.T) {
	// The /oauth2/authorize endpoint is not supported by the apitest package.
	// This is because it's the user-facing endpoint for Twitch's OAuth2 flow and is not used by the API client.
	// This test ensures that requests to unsupported endpoints are properly handled by the mock server.
	req, err := http.NewRequest(http.MethodPost, "http://id.twitch.tv/oauth2/authorize", nil)
	require.NoError(t, err)

	mock := apitest.NewMockAPI(t)
	res, err := mock.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	bs, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Exactly(t, "404 Not Found", string(bs))
}

func TestMockAPI_OAuth2_InvalidClient(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://id.twitch.tv/oauth2/token", nil)
	require.NoError(t, err)

	mock := apitest.NewMockAPI(t)
	res, err := mock.Client().Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	bs, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"status": 400, "message": "invalid client"}`, string(bs))
}

func TestMockAPI_OAuth2_InvalidGrantType(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	ctx, cancel := context.WithCancel(context.WithValue(t.Context(), oauth2.HTTPClient, mock.Client()))
	defer cancel()

	oauthEndpoint := mock.OAuthTokenEndpoint()
	clientID, clientSecret, err := mock.RegisterApplication()
	require.NoError(t, err)
	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     twitch.Endpoint,
	}

	oauth2Config.Endpoint.AuthStyle = oauth2.AuthStyleInParams
	token, err := oauth2Config.PasswordCredentialsToken(ctx, "username", "password")
	require.Nil(t, token)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid grant type")
	require.Exactly(t, 1, oauthEndpoint.TimesCalled)
	require.Exactly(t, 0, oauthEndpoint.Successes)
	require.Exactly(t, 1, oauthEndpoint.Failures)
}
