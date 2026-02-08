package api_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
)

func TestAPI_Client(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.WithTLS())
	clientID, clientSecret, err := mock.RegisterApplication()
	require.NoError(t, err)

	endpoint := apitest.SetMockResponse(mock, http.MethodGet, api.EndpointChatGetChatters, &api.ResponseData[api.UserInfo]{
		Total: 1,
		Data: []api.UserInfo{
			{UserID: "3456", UserLogin: "testuser", UserName: "TestUser"},
		},
	}, apitest.RequireHeader("Authorization"), apitest.RequireQueryParam("broadcaster_id"), apitest.RequireQueryParam("moderator_id"))

	token, err := mock.NewBearerToken(clientID)
	require.NoError(t, err)

	oauthEndpoint := mock.OAuthTokenEndpoint()
	appAccess := api.AppAccess(clientID, clientSecret)
	client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithDefaultAuthorization(appAccess))
	res, err := client.Chat.Chatters.List("1234", "5678").Do(context.Background(), api.WithBearerToken(token.AccessToken))
	require.NoError(t, err)
	require.Len(t, res.Data, res.Total)
	require.Exactly(t, 0, oauthEndpoint.TimesCalled)
	require.Exactly(t, 0, oauthEndpoint.Successes)
	require.Exactly(t, 0, oauthEndpoint.Failures)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 1, endpoint.Successes)
	require.Exactly(t, 0, endpoint.Failures)
}

func BenchmarkAPI_Client(b *testing.B) {
	mock := apitest.NewMockAPI(b, apitest.WithTLS())
	clientID, _, err := mock.RegisterApplication()
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		_ = api.New(clientID)
	}
}

func TestAPI_Client_APIError(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.WithTLS())
	endpoint := apitest.SetMockResponse(mock, http.MethodGet, api.EndpointChatGetChatters, &api.ResponseData[api.UserInfo]{
		Total: 1,
		Data: []api.UserInfo{
			{UserID: "3456", UserLogin: "testuser", UserName: "TestUser"},
		},
	})

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.Error(t, err)
	require.Equal(t, http.StatusUnauthorized, api.CodeOf(err))
	require.Equal(t, "twitchapi: 401 Unauthorized - OAuth token is missing", err.Error())
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestAPI_Client_MissingQueryParam(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.WithTLS())
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters, apitest.RequireQueryParam("first"))

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	token, err := mock.NewBearerToken(clientID)
	require.NoError(t, err)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background(), api.WithBearerToken(token.AccessToken))
	require.Error(t, err)
	require.Equal(t, http.StatusBadRequest, api.CodeOf(err))
	require.Equal(t, "twitchapi: 400 Bad Request - missing query parameter: first", err.Error())
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Successes)
	require.Exactly(t, 1, endpoint.Failures)
}

func TestAPI_Options_SetHeader(t *testing.T) {
	r := &http.Request{Header: make(http.Header)}
	require.NoError(t, api.SetHeader("X-Custom-Header", "CustomValue")(r))
	require.Equal(t, "CustomValue", r.Header.Get("X-Custom-Header"))
}

func TestAPI_Options_AddHeader(t *testing.T) {
	r := &http.Request{Header: make(http.Header)}
	require.NoError(t, api.AddHeader("X-Custom-Header", "MyCustomValue")(r))
	require.ElementsMatch(t, []string{"MyCustomValue"}, r.Header.Values("X-Custom-Header"))
	require.NoError(t, api.AddHeader("X-Custom-Header", "MyOtherCustomValue")(r))
	require.ElementsMatch(t, []string{"MyCustomValue", "MyOtherCustomValue"}, r.Header.Values("X-Custom-Header"))
}

func TestAPI_ConduitError(t *testing.T) {
	err := api.ConduitError{
		ShardID: "1",
		Message: "The shard id is outside the conduit's range",
		Code:    "invalid_parameter",
	}
	require.Equal(t, "conduits: shard 1 invalid_parameter - The shard id is outside the conduit's range", err.Error())
}

func TestAPI_CharityCampaignAmount(t *testing.T) {
	data := api.CharityCampaignAmount{Value: 123456, Decimal: 2, Currency: "USD"}
	require.InDelta(t, 1234.56, data.Amount(), 0.00)
}

func TestAPI_VideoDuration(t *testing.T) {
	var duration api.VideoDuration
	require.NoError(t, json.Unmarshal([]byte("\"3h15m30s\""), &duration))
	require.Equal(t, time.Hour*3+time.Minute*15+time.Second*30, duration.AsDuration())
}

func TestAPI_WithWebhookTransport(t *testing.T) {
	bs := make([]byte, 22)
	if _, err := rand.Read(bs); err != nil {
		return
	}

	secret := fmt.Sprintf("%x", bs)[2:]
	callback := "http://localhost:8080/callback"
	transport := api.WithWebhookTransport(callback, secret)
	require.Equal(t, "webhook", transport.Method)
	require.Equal(t, callback, transport.Callback)
	require.Equal(t, secret, transport.Secret)
	require.Empty(t, transport.SessionID)
	require.Nil(t, transport.ConnectedAt)
	require.Nil(t, transport.DisconnectedAt)
}

func TestAPI_WithWebSocketTransport(t *testing.T) {
	bs := make([]byte, 10)
	if _, err := rand.Read(bs); err != nil {
		return
	}

	sessionID := fmt.Sprintf("%x", bs)[2:]
	transport := api.WithWebSocketTransport(sessionID)
	require.Equal(t, "websocket", transport.Method)
	require.Equal(t, sessionID, transport.SessionID)
	require.Empty(t, transport.Secret)
	require.Empty(t, transport.Callback)
	require.Nil(t, transport.ConnectedAt)
	require.Nil(t, transport.DisconnectedAt)
}

func TestAPI_WithChoice(t *testing.T) {
	choice := api.WithChoice("option1")
	require.Equal(t, "option1", choice.Title)
}

func TestAPI_WithPermanentBan(t *testing.T) {
	timeout := api.WithPermanentBan("1234", "You're banned!")
	require.Equal(t, "1234", timeout.UserID)
	require.Equal(t, "You're banned!", timeout.Reason)
	require.Nil(t, timeout.Duration)
}

func TestAPI_WithTimeout(t *testing.T) {
	timeout := api.WithTimeout("1234", "Stop breaking the rules", 30*time.Second)
	require.Equal(t, "1234", timeout.UserID)
	require.Equal(t, "Stop breaking the rules", timeout.Reason)
	require.Equal(t, 30, *timeout.Duration)
}

func TestAPI_CodeOf_Invalid(t *testing.T) {
	require.Equal(t, 500, api.CodeOf(&time.ParseError{}))
}
