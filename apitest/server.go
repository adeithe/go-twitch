// Package apitest provides utilities for testing the Twitch API Client.
package apitest

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/adeithe/go-twitch/api"
	"golang.org/x/oauth2"
)

// TestingT is an interface for the [testing.T] type.
type TestingT interface {
	Cleanup(func())
}

// MockTwitchAPI represents a mock Twitch API server.
type MockTwitchAPI struct {
	mux          *http.ServeMux
	server       *httptest.Server
	applications map[string]string
	tokens       map[string]string
	handlers     map[string]*MockTwitchAPIEndpoint
	oauthToken   *MockTwitchAPIEndpoint
	mx           sync.RWMutex
	tls          bool

	BaseURL string
}

// MockTwitchAPIEndpoint represents available data for a mock Twitch API endpoint.
type MockTwitchAPIEndpoint struct {
	data       any
	validators []ValidatorFunc

	TimesCalled int
	Successes   int
	Failures    int
}

// ValidatorFunc is a function type for validating incoming HTTP requests.
type ValidatorFunc func(req *http.Request) error

// SetMockValidator adds a mock validator for the given path.
//
// Successful requests will return an empty response.
//
// If a response has already been added for the given path and method, it will be overwritten.
func SetMockValidator(m *MockTwitchAPI, method, path string, validators ...ValidatorFunc) *MockTwitchAPIEndpoint {
	return SetMockResponse(m, method, path, &api.ResponseData[any]{}, validators...)
}

// SetMockResponse adds a mock response for the given path.
//
// If a response has already been added for the given path and method, it will be overwritten.
func SetMockResponse[T any, V api.ResponseData[T]](m *MockTwitchAPI, method, path string, data *V, validators ...ValidatorFunc) *MockTwitchAPIEndpoint {
	m.mx.RLock()
	key := fmt.Sprintf("%s %s", strings.ToUpper(method), path)
	handler, ok := m.handlers[key]
	m.mx.RUnlock()
	if !ok {
		m.mx.Lock()
		handler = &MockTwitchAPIEndpoint{}
		m.handlers[key] = handler
		m.mx.Unlock()
	}
	handler.data = data
	handler.validators = validators
	return handler
}

// NewMockAPI creates a mock Twitch API server for testing purposes.
func NewMockAPI(t TestingT, opts ...MockTwitchAPIOption) *MockTwitchAPI {
	mux := http.NewServeMux()
	srv := httptest.NewUnstartedServer(http.HandlerFunc(mux.ServeHTTP))
	mock := &MockTwitchAPI{
		mux:          mux,
		server:       srv,
		applications: make(map[string]string),
		tokens:       make(map[string]string),
		handlers:     make(map[string]*MockTwitchAPIEndpoint),
		oauthToken:   &MockTwitchAPIEndpoint{},

		BaseURL: srv.URL,
	}

	for _, opt := range opts {
		opt(mock)
	}

	if mock.tls {
		srv.StartTLS()
	} else {
		srv.Start()
	}

	t.Cleanup(srv.Close)
	client := mock.Client()
	url, _ := url.Parse(srv.URL)
	client.Transport = newMockTransport(client, url)
	mux.Handle("/id/twitch/tv/oauth2/token", mock)
	mux.Handle("/api/twitch/tv/{rest...}", mock)
	return mock
}

// OAuthTokenEndpoint returns the mock endpoint for the Twitch API OAuth token URL.
func (m *MockTwitchAPI) OAuthTokenEndpoint() *MockTwitchAPIEndpoint {
	return m.oauthToken
}

// RegisterApplication simulates registering a new Twitch application and returns a client ID and secret.
func (m *MockTwitchAPI) RegisterApplication() (clientID, clientSecret string, err error) {
	c := make([]byte, 16)
	if _, err = rand.Read(c); err != nil {
		return
	}
	clientID = fmt.Sprintf("%x", c)[2:]

	s := make([]byte, 15)
	if _, err = rand.Read(s); err != nil {
		return
	}
	clientSecret = fmt.Sprintf("%x", s)[2:]

	m.mx.Lock()
	defer m.mx.Unlock()
	m.applications[clientID] = clientSecret
	return
}

// NewBearerToken simulates generating a new bearer token for the given client ID.
func (m *MockTwitchAPI) NewBearerToken(clientID string, scopes ...string) (*oauth2.Token, error) {
	tokenBs := make([]byte, 15)
	if _, err := rand.Read(tokenBs); err != nil {
		return nil, err
	}

	refreshBs := make([]byte, 15)
	if _, err := rand.Read(refreshBs); err != nil {
		return nil, err
	}

	m.mx.Lock()
	defer m.mx.Unlock()
	accessToken := fmt.Sprintf("%x", tokenBs)[2:]
	m.tokens[accessToken] = clientID
	token := &oauth2.Token{
		TokenType:    "bearer",
		AccessToken:  accessToken,
		RefreshToken: fmt.Sprintf("%x", refreshBs)[2:],
		Expiry:       time.Now().UTC().Add(time.Hour * 24),
	}
	return token.WithExtra(map[string]any{"scope": scopes}), nil
}

