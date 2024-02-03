package met

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
	"github.com/jpbruinsslot/weather/utils/conv"
)

type Met struct {
	Location  geocoder.Location
	UserAgent string
	Units     string
}

func New(l geocoder.Location, userAgent string, units string) Met {
	return Met{
		Location:  l,
		UserAgent: userAgent,
		Units:     units,
	}
}

func (m Met) GetForecast() (forecaster.Forecast, error) {
	var forecast forecaster.Forecast

	var apiURL bytes.Buffer
	tmpl, _ := template.New("forecast").Parse("{{.URI}}?lat={{.Lat}}&lon={{.Lon}}")
	err := tmpl.Execute(&apiURL, map[string]string{
		"URI": URL,
		"Lat": fmt.Sprintf("%f", m.Location.Latitude),
		"Lon": fmt.Sprintf("%f", m.Location.Longitude),
	})
	if err != nil {
		return forecast, fmt.Errorf("error executing template: %s", err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return forecast, fmt.Errorf("error constructing request: %s", err)
	}

	req.Header.Set("User-Agent", m.UserAgent)

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

	var metResponse Response
	if err := json.Unmarshal(body, &metResponse); err != nil {
		return forecast, fmt.Errorf("error decoding forecast: %s", err)
	}

	slog.Debug(
		"met.no response",
		"url", apiURL.String(),
		"body", string(body),
		"status", resp.Status,
		"status_code", resp.StatusCode,
	)

	// Set icon
	forecast.Icon = CodesToIcons[metResponse.Properties.Timeseries[0].Data.Next1Hours.Summary.SymbolCode]

	// Set temperature
	temperature := metResponse.Properties.Timeseries[0].Data.Instant.Details.AirTemperature
	if m.Units == conv.Imperial {
		temperature = conv.CelsiusToFahrenheit(temperature)
	}
	forecast.Temperature = temperature

	// Set rain
	if metResponse.Properties.Timeseries[0].Data.Next1Hours.Details.PrecipitationAmount > 0 {
		forecast.Rain = true
	}

	return forecast, nil
}
