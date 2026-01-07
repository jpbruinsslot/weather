package forecaster

type Forecast struct {
	Temperature   float64
	Icon          string
	Rain          bool
	WindDirection string
	WindSpeed     float64
}

type Forecaster interface {
	GetForecast() (Forecast, error)
}
