package strategy

import "net/http"

type RequestStrategy interface {
	CloneWithDifferentPayload(idx int, req string, payload []string, strategyClones *[]RequestStrategy)
	CreateRequest(path string) (*http.Request, error)
	ToString() (string, string)
}

func ChooseStrategy(method string) RequestStrategy {
	switch method {
	case "GET":
		return &GetRequestStrategy{}
	case "POST":
		return &PostRequestStrategy{}
	default:
		panic("Method not implemented")
	}
}
