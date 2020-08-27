package kraken

type Error struct {
	Error   string `json:"error,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
