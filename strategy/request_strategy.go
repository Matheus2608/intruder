package strategy

import "net/http"

type RequestStrategy interface {
	CloneWithDifferentPayload(idx int, req string, payload string, strategyClones *[]RequestStrategy)
	CreateRequest(path string) (*http.Request, error)
	ToString() (string, string)
}
