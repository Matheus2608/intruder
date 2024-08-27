package structs

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewResponse(t *testing.T) {
	res := &http.Response{
		StatusCode:    200,
		ContentLength: 123,
	}

	expectedPayload := "test-payload"
	expectedIdx := 1
	expectedElapsedTime := 100 * time.Millisecond
	expectedResponseData := ResponseData{
		RequestId:   uint16(expectedIdx + 1),
		Payload:     expectedPayload,
		StatusCode:  uint8(res.StatusCode),
		TimeElapsed: uint32(expectedElapsedTime.Milliseconds()),
		Err:         false,
		Length:      uint32(res.ContentLength),
	}

	responseData := NewResponse(res, expectedPayload, expectedIdx, expectedElapsedTime)

	if !reflect.DeepEqual(responseData, expectedResponseData) {
		t.Errorf("Expected responseData to be %+v, got %+v", expectedResponseData, responseData)
	}
}
