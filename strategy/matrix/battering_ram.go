package matrix

import (
	"fmt"
	"net/http"
	"strings"
)

type BatteringRam struct{}

func (b *BatteringRam) BuildMatrix(req *http.Request) ([][]string, error) {
	originalMatrix, err := makeOriginalMatrix(req)
	if err != nil {
		return nil, fmt.Errorf("battering ram: failed to read payloads: %w", err)
	}

	// Battering Ram uses only the first payload list
	if len(originalMatrix) == 0 {
		if strings.Contains(req.FormValue("requestData"), "§") {
			return nil, fmt.Errorf("battering ram: payload markers found but payload list 1 is missing or empty")
		}
		return [][]string{}, nil // No markers, no payloads
	}

    payloadList1 := originalMatrix[0]
    // Allow empty list 1

	numberOfPayloadPositions := strings.Count(req.FormValue("requestData"), "§") / 2
	if numberOfPayloadPositions == 0 {
		return nil, fmt.Errorf("battering ram: no payload markers (§...§) found in the request template")
	}

	return buildBatteringRamMatrix(payloadList1, numberOfPayloadPositions), nil
}

// single set of payloads (uses payloadList1)
// iterates through the payloads
// and places the *same* payload into *all* defined positions at once per request
func buildBatteringRamMatrix(payloadList []string, numberOfPayloadPositions int) [][]string {
	if numberOfPayloadPositions <= 0 {
		return [][]string{}
	}

	numberOfRows := len(payloadList) // Each row in the final matrix corresponds to one payload
	batteringRamMatrix := make([][]string, numberOfRows)

	for rowIdx, payload := range payloadList {
		// Create the payload array for this request
		requestPayloads := make([]string, numberOfPayloadPositions)
		for col := 0; col < numberOfPayloadPositions; col++ {
			requestPayloads[col] = payload // Use the same payload for all positions
		}
		batteringRamMatrix[rowIdx] = requestPayloads
	}

	return batteringRamMatrix
}