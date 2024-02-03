package ipgeolocation

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/jpbruinsslot/weather/geocoder"
)

type IPGeolocation struct {
	APIKey string
}

func New(apiKey string) *IPGeolocation {
	return &IPGeolocation{
		APIKey: apiKey,
	}
}

func (i *IPGeolocation) GetLocation() (geocoder.Location, error) {
	var location geocoder.Location

	// Get location
	resp, err := http.Get(fmt.Sprintf("%s?apiKey=%s", URL, i.APIKey))
	if err != nil {
		return location, fmt.Errorf("error: %v", err)
	}
	defer resp.Body.Close()

	var geoResp Response
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return location, fmt.Errorf("error: %v", err)
	}

	location.Latitude, err = strconv.ParseFloat(geoResp.Latitude, 64)
	if err != nil {
		return location, fmt.Errorf("error: %v", err)
	}

	location.Longitude, err = strconv.ParseFloat(geoResp.Longitude, 64)
	if err != nil {
		return location, fmt.Errorf("error: %v", err)
	}

	return location, nil
}

func (i *IPGeolocation) getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
