package forecaster

type Forecast struct {
	Temperature float64
	Icon        string
	Rain        bool
}

type Forecaster interface {
	GetForecast() (Forecast, error)
}
