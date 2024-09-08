package handlers

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"intruder/strategy/matrix"
	strategy "intruder/strategy/request"
	"intruder/structs"
	"net/http"
	"strings"
	"sync"
)

func getMethodAndPath(requestLine string) (string, string) {
	requestLineList := strings.SplitN(requestLine, " ", 3)
	return requestLineList[0], requestLineList[1]
}

func AttackHandler(res http.ResponseWriter, req *http.Request) {
	// Parse the form data
	if err := req.ParseForm(); err != nil {
		http.Error(res, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// request variables
	reqString := req.FormValue("requestData")
	reqStringList := strings.Split(req.FormValue("requestData"), "\r\n")
	requestLine := reqStringList[0]
	method, path := getMethodAndPath(requestLine)

	// payload variables
	matrix_strategy := matrix.GetMatrixStrategy(req.FormValue("typeOfAttack"))
	payloads := matrix_strategy.BuildMatrix(req)
	lenPayloads := len(payloads)

	// Clone Variables
	strategyClones := make([]strategy.RequestStrategy, lenPayloads)
	var clonesWG sync.WaitGroup

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	originalStrategy := strategy.ChooseStrategy(method)
	for idx, payload := range payloads {
		// payload = strings.TrimSpace(payload)
		clonesWG.Add(1)
		go func() {
			defer clonesWG.Done()
			originalStrategy.CloneWithDifferentPayload(idx, reqString, payload, &strategyClones)
		}()
	}

	clonesWG.Wait()

	url := false
	responseList := structs.NewResponses(lenPayloads) // TODO
	for idx, clone := range strategyClones {
		clonesWG.Add(1)
		go func() {
			defer clonesWG.Done()

			httpReq, err := clone.CreateRequest(path)
			if err != nil {
				panic("Error creating request:" + err.Error())
			}

			if !url {
				responseList.URL = httpReq.URL.String()
				url = true
			}

			httpRes, elapsedTime, err := strategy.SendRequest(client, httpReq)
			if err != nil {
				fmt.Println("Error sending request:", err)
			}

			cloneHttpReq, clonePayload := clone.ToString()
			response := structs.NewResponse(
				httpRes,
				elapsedTime,
				cloneHttpReq,
				clonePayload,
				idx,
				err == nil)
			responseList.AddResponse(response)
		}()
	}

	// Wait for all workes finish to send the response
	clonesWG.Wait()

	tmp, err := template.ParseFiles("templates/attack.html")
	if err != nil {
		fmt.Println("Error parsing template:", err)
	}
	err = tmp.Execute(res, responseList)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}
