package facades

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestCreateRequestWithoutHostHeader(t *testing.T) {
	path, m := "doesntmatter", make(map[string]string)
	_, err := CreateRequest(path, m)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCreateRequestBestCase(t *testing.T) {
	path, m := "doesntmatter", map[string]string{"Host": "localhost"}
	expectedURL := "https://" + m["Host"] + path
	output, err := CreateRequest(path, m)
	if err != nil {
		t.Error("Expected nil, got", err)
	}

	if output.URL.String() != expectedURL {
		t.Error("Expected", expectedURL, "got", output.URL.String())
	}

	if output.Method != "GET" {
		t.Error("Expected GET, got", output.Method)
	}
}

func TestCreateRequestWithInvalidUrl(t *testing.T) {
	path, m := "this is a test", map[string]string{"Host": "localhost"}
	_, err := CreateRequest(path, m)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCloneHeaderOnly(t *testing.T) {
	req := &http.Request{
		Header: http.Header{
			"Content-Type": []string{"application/json"},
			"User-Agent":   []string{"Go-http-client/1.1"},
		},
	}

	clonedReq := cloneHeaderOnly(req)

	if clonedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', got '%s'", clonedReq.Header.Get("Content-Type"))
	}

	if clonedReq.Header.Get("User-Agent") != "Go-http-client/1.1" {
		t.Errorf("Expected User-Agent to be 'Go-http-client/1.1', got '%s'", clonedReq.Header.Get("User-Agent"))
	}

	req.Header.Set("Content-Type", "text/plain")
	if clonedReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected cloned Content-Type to remain 'application/json', got '%s'", clonedReq.Header.Get("Content-Type"))
	}
}

func TestChangeHeaderWithPayload(t *testing.T) {
	header, payload := "session_id=§dynamic_input§;", "123456"
	expected := "session_id=123456;"

	output := changeHearderWithPayload(header, payload)
	if output != expected {
		t.Error("Expected", expected, "got", output)
	}
}

func TestChangeHeader(t *testing.T) {
	req := &http.Request{
		Header: http.Header{
			"Content-Type": []string{"application/json"},
			"User-Agent":   []string{"Go-http-client/1.1"},
			"Cookies":      []string{"session_id=§dynamic_input§;"},
		},
	}

	headersToBeChanged := []string{"Cookies"}
	payload := "123456"

	newReq := ChangeHeader(req, headersToBeChanged, payload)

	if newReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', got '%s'", newReq.Header.Get("Content-Type"))
	}

	if newReq.Header.Get("User-Agent") != "Go-http-client/1.1" {
		t.Errorf("Expected User-Agent to be 'Go-http-client/1.1', got '%s'", newReq.Header.Get("User-Agent"))
	}

	if newReq.Header.Get("Cookies") != "session_id=123456;" {
		t.Errorf("Expected Cookies to be 'session_id=123456;', got '%s'", newReq.Header.Get("Cookies"))
	}
}

// Parses

func TestParseRequest(t *testing.T) {
	input := "GET / HTTP/1.1\r\nHost: localhost\r\nContent-Type: application/json\r\nUser-Agent: §Go-http-client/1.1§\r\n" +
		"Cookie: session_id=§id§;\r\n\r\n"

	expectedPath := "/"
	expectedHeadersToBeChanged := []string{"User-Agent", "Cookie"}
	expectedMapRequest := map[string]string{
		"Host":         "localhost",
		"Content-Type": "application/json",
		"User-Agent":   "§Go-http-client/1.1§",
		"Cookie":       "session_id=§id§;",
	}

	requestMap, headersToBeChanged, path := ParseRequest(input)

	if path != expectedPath {
		t.Errorf("Expected path to be '%s', got '%s'", expectedPath, path)
	}

	if !reflect.DeepEqual(expectedHeadersToBeChanged, headersToBeChanged) {
		t.Errorf("Expected headers to be changed %v, got %v", expectedHeadersToBeChanged, headersToBeChanged)
	}

	if !reflect.DeepEqual(expectedMapRequest, requestMap) {
		t.Errorf("Expected request map %v, got %v", expectedMapRequest, requestMap)
	}
}

// Helper function to create a gzip compressed reader
func createGzipReader(content string) io.ReadCloser {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write([]byte(content))
	_ = gz.Close()
	return io.NopCloser(&buf)
}

func TestParseBody(t *testing.T) {
	tests := []struct {
		name            string
		body            io.ReadCloser
		contentEncoding string
		expected        string
		expectError     bool
	}{
		{
			name:            "Nil body",
			body:            nil,
			contentEncoding: "",
			expected:        "",
			expectError:     true,
		},
		{
			name:            "Plain text body",
			body:            io.NopCloser(strings.NewReader("plain text")),
			contentEncoding: "",
			expected:        "plain text",
			expectError:     false,
		},
		{
			name:            "Gzip compressed body",
			body:            createGzipReader("gzip text"),
			contentEncoding: "gzip",
			expected:        "gzip text",
			expectError:     false,
		},
		{
			name:            "Invalid gzip body",
			body:            io.NopCloser(strings.NewReader("invalid gzip")),
			contentEncoding: "gzip",
			expected:        "",
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBody(tt.body, tt.contentEncoding)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
			if result != tt.expected {
				t.Errorf("Expected result: %s, got: %s", tt.expected, result)
			}
		})
	}
}
