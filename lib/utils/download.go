package utils

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/api/drive/v2"
)

func DownloadFile(f *drive.File) ([]byte, error) {
	downloadUrl := f.DownloadUrl
	if downloadUrl == "" {
		// If there is no downloadUrl, there is no body
		log.Printf("An error occurred: File is not downloadable")
		return nil, errors.New("an error occurred: File is not downloadable")
	}

	resp, err := http.Get(downloadUrl)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}

	// // Make sure we close the Body later
	defer resp.Body.Close()
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
		return nil, err
	}

	return body, nil
}
