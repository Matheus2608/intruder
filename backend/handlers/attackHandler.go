package handlers

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"intruder/backend/facades"
	"intruder/backend/structs"
	"net/http"
	"strings"
	"sync"
)

func cloneHeaderOnly(r *http.Request) *http.Request {
	if r.Header == nil {
		panic("http: Request.Header is nil")
	}
	r2 := new(http.Request)
	*r2 = *r

	r2.Header = r.Header.Clone()
	return r2
}

func initialInstanciations(req *http.Request) (*http.Request, *http.Client, []string, []string) {
	// create the request once, because only some header need to be changed later
	// this will save time and space
	requestMap, headersWhichNeedToBeChanged, path := facades.ParseRequest(req.FormValue("request-data-textarea"))
	payload := strings.Split(req.FormValue("payload-input"), "\n")
	httpReq, err := facades.CreateRequest(path, requestMap)
	if err != nil {
		panic("Error creating request:" + err.Error())
	}

	// same logic for the client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return httpReq, client, payload, headersWhichNeedToBeChanged
}

func AttackHandler(res http.ResponseWriter, req *http.Request) {
	// Parse the form data
	if err := req.ParseForm(); err != nil {
		http.Error(res, "Unable to parse form", http.StatusBadRequest)
		return
	}

	httpReq, client, payload, headersWhichNeedToBeChanged := initialInstanciations(req)

	responseList := structs.NewResponses(len(payload), httpReq.URL.String())
	// HTTPS workers
	var httpsWG sync.WaitGroup

	for idx, input := range payload {
		input = strings.TrimSpace(input)
		httpsWG.Add(1)

		go func() {

			defer httpsWG.Done()
			//fmt.Println("Sending request", idx, "with payload:", payload)

			// Create new http for all iteration because the header will be changed
			newHttpReq := cloneHeaderOnly(httpReq)
			httpRes, elapsedTime, err := facades.SendRequestWithPayload(client, newHttpReq, input, headersWhichNeedToBeChanged)
			if err != nil {
				fmt.Println("Error sending request:", idx, "with payload", payload, ":", err)
				return
			}

			response := structs.NewResponse(
				newHttpReq,
				httpRes,
				input,
				idx,
				elapsedTime)

			responseList.AddResponse(response, idx)
		}()

	}

	// Wait for all workes finish to send the response
	httpsWG.Wait()

	tmp, err := template.ParseFiles("frontend/templates/attack.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
	}
	err = tmp.Execute(res, responseList)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}
