package xdg

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func ConfigFile(name string) (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")

	// When XDG_CONFIG_HOME is not set
	if configHome == "" {
		switch runtime.GOOS {
		case "windows":
			configHome = filepath.Join(os.Getenv("APPDATA"))
		case "darwin":
			configHome = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
		case "linux":
			configHome = filepath.Join(os.Getenv("HOME"), ".config")
		default:
			return "", fmt.Errorf("unsupported operating system")
		}
	}

	// Check if file exists
	if _, err := os.Stat(filepath.Join(configHome, name)); err != nil {
		dir := filepath.Dir(filepath.Join(configHome, name))

		// Create directory
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", fmt.Errorf("couldn't create the config directory: (%v)", err)
		}

		// Create file
		if _, err := os.Create(filepath.Join(configHome, name)); err != nil {
			return "", fmt.Errorf("couldn't create the config file: (%v)", err)
		}
	}

	return filepath.Join(configHome, name), nil
}

func DataFile(name string) (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")

	// When XDG_DATA_HOME is not set
	if dataHome == "" {
		switch runtime.GOOS {
		case "windows":
			dataHome = filepath.Join(os.Getenv("APPDATA"))
		case "darwin":
			dataHome = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
		case "linux":
			dataHome = filepath.Join(os.Getenv("HOME"), ".local", "share")
		default:
			return "", fmt.Errorf("unsupported operating system")
		}
	}

	// Check if file exists
	if _, err := os.Stat(filepath.Join(dataHome, name)); err != nil {
		dir := filepath.Dir(filepath.Join(dataHome, name))

		// Create directories
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", fmt.Errorf("couldn't create the data directory: (%v)", err)
		}

		// Create file
		if _, err := os.Create(filepath.Join(dataHome, name)); err != nil {
			return "", fmt.Errorf("couldn't create the data file: (%v)", err)
		}
	}

	return filepath.Join(dataHome, name), nil
}
