package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HTTPRequest stores data for a HTTP request.
type HTTPRequest struct {
	BaseURL string
	Path    string
	Method  string
	Headers map[string]string
	Body    []byte
}

// HTTPResponse contains data about a previously handled HTTP request.
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// HTTPClient used to handle HTTP Requests.
var HTTPClient http.Client = http.Client{Timeout: time.Duration(time.Second * 5)}

// New prepares data for a HTTPRequest.
func New(method string, url string, path string) *HTTPRequest {
	return &HTTPRequest{
		BaseURL: url,
		Path:    path,
		Method:  method,
		Headers: make(map[string]string),
	}
}

// Do the HTTP Request.
func (req HTTPRequest) Do() (HTTPResponse, error) {
	response := &HTTPResponse{Headers: make(map[string]string)}
	url := strings.TrimSuffix(req.BaseURL, "/") + "/" + strings.TrimPrefix(req.Path, "/")
	var reqBody *bytes.Buffer
	if len(req.Body) > 0 {
		reqBody = bytes.NewBuffer(req.Body)
		req.Headers["Content-Length"] = fmt.Sprint(len(req.Body))
		if _, ok := req.Headers["Content-Type"]; !ok {
			req.Headers["Content-Type"] = "application/json"
		}
	}
	r, err := http.NewRequest(strings.ToUpper(req.Method), url, reqBody)
	if err != nil {
		return *response, err
	}
	for key, value := range req.Headers {
		r.Header.Set(key, value)
	}
	resp, err := HTTPClient.Do(r)
	if err != nil {
		return *response, err
	}
	defer resp.Body.Close()
	response.StatusCode = resp.StatusCode
	for key, value := range resp.Header {
		response.Headers[key] = strings.Join(value, ",")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return *response, err
	}
	response.Body = body
	return *response, IsError(response.Body)
}
