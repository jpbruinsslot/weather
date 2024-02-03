package openweather

const (
	URL = "https://api.openweathermap.org/data/2.5/weather"
)

type Response struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`

	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`

	Base string `json:"base"`

	Main struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`

	Visibility int `json:"visibility"`

	Wind struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
	} `json:"wind"`

	Rain struct {
		OneHour float64 `json:"1h"`
	}

	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`

	Dt int `json:"dt"`

	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`

	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}
