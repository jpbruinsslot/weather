package weather

import (
	"fmt"
	"log/slog"

	"github.com/jpbruinsslot/weather/icons"
	"github.com/jpbruinsslot/weather/utils/colors"
	"github.com/jpbruinsslot/weather/utils/conv"
)

type Output struct {
	Colors      colors.Colors
	Units       string
	Temperature float64
	Rain        string
	Icon        string
}

func (w *Weather) Temperature() float64 {
	return w.Forecast.Temperature
}

func (w *Weather) Units() string {
	if w.Config.Units == conv.Imperial {
		return w.Config.Icons["fahrenheit"]
	}
	return w.Config.Icons["celsius"]
}

func (w *Weather) Icon() string {
	slog.Debug(
		fmt.Sprintf("icon used: %s", w.Forecast.Icon),
	)
	return w.Config.Icons[w.Forecast.Icon]
}

func (w *Weather) Rain() string {
	if !w.Forecast.Rain {
		return ""
	}

	return w.Config.Icons[icons.RainIndicator]
}

func (w *Weather) Colors() colors.Colors {
	return w.Config.Colors
}

func (w *Weather) GenerateOutput() Output {
	return Output{
		Temperature: w.Temperature(),
		Rain:        w.Rain(),
		Icon:        w.Icon(),
		Units:       w.Units(),
		Colors:      w.Colors(),
	}
}
