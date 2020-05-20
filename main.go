package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	log "github.com/golang/glog"
	"github.com/mmcdole/gofeed"
)

func main() {
	cmdOpts := parseCmdOptions()
	config := readConfigurationFile(cmdOpts.ConfigFile)

	if cmdOpts.ValidateConfig {
		if err := validateConfig(config); err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Configuration file %s is valid\n", cmdOpts.ConfigFile)
		os.Exit(0)
	}
	for {
		executeFeeds(config)
		if !config.DeamonMode {
			break
		}
		time.Sleep(config.RefreshRate)
	}
}

func executeFeeds(config config) {
	for k, v := range config.Feeds {
		log.Infof("Downloading feed %s", k)
		populateDefaults(&v, config.commonConfig)
		processFeed(k, &v)
	}
}

func processFeed(feedName string, feed *feed) {
	fp := gofeed.NewParser()
	feedInfo, err := fp.ParseURL(feed.URL)

	if err != nil {
		log.Errorf("[%s] %s", feedName, err.Error())
		return
	}
	if len(feedInfo.Items) == 0 {
		log.Infof("[%s] no new items", feedName)
		return
	}
	for _, v := range feedInfo.Items {
		if v.PublishedParsed.Before(feed.lastUpdated) {
			log.Infof("[%s] skipping %s already downloaded", feedName, v.Title)
			continue
		}
		filename, url := extractAtom(*v)
		err := DownloadFile(filename, url)
		if os.IsExist(err) {
			log.Infof("[%s] already present %s", feedName, v.Title)
		} else if err != nil {
			log.Infof("[%s] failed to download %s. due to error %s", feedName, v.Title, err.Error())
		} else {
			log.Infof("[%s] downloaded %s", feedName, v.Title)
		}
	}
	lastFeed := feedInfo.Items[len(feedInfo.Items)-1]
	feed.lastUpdated = *lastFeed.PublishedParsed
}

func extractAtom(feedInfo gofeed.Item) (string, string) {
	url := feedInfo.GUID
	filename, err := FileNameFromContentDisposition(url)
	if err != nil {
		filename = feedInfo.Title
	}
	if match, _ := regexp.MatchString(`\.(torrent|magnet)$`, filename); !match {
		filename = filename + ".torrent"
	}
	return filename, url
}

func populateDefaults(feed *feed, commonConfig commonConfig) {
	if feed.DownloadDir == "" {
		feed.DownloadDir = commonConfig.DownloadDir
	}
	if feed.MinSize == 0 {
		feed.MinSize = commonConfig.MinSize
	}
	if feed.MaxSize == 0 {
		feed.MaxSize = commonConfig.MaxSize
	}
}
