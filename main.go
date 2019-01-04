package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
		downloadFeeds(config)
		if !config.DeamonMode {
			break
		}
		time.Sleep(config.RefreshRate)
	}
}

func downloadFeeds(config config) {
	for k, v := range config.Feeds {
		log.Printf("Downloading feed %s", k)
		populateDefaults(&v, config.commonConfig)
		downloadFeed(k, &v)
	}
}

func downloadFeed(feedName string, feed *feed) {
	fp := gofeed.NewParser()
	feedInfo, err := fp.ParseURL(feed.URL)

	if err != nil {
		log.Printf("[%s] could not retrive feed. due to error %s", feedName, err.Error())
		return
	}
	if len(feedInfo.Items) == 0 {
		log.Printf("[%s] no new items", feedName)
		return
	}
	for _, v := range feedInfo.Items {
		if v.PublishedParsed.Before(feed.lastUpdated) {
			log.Printf("[%s] skipping %s already downloaded", feedName, v.Title)
			continue
		}
		err := downloadAtom(*feed, *v)
		if os.IsExist(err) {
			log.Printf("[%s] already present %s", feedName, v.Title)
		} else if err != nil {
			log.Printf("[%s] failed to download %s. due to error %s", feedName, v.Title, err.Error())
		} else {
			log.Printf("[%s] downloaded %s", feedName, v.Title)
		}
	}
	lastFeed := feedInfo.Items[len(feedInfo.Items)-1]
	feed.lastUpdated = *lastFeed.PublishedParsed
}

func downloadAtom(feed feed, feedInfo gofeed.Item) error {
	url := feedInfo.GUID
	filename, err := fileNameFromContentDisposition(url)
	if err != nil {
		filename = feedInfo.Title
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filePath := filepath.Join(os.ExpandEnv(feed.DownloadDir), filename)

	out, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	// make request on link to mark as read
	resp, _ = http.Get(feedInfo.Link)
	resp.Body.Close()
	return nil
}

func fileNameFromContentDisposition(url string) (string, error) {
	resp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		return "", err
	}
	return params["filename"], nil
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
