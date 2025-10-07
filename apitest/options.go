package apitest

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
