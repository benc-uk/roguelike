package core

import (
	"io"
	"net/http"
	"os"
	"runtime"
)

func IsWASM() bool {
	return runtime.GOARCH == "wasm"
}

func FetchURL(url string) ([]byte, error) {
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

func ReadFile(filePath string) (bytes []byte, err error) {
	var data []byte

	if IsWASM() {
		data, err = FetchURL(filePath)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
