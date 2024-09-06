package matrix

import (
	"net/http"
	"strings"
)

type Sniper struct{}

func (s *Sniper) BuildMatrix(req *http.Request) [][]string {
	originalMatrix := makeOriginalMatrix(req)
	var numberOfPayloadPositions int = strings.Count(req.FormValue("requestData"), "§") / 2
	return buildSniperMatrix(originalMatrix, numberOfPayloadPositions)
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
