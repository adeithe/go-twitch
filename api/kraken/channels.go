package kraken

type Channel struct {
	ID                           int    `json:"_id"`
	Login                        string `json:"name"`
	DisplayName                  string `json:"display_name"`
	Game                         string `json:"game"`
	Status                       string `json:"status"`
	Description                  string `json:"description"`
	Logo                         string `json:"logo"`
	VideoBanner                  string `json:"video_banner"`
	ProfileBanner                string `json:"profile_banner"`
	ProfileBannerBackgroundColor string `json:"profile_banner_background_color"`
	URL                          string `json:"url"`
	Views                        int    `json:"views"`
	Followers                    int    `json:"followers"`
	Partner                      bool   `json:"partner"`
	Mature                       bool   `json:"mature"`
	PrivateVideo                 bool   `json:"private_video"`
	PrivacyOptionsEnabled        bool   `json:"privacy_options_enabled"`
	StreamKey                    string `json:"stream_key"`
	BroadcastLanguage            string `json:"broadcast_language"`
	BroadcastSoftware            string `json:"broadcast_software"`
	Language                     string `json:"language"`
	Type                         string `json:"broadcaster_type"`
	CreatedAt                    string `json:"created_at"`
	UpdatedAt                    string `json:"updated_at"`
}
