package api_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/adeithe/go-twitch/api"
	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
)

func TestAPI_Client(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.WithTLS())
	endpoint := apitest.SetMockResponse(mock, http.MethodGet, api.EndpointChatGetChatters, &api.ResponseData[api.UserInfo]{
		Total: 1,
		Data: []api.UserInfo{
			{UserID: "3456", UserLogin: "testuser", UserName: "TestUser"},
		},
	}, RequireQueryParam(t, "broadcaster_id"), RequireQueryParam(t, "moderator_id"))

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	token, err := mock.NewBearerToken(clientID)
	require.NoError(t, err)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()), api.WithDefaultBearerToken(token))
	res, err := client.Chat.Chatters.List("1234", "5678").Do(context.Background())
	require.NoError(t, err)
	require.Len(t, res.Data, res.Total)
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 0, endpoint.Failures)
	require.Exactly(t, 1, endpoint.Successes)
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
	require.Exactly(t, 1, endpoint.Failures)
	require.Exactly(t, 0, endpoint.Successes)
}

func TestAPI_Client_MissingQueryParam(t *testing.T) {
	mock := apitest.NewMockAPI(t, apitest.WithTLS())
	endpoint := apitest.SetMockValidator(mock, http.MethodGet, api.EndpointChatGetChatters, RequireQueryParam(t, "first"))

	clientID, _, err := mock.RegisterApplication()
	require.NoError(t, err)

	token, err := mock.NewBearerToken(clientID)
	require.NoError(t, err)

	client := api.New(clientID, api.WithHTTPClient(mock.Client()))
	_, err = client.Chat.Chatters.List("1234", "5678").Do(context.Background(), api.WithBearerToken(token))
	require.Error(t, err)
	require.Equal(t, http.StatusBadRequest, api.CodeOf(err))
	require.Equal(t, "twitchapi: 400 Bad Request - missing query parameter: first", err.Error())
	require.Exactly(t, 1, endpoint.TimesCalled)
	require.Exactly(t, 1, endpoint.Failures)
	require.Exactly(t, 0, endpoint.Successes)
}

func TestAPI_Options_SetHeader(t *testing.T) {
	r := &http.Request{Header: make(http.Header)}
	api.SetHeader("X-Custom-Header", "CustomValue")(r)
	require.Equal(t, "CustomValue", r.Header.Get("X-Custom-Header"))
}

func TestAPI_Options_AddHeader(t *testing.T) {
	r := &http.Request{Header: make(http.Header)}
	api.AddHeader("X-Custom-Header", "MyCustomValue")(r)
	require.ElementsMatch(t, []string{"MyCustomValue"}, r.Header.Values("X-Custom-Header"))
	api.AddHeader("X-Custom-Header", "MyOtherCustomValue")(r)
	require.ElementsMatch(t, []string{"MyCustomValue", "MyOtherCustomValue"}, r.Header.Values("X-Custom-Header"))
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

func RequireQueryParam(t *testing.T, key string) apitest.ValidatorFunc {
	return func(req *http.Request) error {
		if req.URL.Query().Get(key) == "" {
			return errors.New("missing query parameter: " + key)
		}
		return nil
	}
}

func RequireBodyParam(t *testing.T, key string) apitest.ValidatorFunc {
	return func(req *http.Request) error {
		bs, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		m := make(map[string]any)
		req.Body = io.NopCloser(bytes.NewReader(bs))
		if err := json.Unmarshal(bs, &m); err != nil {
			return err
		}

		if _, ok := m[key]; !ok {
			return errors.New("missing body parameter: " + key)
		}
		return nil
	}
}

func QueryParamEquals(t *testing.T, key, value string) apitest.ValidatorFunc {
	return func(req *http.Request) error {
		actual := req.URL.Query().Get(key)
		if actual == "" {
			return errors.New("missing query parameter: " + key)
		}

		if actual != value {
			return errors.New("expected query parameter " + key + " to be " + value + ", got " + actual)
		}
		return nil
	}
}

func BodyParamEquals(t *testing.T, key, value string) apitest.ValidatorFunc {
	return func(req *http.Request) error {
		bs, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		m := make(map[string]any)
		req.Body = io.NopCloser(bytes.NewReader(bs))
		if err := json.Unmarshal(bs, &m); err != nil {
			return err
		}

		val, ok := m[key]
		if !ok {
			return errors.New("missing body parameter: " + key)
		}

		if val != value {
			return errors.New("expected body parameter " + key + " to be " + value + ", got " + fmt.Sprintf("%v", val))
		}
		return nil
	}
}
