package kraken

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// TeamOpts stores options for requests to the Twitch Teams API.
type TeamOpts struct {
	Limit  int
	Offset int
}

// TeamsData stores a list of teams returned by the Twitch Teams API.
type TeamsData struct {
	Data []Team `json:"teams"`
}

// Team stores data about a team on Twitch.
type Team struct {
	ID          int    `json:"_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Info        string `json:"info"`
	Logo        string `json:"logo"`
	Banner      string `json:"banner"`
	Background  string `json:"background"`
	Users       []struct {
		ID                           string `json:"_id"`
		Name                         string `json:"name"`
		DisplayName                  string `json:"display_name"`
		Title                        string `json:"status"`
		Game                         string `json:"game"`
		URL                          string `json:"url"`
		Logo                         string `json:"logo"`
		VideoBanner                  string `json:"video_banner"`
		ProfileBanner                string `json:"profile_banner"`
		ProfileBannerBackgroundColor string `json:"profile_banner_background_color"`
		Views                        int    `json:"views"`
		Followers                    int    `json:"followers"`
		IsMature                     bool   `json:"mature"`
		IsPartner                    bool   `json:"partner"`
		BroadcasterLanguage          string `json:"broadcaster_language"`
		Language                     string `json:"language"`
		CreatedAt                    string `json:"created_at"`
		UpdatedAt                    string `json:"updated_at"`
	} `json:"users"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetAllTeams retrieves active teams based on the specified TeamOpts.
//
// See: https://dev.twitch.tv/docs/v5/reference/teams#get-all-teams
func (client *Client) GetAllTeams(opts TeamOpts) (*TeamsData, error) {
	params := ""
	if opts.Limit > 0 {
		params += fmt.Sprint("&limit=", opts.Limit)
	}
	if opts.Offset > 0 {
		params += fmt.Sprint("&offset=", opts.Offset)
	}
	res, err := client.Request(http.MethodGet, fmt.Sprintf("teams?_t=%d%s", time.Now().UTC().Unix(), params), nil)
	if err != nil {
		return nil, err
	}
	teams := &TeamsData{}
	if err := json.Unmarshal(res.Body, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// GetTeam retrieves a specified team object by name.
//
// See: https://dev.twitch.tv/docs/v5/reference/teams#get-team
func (client *Client) GetTeam(name string) (*TeamsData, error) {
	if len(name) < 1 {
		return nil, errors.New("no team name provided")
	}
	res, err := client.Request(http.MethodGet, fmt.Sprintf("teams/%s?_t=%d", name, time.Now().UTC().Unix()), nil)
	if err != nil {
		return nil, err
	}
	teams := &TeamsData{}
	if err := json.Unmarshal(res.Body, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}
