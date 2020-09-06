package helix

import (
	"encoding/json"
)

// Response stores data about a basic response from the Twitch Helix API.
type Response struct {
	Data       json.RawMessage `json:"data"`
	Pagination Pagination      `json:"pagination,omitempty"`
	Error      string          `json:"error,omitempty"`
	Status     int             `json:"status,omitempty"`
	Message    string          `json:"message,omitempty"`
}

// Pagination stores data about a paged response from the Twitch Helix API.
type Pagination struct {
	Cursor string `json:"cursor,omitempty"`
}
