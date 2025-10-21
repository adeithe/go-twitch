package apitest_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/adeithe/go-twitch/apitest"
	"github.com/stretchr/testify/require"
)

func TestMockAPI_RequireQueryParam(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		r := &http.Request{URL: &url.URL{}}
		query := r.URL.Query()
		query.Set("key", "value")
		r.URL.RawQuery = query.Encode()
		require.NoError(t, apitest.RequireQueryParam("key")(r))
	})

	t.Run("Invalid", func(t *testing.T) {
		r := &http.Request{URL: &url.URL{}}
		require.Error(t, apitest.RequireQueryParam("key")(r))
	})
}

func TestMockAPI_QueryParamEquals(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		r := &http.Request{URL: &url.URL{}}
		query := r.URL.Query()
		query.Set("key", "value")
		r.URL.RawQuery = query.Encode()
		require.NoError(t, apitest.QueryParamEquals("key", "value")(r))
	})

	t.Run("Mismatch", func(t *testing.T) {
		r := &http.Request{URL: &url.URL{}}
		query := r.URL.Query()
		query.Set("key", "invalid")
		r.URL.RawQuery = query.Encode()
		require.Error(t, apitest.QueryParamEquals("key", "value")(r))
	})

	t.Run("Invalid", func(t *testing.T) {
		r := &http.Request{URL: &url.URL{}}
		require.Error(t, apitest.QueryParamEquals("key", "value")(r))
	})
}

func TestMockAPI_RequireBodyParam(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		r := &http.Request{Body: io.NopCloser(bytes.NewReader([]byte(`{"key":"value"}`)))}
		err := apitest.RequireBodyParam("key")(r)
		require.NoError(t, err)
	})

	t.Run("Invalid", func(t *testing.T) {
		r := &http.Request{Body: io.NopCloser(bytes.NewReader([]byte("{}")))}
		require.Error(t, apitest.RequireBodyParam("key")(r))
	})
}

func TestMockAPI_BodyParamEquals(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		r := &http.Request{Body: io.NopCloser(bytes.NewReader([]byte(`{"key":"value"}`)))}
		require.NoError(t, apitest.BodyParamEquals("key", "value")(r))
	})

	t.Run("Mismatch", func(t *testing.T) {
		r := &http.Request{Body: io.NopCloser(bytes.NewReader([]byte(`{"key":"invalid"}`)))}
		require.Error(t, apitest.BodyParamEquals("key", "value")(r))
	})

	t.Run("Invalid", func(t *testing.T) {
		r := &http.Request{Body: io.NopCloser(bytes.NewReader([]byte(`{}`)))}
		require.Error(t, apitest.BodyParamEquals("key", "value")(r))
	})
}
