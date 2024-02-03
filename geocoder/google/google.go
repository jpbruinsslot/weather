package google

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jpbruinsslot/weather/geocoder"
)

type Google struct {
	APIKey         string
	LocationString string
}

func New(apiKey, locationString string) *Google {
	return &Google{
		APIKey:         apiKey,
		LocationString: locationString,
	}
}

func (g Google) GetLocation() (geocoder.Location, error) {
	var location geocoder.Location

	data := url.Values{"address": {g.LocationString}, "key": {g.APIKey}}
	apiURL := fmt.Sprintf("%s?%s", URL, data.Encode())

	resp, err := http.Get(apiURL)
	if err != nil {
		return location, fmt.Errorf("error: %v", err)
	}
	defer resp.Body.Close()

	var geoResp Response
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return location, fmt.Errorf("error: %v", err)
	}

	if geoResp.ErrorMessage != "" {
		return location, fmt.Errorf("error: %v", geoResp.ErrorMessage)
	}

	if len(geoResp.Results) <= 0 {
		return location, fmt.Errorf("error: %v", "No results")
	}

	result := geoResp.Results[0]
	location.Latitude = result.Geometry.Location.Latitude
	location.Longitude = result.Geometry.Location.Longitude

	return location, nil
}

func AutoLocate() {}

func IPLocate() {}

func Locate() {}
