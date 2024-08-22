package handlers

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func createRequest(req string) (*http.Request, []string, error) {
	requestMap, headersWhichNeedToBeChanged, path := parseRequest(req)

	host, hostExists := requestMap["Host"]
	if !hostExists {
		fmt.Println("Host key does not exist in requestMap")
		return nil, nil, fmt.Errorf("host key does not exist in requestMap")
	}

	// Ensure the URL is properly formatted
	url := "https://" + host + path

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err == nil {
		for headerType, value := range requestMap {
			httpReq.Header.Add(headerType, value)
		}
	}

	return httpReq, headersWhichNeedToBeChanged, err
}

func parseRequest(req string) (map[string]string, []string, string) {
	requestMap := make(map[string]string)
	headersWhichNeedToBeChanged := []string{}

	lines := strings.Split(req, "\n")
	path := strings.SplitN(lines[0], " ", 3)[1]

	for _, header := range lines[1:] {
		if len(header) < 3 {
			break
		}

		headerList := strings.SplitN(header, ": ", 2)
		if len(headerList) == 2 {
			typeHeader, value := headerList[0], headerList[1]
			if strings.Contains(value, "§") {
				// fmt.Println("Header contains payload:", typeHeader)
				headersWhichNeedToBeChanged = append(headersWhichNeedToBeChanged, typeHeader)
			}
			requestMap[typeHeader] = strings.TrimSpace(value)
		} else {
			fmt.Println("Invalid header format:", header)
		}
	}

	return requestMap, headersWhichNeedToBeChanged, path
}

func makeRequestWithPayload(client *http.Client, req *http.Request, payload string, headersWhichNeedToBeChanged []string) error {

	for _, header := range headersWhichNeedToBeChanged {
		// fmt.Println(header)
		headerChanged := changeHearderWithPayload(req.Header.Get(header), payload)
		// fmt.Println("Header changed:", headerChanged)
		req.Header.Set(header, headerChanged)
	}

	// Send the request
	httpRes, err := client.Do(req)

	// Close the body in all circustances
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}

	defer httpRes.Body.Close()

	// using switch to handle different content encodings
	var reader io.ReadCloser
	switch httpRes.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(httpRes.Body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return err
		}
	default:
		reader = httpRes.Body
	}

	// Read the response body
	_, err = io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// fmt.Printf("%s", body)

	return nil
}

func changeHearderWithPayload(header string, payload string) string {
	separetedString := strings.SplitN(header, "§", 3)
	return separetedString[0] + payload + separetedString[2]
}

func PostHandler(res http.ResponseWriter, req *http.Request) {
	// Parse the form data
	if err := req.ParseForm(); err != nil {
		http.Error(res, "Unable to parse form", http.StatusBadRequest)
		return
	}

	httpReq, headersWhichNeedToBeChanged, err := createRequest(req.Form["requestData"][0])
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	err = makeRequestWithPayload(client, httpReq, req.Form["payload"][0], headersWhichNeedToBeChanged)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
}
