package latlng

import (
	"github.com/jpbruinsslot/weather/geocoder"
)

type LatLng struct {
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lng,omitempty"`
}

func (l LatLng) GetLocation() (geocoder.Location, error) {
	return geocoder.Location{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
	}, nil
}
