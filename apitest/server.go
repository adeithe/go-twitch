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

	"github.com/adeithe/go-twitch/api"
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
	mx           sync.RWMutex
	tls          bool

	BaseURL string
}

// MockTwitchAPIEndpoint represents available data for a mock Twitch API endpoint.
type MockTwitchAPIEndpoint struct {
	data any

	TimesCalled, Successes, Failures int
}

// SetMockResponse adds a mock response for the given path.
//
// If a response has already been added for the given path and method, it will be overwritten.
func SetMockResponse[T any, V api.ResponseData[T]](m *MockTwitchAPI, method string, path string, data *V) *MockTwitchAPIEndpoint {
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
	mux.Handle("/{rest...}", mock)
	return mock
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
func (m *MockTwitchAPI) NewBearerToken(clientID string) (token string, err error) {
	t := make([]byte, 15)
	if _, err = rand.Read(t); err != nil {
		return
	}

	m.mx.Lock()
	defer m.mx.Unlock()
	token = fmt.Sprintf("%x", t)[2:]
	m.tokens[token] = clientID
	return
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
	writer := json.NewEncoder(res)
	data := &api.ResponseData[any]{Status: http.StatusOK}
	defer func() {
		if data.Status < 400 {
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
	handler, ok := m.handlers[fmt.Sprintf("%s %s", strings.ToUpper(req.Method), req.URL.Path)]
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

	if clientID == "" {
		handler.Failures++
		data.Status = http.StatusUnauthorized
		data.Message = "Client ID is missing"
		return
	}

	if clientID != authClient {
		handler.Failures++
		data.Status = http.StatusUnauthorized
		data.Message = "Client ID does not match OAuth token"
		return
	}

	handler.Successes++
	res.WriteHeader(data.Status)
	_ = writer.Encode(handler.data)
}
