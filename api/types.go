package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ResponseData[T any] struct {
	Total  *int `json:"total,omitempty"`  // Only present in some endpoints.
	Points *int `json:"points,omitempty"` // Only present in some endpoints.

	Data       []T         `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`

	Status  int    `json:"status"`            // If not provided by Twitch, defaults to HTTP status code.
	Code    string `json:"error"`             // If not provided by Twitch, defaults to HTTP status text.
	Message string `json:"message,omitempty"` // Only present if status is non-200
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"error"`
	Message string `json:"message"`
}

func decodeResponse[T any](res *http.Response) (*ResponseData[T], error) {
	var data ResponseData[T]
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.Status == 0 {
		data.Status = res.StatusCode
		data.Code = http.StatusText(res.StatusCode)
	}

	if err := data.asError(); err != nil {
		return nil, err
	}
	return &data, nil
}

func (data ResponseData[T]) asError() error {
	if data.Status <= 400 {
		return nil
	}
	return &APIError{data.Status, data.Code, data.Message}
}

func (err APIError) Error() string {
	return fmt.Sprintf("twitchapi: %d %s - %s", err.Status, err.Code, err.Message)
}

// CodeOf returns the HTTP status code of the given error.
// If the error is not an API error, it returns http.StatusInternalServerError.
func CodeOf(err error) int {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Status
	}
	return http.StatusInternalServerError
}
