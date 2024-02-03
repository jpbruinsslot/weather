package geocoder

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Elevation float64 `json:"elevation"`
}

type Locator interface {
	GetLocation() (Location, error)
}
