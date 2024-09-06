package structs

import (
	strategy "intruder/strategy/request"
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
	HttpReq     string
	HttpRes     string
}

func NewResponse(httpRes *http.Response, elapsedTime time.Duration, httpReq string, payload string, idx int, isError bool) ResponseData {

	return ResponseData{
		RequestId:   uint16(idx + 1),
		Payload:     payload,
		StatusCode:  uint8(httpRes.StatusCode),
		TimeElapsed: uint32(elapsedTime.Milliseconds()),
		Err:         false, // TODO check if there is an actual error
		Length:      uint32(httpRes.ContentLength),
		HttpReq:     httpReq,
		HttpRes:     makeHttpResString(httpRes),
	}

}

func makeHttpResString(res *http.Response) string {
	statusLine := res.Proto + " " + res.Status + "\n"

	headersString := ""
	for headerType, value := range res.Header {
		headersString += headerType + ": " + value[0] + "\n"
	}

	body, err := strategy.ParseBody(res.Body, res.Header.Get("Content-Encoding"))
	if err != nil {
		panic(err) // Handle the error appropriately in real code
	}

	// fmt.Println("Body:", body)

	return statusLine + headersString + "\n" + body
}
