package api_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Adeithe/go-twitch/api"
	"github.com/stretchr/testify/assert"
)

func TestAPI_ClipsDuration(t *testing.T) {
	tests := []struct {
		Input    float64
		Expected time.Duration
	}{
		{0.0, 0},
		{0.1, time.Millisecond * 100},
		{0.5, time.Millisecond * 500},
		{1.0, time.Second},
		{1.5, time.Millisecond * 1500},
		{2.0, time.Second * 2},
		{2.5, time.Millisecond * 2500},
		{20.4, time.Millisecond * 20400},
	}

	for _, tt := range tests {
		actual := struct {
			Duration api.ClipDuration `json:"duration"`
		}{}

		if err := json.Unmarshal([]byte(fmt.Sprintf(`{"duration":%f}`, tt.Input)), &actual); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assert.Equal(t, tt.Expected, actual.Duration.AsDuration())
	}
}
