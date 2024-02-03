package google

const (
	URL = "https://maps.googleapis.com/maps/api/geocode/json"
)

type Response struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds map[string]struct {
				Latitude  float64 `json:"lat"`
				Longitude float64 `json:"lng"`
			} `json:"bounds"`
			Location struct {
				Latitude  float64 `json:"lat"`
				Longitude float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     map[string]struct {
				Latitude  float64 `json:"lat"`
				Longitude float64 `json:"lng"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`

	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}
