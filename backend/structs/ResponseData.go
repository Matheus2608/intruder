package structs

import (
	"fmt"
	"intruder/backend/facades"
	"net/http"
	"time"
)

type ResponseData struct {
	RequestId   uint16
	Payload     string // payload user suplied
	StatusCode  uint8
	TimeElapsed uint32 // The time taken to receive or complete the response
	Err         bool
	Length      uint32
	HttpReq     string
	HttpRes     string
}

func NewResponse(httpReq *http.Request, httpRes *http.Response, payload string, idx int, elapsedTime time.Duration) ResponseData {

	// requestBody, reqErr := readReqBody(httpReq)
	// if reqErr != nil {
	// 	fmt.Println("Error reading request body:", reqErr)
	// }

	responseBody, resErr := readResBody(httpRes)
	if resErr != nil {
		fmt.Println("Error reading response body:", resErr)
	}

	httpReqString := makehttpString(httpReq.Header, "")
	httpResString := makehttpString(httpRes.Header, responseBody)

	return ResponseData{
		RequestId:   uint16(idx + 1),
		Payload:     payload,
		StatusCode:  uint8(httpRes.StatusCode),
		TimeElapsed: uint32(elapsedTime.Milliseconds()),
		Err:         false,
		Length:      uint32(httpRes.ContentLength),
		HttpReq:     httpReqString,
		HttpRes:     httpResString}

}

func makehttpString(headers map[string][]string, body string) string {
	httpResString := ""
	for headerType, value := range headers {
		httpResString += headerType + ": " + value[0] + "\n"
	}
	httpResString += "\n" + body
	return httpResString
}

// Exemplos de uso para uma resposta HTTP
func readResBody(httpRes *http.Response) (string, error) {
	return facades.ParseBody(httpRes.Body, httpRes.Header.Get("Content-Encoding"))
}

// // Exemplos de uso para uma requisição HTTP
// func readReqBody(httpReq *http.Request) (string, error) {
// 	return facades.ParseBody(httpReq.Body, httpReq.Header.Get("Content-Encoding"))
// }
