package helix

import (
	"encoding/json"
)

type Response struct {
	Data       json.RawMessage `json:"data"`
	Pagination Pagination      `json:"pagination,omitempty"`
	Error      string          `json:"error,omitempty"`
	Status     int             `json:"status,omitempty"`
	Message    string          `json:"message,omitempty"`
}

type Pagination struct {
	Cursor string `json:"cursor,omitempty"`
}
