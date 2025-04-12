package strategy

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type GetRequestStrategy struct {
	Payload []string // Payloads used for *this* specific request instance
	HttpReq string   // The final HTTP request string *after* payload injection
}

// parseHeaders parses headers from the request string (excluding the first line and body).
// It assumes headers end at the first empty line.
func parseHeaders(headerLines []string) (map[string][]string, error) {
	headersMap := make(map[string][]string)
	for _, headerLine := range headerLines {
		headerLine = strings.TrimSpace(headerLine)
		if headerLine == "" {
			break // End of headers
		}

		parts := strings.SplitN(headerLine, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed header line: %q", headerLine)
		}

		headerName := strings.TrimSpace(parts[0])
		headerValue := strings.TrimSpace(parts[1])

		if headerName == "" {
			return nil, fmt.Errorf("empty header name found in line: %q", headerLine)
		}

		// Check for remaining markers - should not happen if ReplaceDynamicInput worked
		if strings.Contains(headerValue, "§") || strings.Contains(headerName, "§") {
			log.Printf("Warning: Found payload marker '§' in parsed header: %s: %s", headerName, headerValue)
			// Return error? Or proceed? Let's proceed but log warning.
			// return nil, fmt.Errorf("unexpected payload marker '§' found after replacement in header: %s", headerName)
		}

		// Use Header.Add equivalent logic (append to existing values)
		headersMap[headerName] = append(headersMap[headerName], headerValue)
	}
	return headersMap, nil
}

// CloneWithDifferentPayload creates a new instance with payloads applied.
func (strategy *GetRequestStrategy) CloneWithDifferentPayload(idx int, baseReq string, payload []string, strategyClones *[]RequestStrategy) error {
	newReqStr, err := ReplaceDynamicInput(baseReq, payload)
	if err != nil {
		// Don't store the clone if replacement fails
		return fmt.Errorf("failed to replace dynamic input for index %d: %w", idx, err)
	}

	// Create a new instance for the clone
	cloneStrategy := &GetRequestStrategy{
		Payload: payload,
		HttpReq: newReqStr,
	}

	// Store the pointer to the new instance in the pre-allocated slice
	if idx >= 0 && idx < len(*strategyClones) {
		(*strategyClones)[idx] = cloneStrategy
	} else {
		// This should not happen if slices are sized correctly
		return fmt.Errorf("internal error: index %d out of bounds for strategyClones slice", idx)
	}
	return nil
}

// CreateRequest builds the http.Request object.
func (strategy *GetRequestStrategy) CreateRequest(path string) (*http.Request, error) {
	// Split request string into lines
	lines := strings.Split(strings.ReplaceAll(strategy.HttpReq, "\r\n", "\n"), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("cannot create GET request: request string is empty")
	}
	// Headers start from the second line
	headerLines := []string{}
	if len(lines) > 1 {
		headerLines = lines[1:]
	}

	// Parse headers
	headersMap, err := parseHeaders(headerLines)
	if err != nil {
		return nil, fmt.Errorf("cannot create GET request: failed to parse headers: %w", err)
	}

	// Get Host header
	hostValues, hostExists := headersMap["Host"]
	if !hostExists || len(hostValues) == 0 || hostValues[0] == "" {
		// Try finding host in the request line if necessary? Usually it's a header.
		return nil, fmt.Errorf("cannot create GET request: 'Host' header is missing or empty")
	}
	host := hostValues[0] // Use the first Host header value

	// Construct URL - Assume HTTPS for now, could be configurable
	// Consider parsing the host to see if it includes a scheme already
	scheme := "https"
	if strings.HasPrefix(host, "http://") {
		scheme = "http"
		// host = strings.TrimPrefix(host, "http://") // Host header shouldn't contain scheme usually
	} else if strings.HasPrefix(host, "https://") {
		scheme = "https"
		// host = strings.TrimPrefix(host, "https://")
	}
	// Ensure path starts with '/'
	if !strings.HasPrefix(path, "/") {
        log.Printf("Warning: Path '%s' does not start with '/'. Prepending '/'.", path)
		path = "/" + path
	}

	// Use net/http standard library way to build URL to handle encoding etc.
    urlStr := fmt.Sprintf("%s://%s%s", scheme, host, path)
    // Basic validation if URL is parsable, NewRequest will do more thorough checks
    // _, urlErr := url.Parse(urlStr)
    // if urlErr != nil {
    //     return nil, fmt.Errorf("cannot create GET request: constructed URL '%s' is invalid: %w", urlStr, urlErr)
    // }


	// Create the request object (nil body for GET)
	httpReq, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request object: %w", err)
	}

	// Add headers (handle multi-value headers correctly)
	for name, values := range headersMap {
		// Special case: Host header is set automatically by NewRequest based on URL.
		// Setting it manually via Header.Add can cause issues or might be ignored.
		// The Go http client sets the Host header from req.URL.Host.
		// However, we *might* want to override it if the user explicitly set a *different* Host header.
		// Let's trust Go's default handling for now, unless 'Host' is explicitly different from URL host.
		if strings.EqualFold(name, "Host") {
			// If the parsed Host header value is different from the URL's host, set req.Host explicitly.
			if httpReq.URL.Host != values[0] {
				httpReq.Host = values[0] // Explicitly set the Host field for the HTTP request Host header
				// log.Printf("Explicitly setting Host header to: %s (URL host was: %s)", httpReq.Host, httpReq.URL.Host)
			}
			continue // Skip adding Host via Header.Add
		}

		// Add other headers
		for _, value := range values {
			httpReq.Header.Add(name, value)
		}
	}

	return httpReq, nil
}

// ToString returns the request string and payload string.
func (strategy *GetRequestStrategy) ToString() (string, string) {
	payloadStr := strings.Join(strategy.Payload, ", ") // Simple comma join
	return strategy.HttpReq, payloadStr
}