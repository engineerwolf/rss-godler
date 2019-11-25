package main

import (
	"io"
	"mime"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	log.Infof("Downloading %s to %s\n", url, filepath)
	client := &http.Client{}
	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	log.Infof("Download %s complete\n", filepath)
	return err
}

// FileNameFromContentDisposition extracts filename from content disposition
func FileNameFromContentDisposition(url string) (string, error) {
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
