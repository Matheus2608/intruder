package facades

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

//
// Actions
//

func CreateRequest(path string, requestMap map[string]string) (*http.Request, error) {

	host, hostExists := requestMap["Host"]
	if !hostExists {
		fmt.Println("Host key does not exist in requestMap")
		return nil, fmt.Errorf("host key does not exist in requestMap")
	}

	// Ensure the URL is properly formatted
	url := "https://" + host + path

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		for headerType, value := range requestMap {
			httpReq.Header.Add(headerType, value)
		}
	}

	return httpReq, err
}

func changeHearderWithPayload(header string, payload string) string {
	separetedString := strings.SplitN(header, "§", 3)
	return separetedString[0] + payload + separetedString[2]
}

// TODO make headersWhichNeedToBeChanged a global variable
func SendRequestWithPayload(client *http.Client, req *http.Request, payload string, headersWhichNeedToBeChanged []string) (*http.Response, time.Duration, error) {

	for _, header := range headersWhichNeedToBeChanged {
		originalHeader := req.Header.Get(header)
		// fmt.Println("Original header:", originalHeader)
		headerChanged := changeHearderWithPayload(originalHeader, payload)
		// fmt.Println("Header changed:", headerChanged)
		req.Header.Set(header, headerChanged)
	}

	startTime := time.Now() // Captura o tempo de início

	httpRes, err := client.Do(req) // Send the request

	endTime := time.Now()                 // Captura o tempo de fim
	elapsedTime := endTime.Sub(startTime) // Calcula o tempo decorrido

	return httpRes, elapsedTime, err
}

//
// Parses
//

func ParseRequest(req string) (map[string]string, []string, string) {
	requestMap := make(map[string]string)
	headersWhichNeedToBeChanged := []string{}

	lines := strings.Split(req, "\n")
	path := strings.SplitN(lines[0], " ", 3)[1]

	for _, header := range lines[1:] {
		if len(header) < 3 {
			break
		}

		headerList := strings.SplitN(header, ": ", 2)
		typeHeader, value := headerList[0], headerList[1]
		if strings.Contains(value, "§") {
			// fmt.Println("Header contains payload:", typeHeader)
			headersWhichNeedToBeChanged = append(headersWhichNeedToBeChanged, typeHeader)
		}
		requestMap[typeHeader] = strings.TrimSpace(value)
	}

	return requestMap, headersWhichNeedToBeChanged, path
}

func ParseBody(body io.ReadCloser, contentEncoding string) (string, error) {
	if body == nil {
		return "", fmt.Errorf("body is nil")
	}
	defer body.Close() // Garantir que o corpo seja fechado

	var reader io.ReadCloser
	var err error

	if contentEncoding == "gzip" {
		reader, err = gzip.NewReader(body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return "", err
		}
		defer reader.Close() // Garantir que o reader seja fechado
	} else {
		reader = body
	}

	// Ler o corpo
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
