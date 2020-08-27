package helix

import (
	"encoding/json"
)

type Response struct {
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error,omitempty"`
	Status  int             `json:"status,omitempty"`
	Message string          `json:"message,omitempty"`
}
