package handlers

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type ResponseData struct {
	RequestId   uint16
	Payload     string // payload user suplied
	StatusCode  uint8
	TimeElapsed uint16 // The time taken to receive or complete the response
	Err         bool
	Length      uint32
	HttpReq     string
	HttpRes     string
}

type ResponseList struct {
	Responses []ResponseData
}

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

func readResBody(httpRes *http.Response) (string, error) {
	if httpRes.Body == nil {
		return "", fmt.Errorf("response body is nil")
	}
	defer httpRes.Body.Close() // Garantir que o corpo da resposta seja fechado

	var reader io.ReadCloser
	var err error

	if httpRes.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(httpRes.Body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return "", err
		}
		defer reader.Close() // Garantir que o reader seja fechado
	} else {
		reader = httpRes.Body
	}

	// Ler o corpo da resposta
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func readReqBody(httpReq *http.Request) (string, error) {
	if httpReq.Body == nil {
		return "", fmt.Errorf("request body is nil")
	}
	defer httpReq.Body.Close() // Garantir que o corpo da requisição seja fechado

	var reader io.ReadCloser
	var err error

	if httpReq.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(httpReq.Body)
		if err != nil {
			fmt.Println("Error creating gzip reader:", err)
			return "", err
		}
		defer reader.Close() // Garantir que o reader seja fechado
	} else {
		reader = httpReq.Body
	}

	// Ler o corpo da requisição
	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func populateResponse(httpReq *http.Request, httpRes *http.Response, payload string, idx int) ResponseData {
	requestBody, _ := readReqBody(httpReq)
	// if reqErr != nil {
	// 	fmt.Println("Error reading request body:", reqErr)
	// }

	responseBody, _ := readResBody(httpRes)
	// if resErr != nil {
	// 	fmt.Println("Error reading response body:", resErr)
	// }

	httpReqString := makehttpString(httpReq.Header, requestBody)
	httpResString := makehttpString(httpRes.Header, responseBody)

	return ResponseData{
		RequestId:   uint16(idx),
		Payload:     payload,
		StatusCode:  uint8(httpRes.StatusCode),
		TimeElapsed: 2,
		Err:         false,
		Length:      uint32(httpRes.ContentLength),
		HttpReq:     httpReqString,
		HttpRes:     httpResString,
	}
}

func makehttpString(headers map[string][]string, body string) string {
	httpResString := ""
	for headerType, value := range headers {
		httpResString += headerType + ": " + value[0] + "\n"
	}
	httpResString += "\n" + body
	return httpResString
}

// TODO make headersWhichNeedToBeChanged a global variable
func sendRequestWithPayload(client *http.Client, req *http.Request, payload string, headersWhichNeedToBeChanged []string) (*http.Response, error) {

	for _, header := range headersWhichNeedToBeChanged {
		// fmt.Println(header)
		headerChanged := changeHearderWithPayload(req.Header.Get(header), payload)
		// fmt.Println("Header changed:", headerChanged)
		req.Header.Set(header, headerChanged)
	}

	// Send the request
	httpRes, err := client.Do(req)

	return httpRes, err
}

func changeHearderWithPayload(header string, payload string) string {
	separetedString := strings.SplitN(header, "§", 3)
	return separetedString[0] + payload + separetedString[2]
}

func AttackHandler(res http.ResponseWriter, req *http.Request) {
	// Parse the form data
	if err := req.ParseForm(); err != nil {
		http.Error(res, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// create the request once, because only some header need to be changed later
	// this will save time and space
	httpReq, headersWhichNeedToBeChanged, err := createRequest(req.FormValue("requestData"))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	responseList := ResponseList{}
	for idx, payload := range strings.Split(req.FormValue("payload"), "\n") {
		httpRes, err := sendRequestWithPayload(client, httpReq, payload, headersWhichNeedToBeChanged)
		if err != nil {
			fmt.Println("Error sending request:", err)
		}
		response := populateResponse(httpReq, httpRes, payload, idx)
		responseList.Responses = append(responseList.Responses, response)
	}

	tmp, err := template.ParseFiles("templates/attack.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
	}
	err = tmp.Execute(res, responseList)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}
