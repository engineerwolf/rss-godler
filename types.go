package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type config struct {
	DeamonMode   bool            `yaml:"deamon-mode"`
	RefreshRate  time.Duration   `yaml:"refresh,omitempty"`
	Feeds        map[string]feed `yaml:"feeds,flow"`
	commonConfig `yaml:",inline"`
}

type feed struct {
	URL          string `yaml:"url"`
	REI          string `yaml:"include,omitempty"`
	REX          string `yaml:"exclude,omitempty"`
	commonConfig `yaml:",inline"`
	lastUpdated  time.Time
}

type commonConfig struct {
	DownloadDir string   `yaml:"download-dir,omitempty"`
	MinSize     FileSize `yaml:"min-size,omitempty"`
	MaxSize     FileSize `yaml:"max-size,omitempty"`
}

type cmdOpts struct {
	ConfigFile     string `short:"c" long:"conf" description:"configuration file path" default:"$HOME/.config/rss-godler"`
	ValidateConfig bool   `long:"validate-config" description:"validate configuration file and exit."`
}

type feedTracker struct {
	Feeds map[string]trackingInfo `yaml:"feeds,flow"`
}

type trackingInfo struct {
	lastUpdated time.Time `yaml:"lastUpdated"`
}

const (
	//Byte 8 bits
	Byte FileSize = 1
	//Kibibyte 1024 Bytes
	Kibibyte = 1024 * Byte
	// Mebibyte 1024 KiloBytes
	Mebibyte = 1024 * Kibibyte
	// Gibibyte 1024 Mebibytes
	Gibibyte = 1024 * Mebibyte
	//Tebibyte 1024 Gibibyte
	Tebibyte = 1024 * Gibibyte

	//Kilobyte 1000 Bytes
	Kilobyte = 1000 * Byte
	// Megabyte 1000 KiloBytes
	Megabyte = 1000 * Kilobyte
	// Gigabyte 1000 Megabyte
	Gigabyte = 1000 * Megabyte
	//Terabyte 1000 Gigabyte
	Terabyte = 1000 * Gigabyte
)

// FileSize represents number of bytes
type FileSize int64

// UnmarshalYAML Implements the Unmarshaler interface of the yaml pkg.
func (fs *FileSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawValue string
	err := unmarshal(&rawValue)
	if err != nil {
		return err
	}
	if rawValue == "None" || rawValue == "" {
		*fs = FileSize(-1)
		return nil
	}
	sizeRegEx := regexp.MustCompile(`^(?P<value>[[:digit:]]+)[ ]?(?P<unit>(b|B|Kb|Kib|Mb|Mib|Gb|Gib|Tb|Tib)?)$`)
	if !sizeRegEx.MatchString(rawValue) {
		return fmt.Errorf("value must be digits optionally followed by one of the units (b, B, Kb, Kib, Mb, Mib, Gb, Gib, Tb, Tib) without any spaces")
	}
	unit := sizeRegEx.ReplaceAllString(rawValue, "${unit}")
	value, err := strconv.ParseFloat(sizeRegEx.ReplaceAllString(rawValue, "${value}"), 64)
	if err != nil {
		return err
	}
	bytes := func(unit string, value float64) float64 {
		switch unit {
		case "Kb":
			return value * float64(Kilobyte)
		case "Kib":
			return value * float64(Kibibyte)
		case "Mb":
			return value * float64(Megabyte)
		case "Mib":
			return value * float64(Mebibyte)
		case "Gb":
			return value * float64(Gigabyte)
		case "Gib":
			return value * float64(Gibibyte)
		case "Tb":
			return value * float64(Terabyte)
		case "Tib":
			return value * float64(Tebibyte)
		case "b":
			fallthrough
		case "B":
			fallthrough
		default:
			return value
		}
	}(unit, value)

	*fs = FileSize(bytes)
	return nil
}
