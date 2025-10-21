package apitest_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
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

			client := api.New(clientID, api.WithHTTPClient(mock.Client()))
			check, err := tt.fetch(client, api.WithBearerToken(token))
			require.NoError(t, err)
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

	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background(), api.WithBearerToken(token))
	require.Error(t, err)
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
	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	token, err := mock.NewBearerToken(clientID)
	require.NoError(t, err)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithDefaultBearerToken(token))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
}

func TestMockAPI_MissingClientID(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)

	client := api.New("", api.WithHTTPClient(mock.Client()), api.WithDefaultBearerToken("invalidtoken"))
	_, err := client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_InvalidClientID(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)

	client := api.New("testclient", api.WithHTTPClient(mock.Client()), api.WithDefaultBearerToken("invalidtoken"))
	_, err := client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_MissingToken(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_InvalidToken(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithDefaultBearerToken("invalidtoken"))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_MalformedToken(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background(), api.SetHeader("Authorization", "Bearer"))
	require.Error(t, err)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestMockAPI_TokenMismatch(t *testing.T) {
	mock := apitest.NewMockAPI(t)
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters)

	clientID1, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	clientID2, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	token, err := mock.NewBearerToken(clientID2)
	require.NoError(t, err)

	client := api.New(clientID1, api.WithHTTPClient(mock.Client()), api.WithDefaultBearerToken(token))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}
