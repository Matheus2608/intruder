package structs

import (
	"bytes"
	strategy "intruder/strategy/request"
	// "io" //
	"log"
	"net/http"
	// "strconv"
	"time"
)

type ResponseData struct {
	RequestId   uint16 // Index + 1
	Payload     string // Joined payload string for this request
	StatusCode  int    // Use int for standard http codes
	TimeElapsed uint32 // Time taken (ms)
	Err         bool   // Flag indicating if an error occurred during request/response handling
	ErrorMsg    string // Specific error message if Err is true
	Length      int64  // Content length (-1 if unknown, or actual size)
	HttpReq     string // String representation of the sent request
	HttpRes     string // String representation of the received response (headers + body)
}

// NewResponse creates a ResponseData for a successful request/response cycle
// Note: Reads and closes the response body.
func NewResponse(httpRes *http.Response, elapsedTime time.Duration, httpReqStr string, payloadStr string, idx int, success bool) ResponseData {
	if httpRes == nil {
		// This case should ideally be handled by NewErrorResponse, but as a fallback:
		log.Printf("Warning: NewResponse called with nil http.Response for index %d", idx)
		return NewErrorResponse(idx, "Internal error: Received nil response object")
	}

	// Ensure body is read and closed
	resString, bodyLen, readErr := makeHttpResStringAndLength(httpRes)
	if readErr != nil {
		log.Printf("ERROR reading response body for index %d: %v", idx, readErr)
		// Return an error response instead
		return NewErrorResponse(idx, "Error reading response body: "+readErr.Error())
	}

	// Use actual body length if ContentLength header is not reliable (-1)
	length := httpRes.ContentLength
	if length < 0 {
		length = bodyLen // Use the length of the body we actually read
	}


	return ResponseData{
		RequestId:   uint16(idx + 1),
		Payload:     payloadStr,
		StatusCode:  httpRes.StatusCode,
		TimeElapsed: uint32(elapsedTime.Milliseconds()),
		Err:         !success, // Mark as error if the sending process failed (passed via 'success')
		ErrorMsg:    "",       // No specific error message for successful sends
		Length:      length,
		HttpReq:     httpReqStr,
		HttpRes:     resString,
	}
}

// NewErrorResponse creates a ResponseData when an error occurred before or during the request
func NewErrorResponse(idx int, errorMsg string) ResponseData {
	return ResponseData{
		RequestId:   uint16(idx + 1),
		Payload:     "N/A", // Payload might not be available or relevant
		StatusCode:  0,     // No status code
		TimeElapsed: 0,     // No timing available
		Err:         true,  // Mark as error
		ErrorMsg:    errorMsg, // Store the specific error
		Length:      0,
		HttpReq:     "Error occurred", // Placeholder request string
		HttpRes:     "Error: " + errorMsg, // Placeholder response string
	}
}

// makeHttpResStringAndLength reads the response body, closes it, and returns the full response string + body length.
func makeHttpResStringAndLength(res *http.Response) (string, int64, error) {
	if res == nil {
		return "Error: Nil Response", 0, nil
	}

	// Start building response string
	var respBuf bytes.Buffer

	// Status Line
	respBuf.WriteString(res.Proto)
	respBuf.WriteString(" ")
	respBuf.WriteString(res.Status) // Includes status code and text
	respBuf.WriteString("\n")

	// Headers
	for key, values := range res.Header {
		for _, value := range values {
			respBuf.WriteString(key)
			respBuf.WriteString(": ")
			respBuf.WriteString(value)
			respBuf.WriteString("\n")
		}
	}
	respBuf.WriteString("\n") // End of headers

	// Body Handling - IMPORTANT: Read and close the body here
	var bodyBytes []byte
	var readErr error
	var bodyLen int64 = 0

	if res.Body != nil {
		// Use the strategy helper for potential decompression
		bodyString, err := strategy.ParseBody(res.Body, res.Header.Get("Content-Encoding"))
		// ParseBody *should* close the body reader (original or gzip)
		if err != nil {
			readErr = err // Store error but continue to build response string so far
			log.Printf("Error parsing response body: %v", err)
			respBuf.WriteString("[Error reading body: ")
			respBuf.WriteString(err.Error())
			respBuf.WriteString("]")
		} else {
			bodyBytes = []byte(bodyString)
			bodyLen = int64(len(bodyBytes))
			respBuf.Write(bodyBytes) // Append the body content
		}
		// Ensure res.Body is closed (ParseBody should handle this via its internal readBody)
		// If ParseBody doesn't guarantee closure, add: defer res.Body.Close() earlier
	} else {
		respBuf.WriteString("[No Body]")
	}


	return respBuf.String(), bodyLen, readErr // Return full string, read body length, and any read error
}