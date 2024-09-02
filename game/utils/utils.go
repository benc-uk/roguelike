package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

func IsWASM() bool {
	host, err := os.Hostname()
	if err != nil {
		return false
	}

	return host == "js"
}

func FetchURL(url string) ([]byte, error) {
	log.Println("Fetching URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
