package kraken

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// IngestsData stores a list of ingests returned by the Twitch Ingests API.
type IngestsData struct {
	Data []IngestServer `json:"ingests"`
}

// IngestServer stores data about a Twitch RTMP Server.
type IngestServer struct {
	ID           int     `json:"_id"`
	Name         string  `json:"name"`
	Availability float32 `json:"availability"`
	URLTemplate  string  `json:"url_template"`
	Default      bool    `json:"default"`
}

// GetIngestServers retrieves a list of RTMP Servers used to stream on Twitch.
//
// See: https://dev.twitch.tv/docs/v5/reference/ingests#get-ingest-server-list
func (client *Client) GetIngestServers() (*IngestsData, error) {
	res, err := client.Request(http.MethodGet, fmt.Sprintf("ingests?_t=%d", time.Now().UTC().Unix()), nil)
	if err != nil {
		return nil, err
	}
	ingests := &IngestsData{}
	if err := json.Unmarshal(res.Body, &ingests); err != nil {
		return nil, err
	}
	return ingests, nil
}
