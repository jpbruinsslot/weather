package openweather

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/jpbruinsslot/weather/forecaster"
	"github.com/jpbruinsslot/weather/geocoder"
)

type OpenWeather struct {
	APIKey    string
	Location  geocoder.Location
	UserAgent string
	Units     string
}

func New(apiKey string, location geocoder.Location, userAgent string, units string) *OpenWeather {
	return &OpenWeather{
		APIKey:    apiKey,
		Location:  location,
		UserAgent: userAgent,
		Units:     units,
	}
}

func (o OpenWeather) GetForecast() (forecaster.Forecast, error) {
	var forecast forecaster.Forecast

	var apiURL bytes.Buffer
	tmpl, _ := template.New("forecast").
		Parse("{{.URI}}/?lat={{.Lat}}&lon={{.Lng}}&appid={{.APIKey}}&units={{.Units}}")

	err := tmpl.Execute(&apiURL, map[string]string{
		"URI":    URL,
		"Lat":    fmt.Sprintf("%f", o.Location.Latitude),
		"Lng":    fmt.Sprintf("%f", o.Location.Longitude),
		"APIKey": o.APIKey,
		"Units":  o.Units,
	})
	if err != nil {
		return forecast, fmt.Errorf("error executing template: %s", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return forecast, fmt.Errorf("error constructing request: %s", err)
	}

	req.Header.Set("User-Agent", o.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return forecast, fmt.Errorf("error executing request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return forecast, fmt.Errorf("error received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return forecast, fmt.Errorf("error reading response body: %s", err)
	}

	var openweatherResp Response
	if err := json.Unmarshal(body, &openweatherResp); err != nil {
		return forecast, fmt.Errorf("error decoding forecast: %s", err)
	}

	// TODO: when night time, use night icons

	slog.Debug(
		"openweather response",
		"url", apiURL.String(),
		"body", string(body),
		"status", resp.Status,
		"status_code", resp.StatusCode,
	)

	// Set icon
	forecast.Icon = CodesToIconsDay[openweatherResp.Weather[0].ID]

	// Set temperature
	forecast.Temperature = openweatherResp.Main.Temp

	// Set rain
	forecast.Rain = false
	if openweatherResp.Rain.OneHour > 0 {
		forecast.Rain = true
	}

	return forecast, nil
}
