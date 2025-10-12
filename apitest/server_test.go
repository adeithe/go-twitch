package apitest_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
)

func TestAPITest(t *testing.T) {
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
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, chatters, res.Total)
					require.ElementsMatch(t, res.Chatters, chatters)
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
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, chatters, res.Total)
					require.ElementsMatch(t, res.Chatters, chatters)
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
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "2345").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, chatters, res.Total)
					require.ElementsMatch(t, res.Chatters, chatters)
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
				})
			},
			func(api *api.Client, opts ...api.RequestOption) (func(t *testing.T), error) {
				res, err := api.Chat.Chatters.List("1234", "5678").Do(context.Background(), opts...)
				return func(t *testing.T) {
					require.Len(t, chatters, res.Total)
					require.ElementsMatch(t, res.Chatters, chatters)
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

			client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithClientSecret(secret))
			check, err := tt.fetch(client, api.WithBearerToken(token))
			require.NoError(t, err)
			require.Exactly(t, 1, endpoint.TimesCalled)
			require.Exactly(t, 1, endpoint.Successes)
			require.Exactly(t, 0, endpoint.Failures)
			check(t)
		})
	}
}
