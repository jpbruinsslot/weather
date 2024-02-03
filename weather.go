package weather

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"text/template"
	"time"

	"github.com/jpbruinsslot/weather/config"
	"github.com/jpbruinsslot/weather/forecaster"
	"github.com/jpbruinsslot/weather/forecaster/met"
	"github.com/jpbruinsslot/weather/forecaster/openweather"
	"github.com/jpbruinsslot/weather/geocoder"
	"github.com/jpbruinsslot/weather/geocoder/google"
	"github.com/jpbruinsslot/weather/geocoder/ipgeolocation"
	"github.com/jpbruinsslot/weather/geocoder/latlng"
	"github.com/jpbruinsslot/weather/utils/xdg"
)

const (
	Version   = "0.1.0"
	URL       = "github.com/jpbruinsslot/weather"
	UserAgent = "weather/" + Version + " (" + URL + ")"
)

type Weather struct {
	Config    *config.Config
	LastCheck time.Time

	Forecaster forecaster.Forecaster
	Forecast   forecaster.Forecast

	Locator  geocoder.Locator
	Location geocoder.Location
}

func New(config *config.Config) (*Weather, error) {
	weather := &Weather{
		Config: config,
	}

	// Geocode based on coordinates
	if config.Latitude != 0 && config.Longitude != 0 {
		weather.Locator = &latlng.LatLng{
			Latitude:  config.Latitude,
			Longitude: config.Longitude,
		}

		err := weather.GetLocation()
		if err != nil {
			return weather, err
		}
	}

	// If we don't have a location yet, try to get one based on the location string
	if (weather.Location == geocoder.Location{}) && config.Location != "" {
		switch {
		case config.Geocoder.Google != (google.Config{}) && config.Geocoder.Google.Enabled:
			if config.Geocoder.Google.APIKey == "" {
				return weather, errors.New("no Google API key provided")
			}

			weather.Locator = google.New(config.Geocoder.Google.APIKey, config.Location)
		default:
			return weather, errors.New("no location geocoder provided")
		}

		err := weather.GetLocation()
		if err != nil {
			return weather, err
		}
	}

	// Get location based on ip
	if (weather.Location == geocoder.Location{}) {
		switch {
		case config.Geocoder.IPGeolocation != (ipgeolocation.Config{}) && config.Geocoder.IPGeolocation.Enabled:
			if config.Geocoder.IPGeolocation.APIKey == "" {
				return weather, errors.New("no IPGeolocation API key provided")
			}
			weather.Locator = ipgeolocation.New(config.Geocoder.IPGeolocation.APIKey)
		default:
			weather.Locator = &latlng.LatLng{
				Latitude:  52.1015441,
				Longitude: 5.1779992,
			}
		}

		err := weather.GetLocation()
		if err != nil {
			return weather, err
		}
	}

	// What weather API are we using?
	switch {
	case config.Forecaster.Met != (met.Config{}) && config.Forecaster.Met.Enabled:
		weather.Forecaster = met.New(
			weather.Location,
			UserAgent,
			config.Units,
		)
	case config.Forecaster.OpenWeather != (openweather.Config{}) && config.Forecaster.OpenWeather.Enabled:
		weather.Forecaster = openweather.New(
			config.Forecaster.OpenWeather.APIKey,
			weather.Location,
			UserAgent,
			config.Units,
		)
	default:
		weather.Forecaster = met.New(
			weather.Location,
			UserAgent,
			config.Units,
		)
	}

	slog.Debug(
		fmt.Sprintf("setting up weather with forecaster %s and locator %s", weather.GetForecaster(), weather.GetLocator()),
		"forecaster", weather.GetForecaster(),
		"locator", weather.GetLocator(),
		"location", weather.Location,
	)

	return weather, nil
}

func (w *Weather) GetLocation() error {
	location, err := w.Locator.GetLocation()
	if err != nil {
		return err
	}

	w.Location = location

	return nil
}

func (w *Weather) GetLocator() string {
	return reflect.TypeOf(w.Locator).String()
}

func (w *Weather) GetForecast() error {
	// Load cached forecast
	err := w.Load()
	if err != nil {
		return err
	}

	// When the last check was less than the interval ago, use the cached forecast
	if time.Since(w.LastCheck) < time.Duration(w.Config.Interval)*time.Second {
		slog.Debug(
			"using cached forecast",
			"forecaster", w.GetForecaster(),
			"locator", w.GetLocator(),
		)

		return nil
	}

	// Get geocode
	if w.Location == (geocoder.Location{}) {
		err := w.GetLocation()
		if err != nil {
			return err
		}
	}

	// Get forecast
	forecast, err := w.Forecaster.GetForecast()
	if err != nil {
		return err
	}

	w.Forecast = forecast
	w.LastCheck = time.Now()

	// Save forecast
	err = w.Save()
	if err != nil {
		return err
	}

	return nil
}

func (w *Weather) GetForecaster() string {
	return reflect.TypeOf(w.Forecaster).String()
}

func (w *Weather) PrintForecast() error {
	output := w.GenerateOutput()

	templ, err := template.New("weather").Parse(w.Config.Format)
	if err != nil {
		return err
	}

	err = templ.Execute(os.Stdout, output)
	if err != nil {
		return err
	}

	return nil
}

func (w *Weather) Save() error {

	dataPath, err := xdg.DataFile("weather/data.gob")
	if err != nil {
		return err
	}

	f, err := os.Create(dataPath)
	if err != nil {
		return fmt.Errorf("couldn't open the data file: (%v)", err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(w); err != nil {
		return fmt.Errorf("couldn't encode the data file: (%v)", err)
	}

	return nil
}

func (w *Weather) Load() error {
	dataPath, err := xdg.DataFile("weather/data.gob")
	if err != nil {
		return err
	}

	f, err := os.Open(dataPath)
	if err != nil {
		return fmt.Errorf("couldn't open the data file: (%v)", err)
	}
	defer f.Close()

	if err := gob.NewDecoder(f).Decode(&w); err != nil {
		// If file is empty, save the initial data file
		if err == io.EOF {
			err = w.Save()
			if err != nil {
				return fmt.Errorf("couldn't save the data file: (%v)", err)
			}
		} else {
			return fmt.Errorf("the data file isn't valid gob: (%v)", err)
		}
	}

	return nil
}

func (w *Weather) GobEncode() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(w.Forecast)
	if err != nil {
		return nil, err
	}

	err = encoder.Encode(w.LastCheck)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (w *Weather) GobDecode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	err := decoder.Decode(&w.Forecast)
	if err != nil {
		return err
	}

	err = decoder.Decode(&w.LastCheck)
	if err != nil {
		return err
	}

	return nil
}
