package strategy

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TODO make with more payloads
func ReplaceDynamicInput(req string, payload string) (string, error) {
	separetedString := strings.Split(req, "§")
	if len(separetedString) != 3 {
		return "", fmt.Errorf("invalid request")
	}
	return separetedString[0] + payload + separetedString[2], nil
}

func SendRequest(client *http.Client, req *http.Request) (*http.Response, time.Duration, error) {

	startTime := time.Now()

	httpRes, err := client.Do(req)

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	return httpRes, elapsedTime, err
}

func readBody(body io.ReadCloser) (string, error) {
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func ParseBody(body io.ReadCloser, contentEncoding string) (string, error) {
	if body == nil {
		return "", fmt.Errorf("body is nil")
	}

	if contentEncoding == "gzip" {
		reader, err := gzip.NewReader(body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return "", err
		}
		return readBody(reader)

	}

	return readBody(body)
}
