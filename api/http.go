package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HTTPRequest struct {
	BaseURL string
	Path    string
	Method  string
	Headers map[string]string
	Body    json.RawMessage
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
}

// HTTPClient used to handle HTTP Requests
var HTTPClient http.Client = http.Client{Timeout: time.Duration(time.Second * 5)}

// NewRequest prepares data for a HTTPRequest
func NewRequest(method string, url string, path string) *HTTPRequest {
	return &HTTPRequest{
		BaseURL: url,
		Path:    path,
		Method:  method,
		Headers: make(map[string]string),
	}
}

// Do the HTTP Request
func (req HTTPRequest) Do() (HTTPResponse, error) {
	response := &HTTPResponse{}
	url := strings.TrimSuffix(req.BaseURL, "/") + "/" + strings.TrimPrefix(req.Path, "/")
	r, err := http.NewRequest(strings.ToUpper(req.Method), url, bytes.NewBuffer(req.Body))
	for key, value := range req.Headers {
		r.Header.Set(key, value)
	}
	resp, err := HTTPClient.Do(r)
	if err != nil {
		return *response, err
	}
	response.StatusCode = resp.StatusCode
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return *response, err
	}
	response.Body = body
	return *response, nil
}
