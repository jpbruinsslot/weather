package colors

type Colors struct {
	Black   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	White   string
	Reset   string
}

// DefaultANSIColors is the default color scheme for the ANSI color
// escape codes.
var DefaultANSIColors = Colors{
	Black:   "\033[30m",
	Red:     "\033[31m",
	Green:   "\033[32m",
	Yellow:  "\033[33m",
	Blue:    "\033[34m",
	Magenta: "\033[35m",
	Cyan:    "\033[36m",
	White:   "\033[37m",
	Reset:   "\033[0m",
}
