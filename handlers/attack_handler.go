package handlers

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"intruder/strategy"
	"intruder/structs"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func getMethodAndPath(requestLine string) (string, string) {
	requestLineList := strings.SplitN(requestLine, " ", 3)
	return requestLineList[0], requestLineList[1]
}

func chooseStrategy(method string) strategy.RequestStrategy {
	switch method {
	case "GET":
		return &strategy.GetRequestStrategy{}
	case "POST":
		return &strategy.PostRequestStrategy{}
	default:
		panic("Method not implemented")
	}
}

func getPayloads(req *http.Request) [][]string {
	var originalPayloads [][]string

	payload1 := strings.Split(req.FormValue("payload1"), "\r\n")
	expectedNumberOfInputs := len(payload1)
	if expectedNumberOfInputs == 0 {
		return nil
	}
	originalPayloads = append(originalPayloads, payload1)

	basePayloadKey := "payload"
	for i := 2; ; i++ {
		payloadKey := basePayloadKey + strconv.Itoa(i)
		payload := req.FormValue(payloadKey)
		if payload == "" {
			break
		}

		payloadList := strings.Split(payload, "\r\n")
		// fmt.Printf("payloadKey: %s, payloadList: %v\n", payloadKey, payloadList)

		if len(payloadList) != expectedNumberOfInputs {
			panic("Different number of payloads")
		}

		originalPayloads = append(originalPayloads, payloadList)
	}

	return convertMatrix(originalPayloads)
}

func convertMatrix(originalMatrix [][]string) [][]string {
	numberOfRows := len(originalMatrix[0])
	numberOfColumns := len(originalMatrix)
	convertedMatrix := make([][]string, numberOfRows)
	for i := range convertedMatrix {
		convertedMatrix[i] = make([]string, numberOfColumns)
	}

	for row, rowList := range originalMatrix {
		for col, value := range rowList {
			convertedMatrix[col][row] = value
		}
	}

	return convertedMatrix
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
	payloads := getPayloads(req)
	lenPayloads := len(payloads)

	// Clone Variables
	strategyClones := make([]strategy.RequestStrategy, lenPayloads)
	var clonesWG sync.WaitGroup

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	originalStrategy := chooseStrategy(method)
	for idx, payload := range payloads {
		// payload = strings.TrimSpace(payload)
		clonesWG.Add(1)
		go func() {
			defer clonesWG.Done()
			originalStrategy.CloneWithDifferentPayload(idx, reqString, payload, &strategyClones)
		}()
	}

	clonesWG.Wait()

	// TODO
	responseList := structs.NewResponses(lenPayloads, "")
	for idx, clone := range strategyClones {
		clonesWG.Add(1)
		go func() {
			defer clonesWG.Done()

			httpReq, err := clone.CreateRequest(path)
			if err != nil {
				panic("Error creating request:" + err.Error())
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
