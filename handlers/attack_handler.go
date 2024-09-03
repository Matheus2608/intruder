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

func makeOriginalMatrix(req *http.Request) [][]string {
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

		// if len(payloadList) != expectedNumberOfInputs {
		// 	fmt.Println("payloadList:", payloadList)
		// 	fmt.Println("expectedNumberOfInputs:", expectedNumberOfInputs)
		// 	panic("Different number of payloads")
		// }

		originalPayloads = append(originalPayloads, payloadList)
	}

	return originalPayloads
}

func getPayloads(req *http.Request) [][]string {
	originalPayloads := makeOriginalMatrix(req)

	switch req.FormValue("typeOfAttack") {
	case "sniper":
		var numberOfPayloadPositions int = strings.Count(req.FormValue("requestData"), "§") / 2
		return buildSniperMatrix(originalPayloads, numberOfPayloadPositions)
	case "battering-ram":
		var numberOfPayloadPositions int = strings.Count(req.FormValue("requestData"), "§") / 2
		return buildBatteringRamMatrix(originalPayloads, numberOfPayloadPositions)
	case "pitch-fork":
		return buildPitchForkMatrix(originalPayloads)
	case "cluster-bomb":
		return buildClusterBombMatrix(originalPayloads)
	default:
		panic("Type of attack not implemented")
	}
}

// single set of payloads
// one or more payload positions
// it places the first payload in the first position
// then the second payload in the seocnd position
// and so on
func buildSniperMatrix(originalMatrix [][]string, numberOfPayloadPositions int) [][]string {
	payloadList := originalMatrix[0]
	var sniperMatrix [][]string

	for _, payload := range payloadList {
		subMatrix := generateSubMatrix(payload, numberOfPayloadPositions)
		sniperMatrix = append(sniperMatrix, subMatrix...)
	}

	return sniperMatrix
}

func generateSubMatrix(payload string, numberOfPayloadPositions int) [][]string {
	subMatrix := make([][]string, numberOfPayloadPositions)
	for i := range subMatrix {
		subMatrix[i] = make([]string, numberOfPayloadPositions)
		subMatrix[i][i] = payload
	}

	return subMatrix
}

// single set of payloads
// iterates through the payloads
// and places the same payload  into all defined positions at once
func buildBatteringRamMatrix(originalMatrix [][]string, numberOfPayloadPositions int) [][]string {
	payloadList := originalMatrix[0]
	numberOfRows := len(payloadList)
	numberOfColumns := numberOfPayloadPositions

	batteringRamMatrix := make([][]string, numberOfRows)
	for row, payload := range payloadList {
		batteringRamMatrix[row] = make([]string, numberOfColumns)
		for col := 0; col < numberOfColumns; col++ {
			batteringRamMatrix[row][col] = payload
		}
	}

	return batteringRamMatrix
}

// multiple set of payloads
// different payload position for each defined position
// iterates through all payloads sets simultanesouly
// uses the first payload from each set
// then the second payload from each set
// and so on
func buildPitchForkMatrix(originalMatrix [][]string) [][]string {
	numberOfRows := len(originalMatrix[0])
	numberOfColumns := len(originalMatrix)
	pitchForkMatrix := make([][]string, numberOfRows)
	for i := range pitchForkMatrix {
		pitchForkMatrix[i] = make([]string, numberOfColumns)
	}

	for row, rowList := range originalMatrix {
		for col, value := range rowList {
			pitchForkMatrix[col][row] = value
		}
	}

	return pitchForkMatrix
}

// multiple set of payloads
// different payload position for each defined position
// iterates through all payloads sets in turn
// so all permutations of payloads combinations are tested
func buildClusterBombMatrix(originalMatrix [][]string) [][]string {
	if len(originalMatrix) == 0 {
		return [][]string{}
	}

	// Inicia a matriz de permutações com os valores da primeira linha
	result := [][]string{}
	for _, val := range originalMatrix[0] {
		result = append(result, []string{val})
	}

	// Itera sobre as demais linhas para gerar as combinações
	for i := 1; i < len(originalMatrix); i++ {
		var temp [][]string
		for _, res := range result {
			for _, val := range originalMatrix[i] {
				// Cria uma nova linha para cada combinação e adiciona ao resultado
				newRow := append([]string{}, res...)
				newRow = append(newRow, val)
				temp = append(temp, newRow)
			}
		}
		result = temp
	}

	return result
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

	responseList := structs.NewResponses(lenPayloads, "") // TODO
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
