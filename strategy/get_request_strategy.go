package strategy

import (
	"fmt"
	"net/http"
	"strings"
)

type GetRequestStrategy struct {
	Payload []string
	HttpReq string
}

func parseRequest(headerAndBody []string) (map[string]string, error) {
	headersMap := make(map[string]string)

	for _, header := range headerAndBody {
		if header == "" {
			break
		}

		headerList := strings.SplitN(header, ": ", 2)

		if len(headerList) < 2 {
			return nil, fmt.Errorf("HTTP is malformed")
		}

		typeHeader, value := headerList[0], headerList[1]

		if strings.Contains(header, "§") {
			return nil, fmt.Errorf("there are still places where dynamic inputs need to be replaced")
		}

		headersMap[typeHeader] = strings.TrimSpace(value)
	}

	return headersMap, nil
}

func (strategy *GetRequestStrategy) CloneWithDifferentPayload(idx int, req string, payload []string, strategyClones *[]RequestStrategy) {
	newReq, err := ReplaceDynamicInput(req, payload)
	fmt.Print("newReq: ", newReq)
	if err != nil {
		panic(err)
	}

	cloneStrategy := new(GetRequestStrategy)
	*cloneStrategy = *strategy

	cloneStrategy.Payload = payload
	cloneStrategy.HttpReq = newReq

	(*strategyClones)[idx] = cloneStrategy
}

// TODO olhar se todos os mapas vao mudar
func (strategy *GetRequestStrategy) CreateRequest(path string) (*http.Request, error) {
	headersAndBody := strings.Split(strategy.HttpReq, "\r\n")[1:]
	headersMap, err := parseRequest(headersAndBody)
	if err != nil {
		return nil, err
	}

	host, hostExists := headersMap["Host"]

	if !hostExists {
		return nil, fmt.Errorf("host key does not exist in requestMap")
	}

	// Ensure the URL is properly formatted
	url := "https://" + host + path

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return httpReq, err
	}

	for headerType, value := range headersMap {
		httpReq.Header.Add(headerType, value)
	}

	return httpReq, err
}

func (strategy *GetRequestStrategy) ToString() (string, string) {
	return strategy.HttpReq, strings.Join(strategy.Payload, ", ")
}
