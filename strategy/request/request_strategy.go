package strategy

import "net/http"

type RequestStrategy interface {
	// CloneWithDifferentPayload prepares a new strategy instance with the given payloads applied
	// to the base request template. It should store the resulting request string internally.
	// Returns an error if payload replacement fails.
	CloneWithDifferentPayload(idx int, baseReq string, payload []string, strategyClones *[]RequestStrategy) error

	// CreateRequest generates the final http.Request object based on the stored request string
	// and the provided path (extracted from the first line).
	CreateRequest(path string) (*http.Request, error)

	// ToString returns the final request string and the payload string used for this instance.
	ToString() (string, string)
}

// ChooseStrategy selects the appropriate request strategy based on the HTTP method.
// Panics if method is not supported (Consider returning error).
func ChooseStrategy(method string) RequestStrategy {
	switch method {
	case "GET":
		return &GetRequestStrategy{}
	case "POST":
		// Add other methods here as needed
		return &PostRequestStrategy{}
	// case "PUT":
	//  return &PutRequestStrategy{} // Example
	default:
		// Consider returning nil, fmt.Errorf("method '%s' not implemented", method)
		panic("Method not implemented: " + method)
	}
}