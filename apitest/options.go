package apitest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// MockTwitchAPIOption represents an option for configuring the MockTwitchAPI.
type MockTwitchAPIOption func(*MockTwitchAPI)

// EnableHTTP2 configures the mock server to support HTTP/2.
func EnableHTTP2() MockTwitchAPIOption {
	return func(m *MockTwitchAPI) {
		m.server.EnableHTTP2 = true
	}
}

// WithTLS configures the mock server to use HTTPS.
func WithTLS() MockTwitchAPIOption {
	return func(m *MockTwitchAPI) {
		m.tls = true
	}
}

// RequireQueryParam returns a ValidatorFunc that checks for the presence of a query parameter.
func RequireQueryParam(key string) ValidatorFunc {
	return func(req *http.Request) error {
		if req.URL.Query().Get(key) == "" {
			return errors.New("missing query parameter: " + key)
		}
		return nil
	}
}

// RequireBodyParam returns a ValidatorFunc that checks for the presence of a body parameter.
func RequireBodyParam(key string) ValidatorFunc {
	return func(req *http.Request) error {
		m := make(map[string]any)
		bs, _ := io.ReadAll(req.Body)
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

// QueryParamEquals returns a ValidatorFunc that checks if a query parameter equals a specific value.
func QueryParamEquals(key, value string) ValidatorFunc {
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

// BodyParamEquals returns a ValidatorFunc that checks if a body parameter equals a specific value.
func BodyParamEquals(key, value string) ValidatorFunc {
	return func(req *http.Request) error {
		m := make(map[string]any)
		bs, _ := io.ReadAll(req.Body)
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
