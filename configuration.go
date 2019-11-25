package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	log "github.com/golang/glog"
	flags "github.com/jessevdk/go-flags"
	yaml "gopkg.in/yaml.v2"
)

var (
	sizeRegEx = regexp.MustCompile("^(?P<value>[[:digit:]]+)(?P<unit>[B|K|M|G|T]?)$")
	version   = "undefined"
)

func parseCmdOptions() (cmdOpts cmdOpts) {
	_, err := flags.Parse(&cmdOpts)

	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			os.Exit(1)
		} else {
			log.Fatalf(err.Error())
		}
	}
	if cmdOpts.PrintVersion {
		fmt.Printf("rss-godler v.%s\n", version)
		os.Exit(0)
	}
	return
}

func readConfigurationFile(filename string) (config config) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to open file %s", filename)
	}
	err2 := yaml.Unmarshal(data, &config)
	if err2 != nil {
		log.Fatalf("cannot unmarshal data: %v", err2)
	}
	setDefaults(&config)
	return
}

func setDefaults(config *config) {
	if config.DownloadDir == "" {
		config.DownloadDir = "./"
	}
	if config.RefreshRate < 0 {
		config.RefreshRate = 15 * time.Minute
	}
	if config.RefreshRate < 1*time.Minute {
		config.RefreshRate = 15 * time.Minute
	}
	if config.MaxSize == 0 {
		config.MaxSize = -1
	}
}

func validateConfig(config config) error {

	for k, v := range config.Feeds {
		if v.REI != "" {
			if _, err := regexp.Compile(v.REI); err != nil {
				return fmt.Errorf("Include expression for %s is invalid. %s", k, v.REI)
			}
		}
		if v.REX != "" {
			if _, err := regexp.Compile(v.REX); err != nil {
				return fmt.Errorf("Exclude expression for %s is invalid. %s", k, v.REX)
			}
		}
	}
	return nil
}

func printConfig(config config) {
	fmt.Println("deamon-mode:", config.DeamonMode)
	fmt.Println("download-dir:", config.DownloadDir)
	fmt.Println("min-size:", config.MinSize)
	fmt.Println("max-size:", config.MaxSize)
	fmt.Println("refresh:", config.RefreshRate)
	for k, v := range config.Feeds {
		fmt.Println("For", k)
		fmt.Println("\tURL:", v.URL)
		fmt.Println("\tdownload-dir:", v.DownloadDir)
		fmt.Println("\tinclude:", v.REI)
		fmt.Println("\texclude:", v.REX)
		fmt.Println("\tmin-size:", v.MinSize)
		fmt.Println("\tmax-size:", v.MaxSize)
	}
}
