package structs

import (
	"reflect"
	"testing"
)

func TestNewResponses(t *testing.T) {
	size := 5
	url := "http://example.com"

	responses := NewResponses(size, url)

	if len(responses.List) != size {
		t.Errorf("Expected responses list to have length %d, got %d", size, len(responses.List))
	}

	if responses.URL != url {
		t.Errorf("Expected responses URL to be %s, got %s", url, responses.URL)
	}
}

func TestAddResponse(t *testing.T) {
	size := 5
	url := "http://example.com"
	responses := NewResponses(size, url)

	response := ResponseData{
		RequestId: 3,
	}

	responses.AddResponse(response)

	if !reflect.DeepEqual(responses.List[2], response) {
		t.Errorf("Expected responses list to have response %+v at index %d, got %+v", response, response.RequestId-1, responses.List[response.RequestId-1])
	}
}
