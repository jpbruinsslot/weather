package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jpbruinsslot/weather"
	"github.com/jpbruinsslot/weather/config"
	"github.com/jpbruinsslot/weather/utils/xdg"
)

const (
	USAGE = `NAME:
    weather - weather forecasts from the command line

USAGE:
    weather [options]

VERSION:
    %s

WEBSITE:
    %s

OPTIONS:
    -config <path>  location of config file
    -config-path    print location of config file
    -debug          enable debug mode
    -help           show this help message
`
)

var (
	flgConfig string
	flgPath   bool
	flgDebug  bool
)

func init() {
	flag.StringVar(
		&flgConfig,
		"config",
		"",
		"location of config file",
	)

	flag.BoolVar(
		&flgPath,
		"config-path",
		false,
		"print location of config file",
	)

	flag.BoolVar(
		&flgDebug,
		"debug",
		false,
		"enable debug mode",
	)

	flag.Usage = func() {
		fmt.Printf(USAGE, weather.Version, weather.URL)
	}

	flag.Parse()
}

func main() {
	var err error

	// Locate config file
	if flgConfig == "" {
		cfgPath, err := xdg.ConfigFile("weather/weather.conf")
		if err != nil {
			log.Fatal(err)
		}
		flgConfig = cfgPath
	}

	if flgPath {
		fmt.Println(flgConfig)
		return
	}

	// Set up logger
	var opts *slog.HandlerOptions
	if flgDebug {
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)

	cfg, err := config.New(flgConfig)
	if err != nil {
		log.Fatal(err)
	}

	w, err := weather.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = w.GetForecast()
	if err != nil {
		log.Fatal(err)
	}

	err = w.PrintForecast()
	if err != nil {
		log.Fatal(err)
	}
}
