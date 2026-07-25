package conv

func DegreesToCompass(degrees float64) string {
	const (
		north     = "wind-north"
		northEast = "wind-north-east"
		east      = "wind-east"
		southEast = "wind-south-east"
		south     = "wind-south"
		southWest = "wind-south-west"
		west      = "wind-west"
		northWest = "wind-north-west"
	)

	switch {
	case degrees >= 348.75 || degrees < 11.25:
		return north
	case degrees >= 11.25 && degrees < 33.75:
		return northEast
	case degrees >= 33.75 && degrees < 56.25:
		return east
	case degrees >= 56.25 && degrees < 78.75:
		return southEast
	case degrees >= 78.75 && degrees < 101.25:
		return south
	case degrees >= 101.25 && degrees < 123.75:
		return southWest
	case degrees >= 123.75 && degrees < 146.25:
		return west
	case degrees >= 146.25 && degrees < 168.75:
		return northWest
	default:
		return north // Fallback case, should not happen
	}
}

func MeterPerSecondToBeaufort(mps float64) int {
	return int(mps * 3.6 / 1.61) // Convert m/s to km/h and then to Beaufort scale
}

func MeterPerSecondToMilesPerHour(mps float64) float64 {
	return mps * 2.23694
}

func MilesPerHourToMeterPerSecond(mph float64) float64 {
	return mph * 0.44704
}
