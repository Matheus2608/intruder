package structs

import (
	"net/http"
	"time"
)

type ResponseData struct {
	RequestId   uint16
	Payload     string // payload user suplied
	StatusCode  uint8
	TimeElapsed uint32 // time taken to receive or complete the response
	Err         bool
	Length      uint32
}

func NewResponse(httpRes *http.Response, payload string, idx int, elapsedTime time.Duration) ResponseData {
	return ResponseData{
		RequestId:   uint16(idx + 1),
		Payload:     payload,
		StatusCode:  uint8(httpRes.StatusCode),
		TimeElapsed: uint32(elapsedTime.Milliseconds()),
		Err:         false, // TODO check if there is an actual error
		Length:      uint32(httpRes.ContentLength),
	}

}
