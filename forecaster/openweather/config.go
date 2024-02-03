package openweather

type Config struct {
	Enabled bool   `json:"enabled,omitempty"`
	APIKey  string `json:"api_key,omitempty"`
}