// Certificate returns the TLS certificate used by the mock server.
func (m *MockTwitchAPI) Certificate() *x509.Certificate {
	return m.server.Certificate()
}

// Client returns a HTTP client to use with the mock server.
func (m *MockTwitchAPI) Client() *http.Client {
	return m.server.Client()
}

// ServeHTTP implements the [net/http.Handler] interface for MockTwitchAPI.
func (m *MockTwitchAPI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/id/twitch/tv/") {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/id/twitch/tv")
		m.handleOAuth(res, req)
		return
	}

	writer := json.NewEncoder(res)
	data := &api.ResponseData[any]{Status: http.StatusOK}
	defer func() {
		if data.Status < http.StatusBadRequest {
			return
		}

		res.WriteHeader(data.Status)
		_ = writer.Encode(api.TwitchAPIError{
			Status:  data.Status,
			Code:    http.StatusText(data.Status),
			Message: data.Message,
		})
	}()

	m.mx.RLock()
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api/twitch/tv/")
	handler, ok := m.handlers[fmt.Sprintf("%s /%s", strings.ToUpper(req.Method), req.URL.Path)]
	m.mx.RUnlock()
	if !ok {
		data.Status = http.StatusNotFound
		return
	}

	m.mx.Lock()
	defer m.mx.Unlock()
	handler.TimesCalled++
	token := strings.TrimSpace(req.Header.Get("Authorization"))
	if token == "" {
		handler.Failures++
		data.Status = http.StatusUnauthorized
		data.Message = "OAuth token is missing"
		return
	}

	authorization := strings.SplitN(token, " ", 2)
	clientID := strings.TrimSpace(req.Header.Get("Client-ID"))
	if clientID == "" {
		handler.Failures++
		data.Status = http.StatusUnauthorized
		data.Message = "Client ID is missing"
		return
	}

	if len(authorization) != 2 || !strings.EqualFold(authorization[0], "Bearer") {
		handler.Failures++
		data.Status = http.StatusUnauthorized
		data.Message = "OAuth token is missing or malformed"
		return
	}

	authClient, ok := m.tokens[authorization[1]]
	if !ok {
		handler.Failures++
		data.Status = http.StatusUnauthorized
		data.Message = "Invalid OAuth token"
		return
	}

	if clientID != authClient {
		handler.Failures++
		data.Status = http.StatusUnauthorized
		data.Message = "Client ID does not match OAuth token"
		return
	}

	for _, validate := range handler.validators {
		if err := validate(req); err != nil {
			handler.Failures++
			data.Status = http.StatusBadRequest
			data.Message = err.Error()
			return
		}
	}

	handler.Successes++
	res.WriteHeader(data.Status)
	_ = writer.Encode(handler.data)
}

func (m *MockTwitchAPI) handleOAuth(res http.ResponseWriter, req *http.Request) {
	m.oauthToken.TimesCalled++
	if req.Method != http.MethodPost || req.URL.Path != "/oauth2/token" {
		m.oauthToken.Failures++
		res.WriteHeader(http.StatusNotFound)
		_, _ = res.Write([]byte("404 Not Found"))
		return
	}

	writer := json.NewEncoder(res)
	if err := req.ParseForm(); err != nil {
		m.oauthToken.Failures++
		res.WriteHeader(http.StatusBadRequest)
		_ = writer.Encode(api.TwitchAPIError{
			Status:  http.StatusBadRequest,
			Message: "malformed request",
		})
		return
	}

	query := req.Form
	grantType, clientID, clientSecret := query.Get("grant_type"), query.Get("client_id"), query.Get("client_secret")
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	m.mx.RLock()
	secret, ok := m.applications[clientID]
	m.mx.RUnlock()
	if !ok || clientSecret != secret {
		m.oauthToken.Failures++
		res.WriteHeader(http.StatusBadRequest)
		_ = writer.Encode(api.TwitchAPIError{
			Status:  http.StatusBadRequest,
			Message: "invalid client",
		})
		return
	}

	var refreshToken string
	scopes := strings.Split(query.Get("scope"), " ")
	token, err := m.NewBearerToken(clientID, scopes...)
	if err != nil {
		m.oauthToken.Failures++
		res.WriteHeader(http.StatusInternalServerError)
		_ = writer.Encode(api.TwitchAPIError{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	switch grantType {
	case "client_credentials":
	case "code":
		refreshToken = token.RefreshToken
	default:
		m.oauthToken.Failures++
		res.WriteHeader(http.StatusBadRequest)
		_ = writer.Encode(api.TwitchAPIError{
			Status:  http.StatusBadRequest,
			Message: "invalid grant type",
		})
		return
	}

	m.oauthToken.Successes++
	res.WriteHeader(http.StatusOK)
	_ = writer.Encode(struct {
		TokenType    string   `json:"token_type"`
		AccessToken  string   `json:"access_token"`
		RefreshToken string   `json:"refresh_token,omitempty"`
		Scope        []string `json:"scope,omitempty"`
		ExpiresIn    int      `json:"expires_in"`
	}{
		TokenType:    "bearer",
		AccessToken:  token.AccessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int((time.Hour * 24 * 60).Seconds()),
		Scope:        token.Extra("scope").([]string),
	})
}
