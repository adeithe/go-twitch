package request

import (
	"encoding/json"
	"fmt"
)

// APIError stores data returned from a Twitch API Error.
type APIError struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// IsError returns an error if the Twitch API responded with one.
func IsError(bytes []byte) error {
	var data APIError
	if err := json.Unmarshal(bytes, &data); err != nil {
		data = APIError{
			Error:   "Internal Server Error",
			Status:  500,
			Message: "Failed to parse response data",
		}
	}
	if data.Status > 0 && len(data.Error) > 0 && len(data.Message) > 0 {
		return fmt.Errorf("%s: %s (status: %d)", data.Error, data.Message, data.Status)
	}
	return nil
}
