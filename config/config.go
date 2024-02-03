package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"

	"github.com/jpbruinsslot/weather/icons"
	"github.com/jpbruinsslot/weather/utils/colors"
	"github.com/jpbruinsslot/weather/utils/conv"
)

type Config struct {
	Path string `json:"path"`

	Location  string  `json:"location"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Colors colors.Colors `json:"-"`
	Format string        `json:"format"`
	Units  string        `json:"units"`

	IconType string      `json:"icon_type"`
	Icons    icons.Icons `json:"icons"`

	Forecaster ForecastProviders `json:"forecaster"`
	Geocoder   GeocodeProviders  `json:"geocoder"`

	Interval int `json:"interval"`
}

func (c *Config) Save() error {
	f, err := os.Create(c.Path)
	if err != nil {
		return fmt.Errorf("couldn't open the config file: (%v)", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("couldn't encode the config file: (%v)", err)
	}

	return nil
}

func New(flgConfigPath string) (*Config, error) {
	cfg := getDefaultConfig()

	f, err := os.Open(flgConfigPath)
	if err != nil {
		return &cfg, fmt.Errorf("couldn't open the config file: (%v)", err)
	}
	defer f.Close()

	cfg.Path = flgConfigPath

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		// If file is empty, save the default config
		if err == io.EOF {
			err = cfg.Save()
			if err != nil {
				return &cfg, fmt.Errorf("couldn't save the config file: (%v)", err)
			}
		} else {
			return &cfg, fmt.Errorf("the config file isn't valid json: (%v)", err)
		}
	}

	// Merge the selected iconset with the custom icons
	var iconSet icons.Icons
	switch cfg.IconType {
	case icons.IconTypeNerdFonts:
		iconSet = icons.NerdFonts()
	case icons.IconTypeUnicode:
		iconSet = icons.Unicode()
	default:
		iconSet = icons.Unicode()
	}

	maps.Copy(iconSet, cfg.Icons)
	cfg.Icons = iconSet

	slog.Debug(fmt.Sprintf("using icon type: %s", cfg.IconType))

	if cfg.Units != conv.Metric && cfg.Units != conv.Imperial {
		return &cfg, fmt.Errorf("the units must be either 'metric' or 'imperial'")
	}

	return &cfg, nil
}

func getDefaultConfig() Config {
	return Config{
		Colors: colors.DefaultANSIColors,
		Format: "{{if .Rain}}{{.Rain}} {{end}}{{.Icon}}  {{printf `%.1f` .Temperature}}{{.Units}}",
		Units:  conv.Metric,

		IconType: icons.IconTypeUnicode,

		Interval: 60,
	}
}
