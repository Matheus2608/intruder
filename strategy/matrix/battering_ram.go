package matrix

import (
	"net/http"
	"strings"
)

type BatteringRam struct{}

func (b *BatteringRam) BuildMatrix(req *http.Request) [][]string {
	originalMatrix := makeOriginalMatrix(req)
	var numberOfPayloadPositions int = strings.Count(req.FormValue("requestData"), "§") / 2
	return buildBatteringRamMatrix(originalMatrix, numberOfPayloadPositions)
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
