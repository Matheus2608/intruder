package strategy

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"
)

type PostRequestStrategy struct {
	Payload []string // Payloads used for *this* specific request instance
	HttpReq string   // The final HTTP request string *after* payload injection
}

// parseRequestAndBody separates headers and body from the request string (excluding first line).
func (strategy *PostRequestStrategy) parseRequestAndBody(lines []string) (map[string][]string, string, error) {
	headersMap := make(map[string][]string)
	bodySeparatorIndex := -1

	// Find headers and body separation
	for idx, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			// Found the separator between headers and body
			bodySeparatorIndex = idx
			break
		}

		// Parse header line
		parts := strings.SplitN(line, ":", 2) // Use original line for parsing
		if len(parts) != 2 {
			return nil, "", fmt.Errorf("malformed header line: %q", line)
		}
		headerName := strings.TrimSpace(parts[0])
		headerValue := strings.TrimSpace(parts[1])

		if headerName == "" {
			return nil, "", fmt.Errorf("empty header name found in line: %q", line)
		}
		// Check for remaining markers - should not happen
		if strings.Contains(headerValue, "§") || strings.Contains(headerName, "§") {
			log.Printf("Warning: Found payload marker '§' in parsed header: %s: %s", headerName, headerValue)
			// return nil, "", fmt.Errorf("unexpected payload marker '§' found after replacement in header: %s", headerName)
		}

		headersMap[headerName] = append(headersMap[headerName], headerValue)
	}

	// Extract body
	var body string
	if bodySeparatorIndex != -1 && bodySeparatorIndex+1 < len(lines) {
		// Body is everything after the separator line, preserving original newlines
		// Re-join using '\n' as we split by '\n' earlier
		body = strings.Join(lines[bodySeparatorIndex+1:], "\n")
	} else if bodySeparatorIndex == -1 {
		// No empty line found, assume no body? This is unusual for POST but possible.
		// Or maybe the entire part after line 1 is headers? Let's assume no body if no separator.
		log.Println("Warning: No empty line separator found between headers and body in POST request.")
		body = ""
        // Add all lines as headers if no separator? The loop above handles this.
	}


	return headersMap, body, nil
}

// CloneWithDifferentPayload creates a new instance with payloads applied.
func (strategy *PostRequestStrategy) CloneWithDifferentPayload(idx int, baseReq string, payload []string, strategyClones *[]RequestStrategy) error {
	newReqStr, err := ReplaceDynamicInput(baseReq, payload)
	if err != nil {
		return fmt.Errorf("failed to replace dynamic input for index %d: %w", idx, err)
	}

	cloneStrategy := &PostRequestStrategy{
		Payload: payload,
		HttpReq: newReqStr,
	}

	if idx >= 0 && idx < len(*strategyClones) {
		(*strategyClones)[idx] = cloneStrategy
	} else {
		return fmt.Errorf("internal error: index %d out of bounds for strategyClones slice", idx)
	}
	return nil
}

// CreateRequest builds the http.Request object.
func (strategy *PostRequestStrategy) CreateRequest(path string) (*http.Request, error) {
	lines := strings.Split(strings.ReplaceAll(strategy.HttpReq, "\r\n", "\n"), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("cannot create POST request: request string is empty")
	}
	// Content (headers and body) starts from the second line
	contentLines := []string{}
	if len(lines) > 1 {
		contentLines = lines[1:]
	}

	// Parse headers and body
	headersMap, body, err := strategy.parseRequestAndBody(contentLines)
	if err != nil {
		return nil, fmt.Errorf("cannot create POST request: failed to parse headers/body: %w", err)
	}

	// Get Host header
	hostValues, hostExists := headersMap["Host"]
	if !hostExists || len(hostValues) == 0 || hostValues[0] == "" {
		return nil, fmt.Errorf("cannot create POST request: 'Host' header is missing or empty")
	}
	host := hostValues[0]

	// Construct URL (similar to GET)
	scheme := "https"
	// Basic scheme check (Host header shouldn't really contain it)
	if strings.HasPrefix(host, "http://") { scheme = "http" }
    if !strings.HasPrefix(path, "/") { path = "/" + path }
	urlStr := fmt.Sprintf("%s://%s%s", scheme, host, path)


	// Create request object with body
    // Use strings.NewReader for the body, it's efficient
	httpReq, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request object: %w", err)
	}

	// Add headers
	// Track if Content-Length was set manually
	hasManualContentLength := false
	for name, values := range headersMap {
		// Handle Host header explicitly (same as GET)
		if strings.EqualFold(name, "Host") {
			if httpReq.URL.Host != values[0] {
				httpReq.Host = values[0]
			}
			continue
		}
		// Check for manual Content-Length
		if strings.EqualFold(name, "Content-Length") {
			hasManualContentLength = true
			// Use the value provided by the user
			// http.NewRequest automatically sets ContentLength based on the reader,
			// but we should respect the user's value if provided.
			// Setting Header["Content-Length"] manually overrides Go's calculation.
		}

		// Add header values
		for _, value := range values {
			httpReq.Header.Add(name, value)
		}
	}

    // If Content-Length wasn't manually set by the user in the template,
    // Go's http.NewRequest will have set it correctly based on the strings.Reader.
    // If it *was* set manually, Go respects the Header map value.
    // We might want to log a warning if the manual Content-Length doesn't match len(body).
    if clHeader, ok := httpReq.Header["Content-Length"]; ok && len(clHeader) > 0 {
         manualLenStr := clHeader[0]
         manualLen, convErr := strconv.Atoi(manualLenStr) // <--- USO DO STRCONV AQUI
         actualLen := len(body)
         if convErr != nil {
              log.Printf("Warning: Invalid manual Content-Length value '%s' in request.", manualLenStr)
         } else if hasManualContentLength && manualLen != actualLen {
             log.Printf("Warning: Manual Content-Length (%d) in request does not match actual body size (%d). Using manual value.", manualLen, actualLen)
         }
    }


	return httpReq, nil
}

// ToString returns the request string and payload string.
func (strategy *PostRequestStrategy) ToString() (string, string) {
	payloadStr := strings.Join(strategy.Payload, ", ")
	return strategy.HttpReq, payloadStr
}