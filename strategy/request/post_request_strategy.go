package strategy

import (
	"fmt"
	"net/http"
	"strings"
)

type PostRequestStrategy struct {
	Payload []string
	HttpReq string
}

func (strategy *PostRequestStrategy) parseRequest(headerAndBody []string) (map[string]string, string, error) {
	headersMap := make(map[string]string)
	var idxBodyStarts int

	for idx, header := range headerAndBody {
		if header == "" {
			idxBodyStarts = idx + 1
			// fmt.Println("idxBodyStarts", idxBodyStarts)
			break
		}

		headerList := strings.SplitN(header, ": ", 2)

		if len(headerList) < 2 {
			return nil, "", fmt.Errorf("HTTP is malformed")
		}

		typeHeader, value := headerList[0], headerList[1]

		if strings.Contains(value, "§") {
			return nil, "", fmt.Errorf("there are still places where dynamic inputs need to be replaced")
		}
		headersMap[typeHeader] = strings.TrimSpace(value)
	}

	body := strings.Join(headerAndBody[idxBodyStarts:], "\n")

	return headersMap, body, nil
}

func (strategy *PostRequestStrategy) CloneWithDifferentPayload(idx int, req string, payload []string, strategyClones *[]RequestStrategy) {
	newReq, err := ReplaceDynamicInput(req, payload)
	if err != nil {
		panic(err)
	}

	cloneStrategy := new(PostRequestStrategy)
	*cloneStrategy = *strategy

	cloneStrategy.Payload = payload
	cloneStrategy.HttpReq = newReq

	(*strategyClones)[idx] = cloneStrategy
}

func (strategy *PostRequestStrategy) CreateRequest(path string) (*http.Request, error) {
	headersAndBody := strings.Split(strategy.HttpReq, "\r\n")[1:]
	headersMap, body, err := strategy.parseRequest(headersAndBody)
	if err != nil {
		return nil, err
	}

	host, hostExists := headersMap["Host"]
	if !hostExists {
		return nil, fmt.Errorf("host key does not exist in requestMap")
	}

	// Ensure the URL is properly formatted
	url := "https://" + host + path
	httpReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))

	if err != nil {
		return httpReq, err
	}

	for headerType, value := range headersMap {
		httpReq.Header.Add(headerType, value)
	}

	return httpReq, err
}

func (strategy *PostRequestStrategy) ToString() (string, string) {
	return strategy.HttpReq, strings.Join(strategy.Payload, ", ")
}
