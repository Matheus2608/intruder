package strategy

import (
	// "bytes" //
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// ReplaceDynamicInput replaces §...§ markers in the request string with payloads.
// It iterates through the payload slice, replacing one marker pair per payload item.
func ReplaceDynamicInput(req string, payload []string) (string, error) {
	currentReq := req
	payloadIndex := 0

	// Keep replacing until no more § markers are found or we run out of payloads for them
	for strings.Contains(currentReq, "§") {
		// Ensure we have a payload for the next marker (unless it's an empty marker pair like §§)
		if payloadIndex >= len(payload) {
			// This happens in Sniper mode where some positions get empty strings.
			// Or if the number of markers doesn't match payload length (which matrix strategies should prevent).
			// Let's assume empty string should be used for remaining markers if payload slice is shorter.
            // A better approach might depend on the attack strategy context.
            // For now, let's replace remaining markers with empty string.
            // log.Printf("Warning: More markers than payloads. Replacing remaining '§...§' with empty string.")
		}

		payloadValue := "" // Default to empty string if no more payloads or payload is empty
		if payloadIndex < len(payload) {
			payloadValue = payload[payloadIndex]
		}


		parts := strings.SplitN(currentReq, "§", 3)
		if len(parts) != 3 {
			// Odd number of markers or malformed input
			return "", fmt.Errorf("invalid request template: unbalanced or adjacent '§' markers near '%s'", parts[0][max(0, len(parts[0])-20):])
		}

		// Replace the marker pair and the content between them
		currentReq = parts[0] + payloadValue + parts[2]
		payloadIndex++ // Move to the next payload for the next marker pair
	}

    // If payloadIndex < len(payload), it means we had more payloads than markers. This is usually fine.

	return currentReq, nil
}


// SendRequest sends the HTTP request and measures the time taken.
// Returns the response, elapsed duration, and any error during sending.
func SendRequest(client *http.Client, req *http.Request) (*http.Response, time.Duration, error) {
	if req == nil {
		return nil, 0, fmt.Errorf("cannot send nil request")
	}

	// Log the request details (optional, can be verbose)
	// log.Printf("Sending request: %s %s", req.Method, req.URL)

	startTime := time.Now()
	httpRes, err := client.Do(req)
	elapsedTime := time.Since(startTime)

	if err != nil {
		// Don't wrap error here, let caller handle specific error types if needed
		return nil, elapsedTime, err // Return elapsed time even on error
	}

	// Log response status (optional)
	// log.Printf("Received response: %s (took %v)", httpRes.Status, elapsedTime)

	return httpRes, elapsedTime, nil
}

// readBody is a helper to read all from a ReadCloser and ensure it's closed.
func readBody(reader io.ReadCloser) ([]byte, error) {
	if reader == nil {
		return nil, nil // Nothing to read
	}
	defer reader.Close() // Ensure closure

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed reading body: %w", err)
	}
	return bodyBytes, nil
}

// ParseBody reads the response body, handling potential gzip decompression.
// It ensures the original body or the gzip reader is closed.
// Returns the body content as a string.
func ParseBody(body io.ReadCloser, contentEncoding string) (string, error) {
	if body == nil {
		// log.Println("ParseBody called with nil body")
		return "", nil // No body content
	}

	var reader io.ReadCloser = body // Start with the original body reader

	// Check Content-Encoding header for gzip
	if strings.Contains(strings.ToLower(contentEncoding), "gzip") {
		// log.Println("Detected gzip encoding, attempting decompression")
		gzipReader, err := gzip.NewReader(body)
		if err != nil {
			log.Printf("Error creating gzip reader (falling back to raw body): %v", err)
			// Don't close original body here, let fallback read it.
			// Return error or try reading raw? Let's try reading raw.
			body.Close() // Close the original body as we failed to wrap it
            return "", fmt.Errorf("failed to create gzip reader: %w", err) // Return error is safer

		}
		// Use the gzipReader from now on. readBody will close it.
		reader = gzipReader
		// Note: The original 'body' will be closed when the gzipReader is closed.
	}

	// Read all content using the appropriate reader (original or gzip)
	bodyBytes, err := readBody(reader) // readBody handles closing the reader
	if err != nil {
		return "", err // Error occurred during reading or closing
	}

	// TODO: Detect character encoding (e.g., from Content-Type header) and decode properly if needed.
	// For now, assume UTF-8 or compatible.
	return string(bodyBytes), nil
}

// Helper for slicing string safely
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}