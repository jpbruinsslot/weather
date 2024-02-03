package config

import (
	"github.com/jpbruinsslot/weather/forecaster/met"
	"github.com/jpbruinsslot/weather/forecaster/openweather"
	"github.com/jpbruinsslot/weather/geocoder/google"
	"github.com/jpbruinsslot/weather/geocoder/ipgeolocation"
)

type ForecastProviders struct {
	OpenWeather openweather.Config `json:"openweather,omitempty"`
	Met         met.Config         `json:"met,omitempty"`
}

type GeocodeProviders struct {
	Google        google.Config        `json:"google,omitempty"`
	IPGeolocation ipgeolocation.Config `json:"ipgeolocation,omitempty"`
}
