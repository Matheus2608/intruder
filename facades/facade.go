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

func cloneHeaderOnly(r *http.Request) *http.Request {
	if r.Header == nil {
		panic("http: Request.Header is nil")
	}
	r2 := new(http.Request)
	*r2 = *r

	r2.Header = r.Header.Clone()
	return r2
}

func changeHearderWithPayload(header string, payload string) string {
	separetedString := strings.SplitN(header, "§", 3)
	return separetedString[0] + payload + separetedString[2]
}

func ChangeHeader(req *http.Request, headersToBeChanged []string, payload string) *http.Request {
	newHttpReq := cloneHeaderOnly(req)

	for _, header := range headersToBeChanged {
		originalHeader := req.Header.Get(header)
		// fmt.Println("Original header:", originalHeader)
		headerChanged := changeHearderWithPayload(originalHeader, payload)
		// fmt.Println("Header changed:", headerChanged)
		newHttpReq.Header.Set(header, headerChanged)
	}

	return newHttpReq
}

func SendRequest(client *http.Client, req *http.Request) (*http.Response, time.Duration, error) {

	startTime := time.Now()

	httpRes, err := client.Do(req)

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	return httpRes, elapsedTime, err
}

//
// Parses
//

func ParseRequest(req string) (map[string]string, []string, string) {
	requestMap := make(map[string]string)
	headersToBeChanged := []string{}

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
			headersToBeChanged = append(headersToBeChanged, typeHeader)
		}
		requestMap[typeHeader] = strings.TrimSpace(value)
	}

	return requestMap, headersToBeChanged, path
}

func ParseBody(body io.ReadCloser, contentEncoding string) (string, error) {
	if body == nil {
		return "", fmt.Errorf("body is nil")
	}

	defer body.Close()

	var reader io.ReadCloser
	var err error

	if contentEncoding == "gzip" {
		reader, err = gzip.NewReader(body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return "", err
		}
	} else {
		reader = body
	}

	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
