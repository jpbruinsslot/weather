package ipgeolocation

const (
	URL = "https://api.ipgeolocation.io/ipgeo"
)

// https://ipgeolocation.io/documentation/ip-geolocation-api.html
type Response struct {
	IP        string `json:"ip"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
