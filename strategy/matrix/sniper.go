package matrix

import (
	"fmt"
	"net/http"
	"strings"
)

type Sniper struct{}

func (s *Sniper) BuildMatrix(req *http.Request) ([][]string, error) {
	originalMatrix, err := makeOriginalMatrix(req)
	if err != nil {
		return nil, fmt.Errorf("sniper: failed to read payloads: %w", err)
	}

	// Sniper uses only the first payload list
	if len(originalMatrix) == 0 {
		// Check if markers were present. If yes, error because payload1 is missing.
		if strings.Contains(req.FormValue("requestData"), "§") {
			return nil, fmt.Errorf("sniper: payload markers found but payload list 1 is missing or empty")
		}
		// No markers, no payloads - return empty list, attack handler should notice.
		return [][]string{}, nil
	}

	payloadList1 := originalMatrix[0]
	if len(payloadList1) == 0 {
         // Payload list 1 was provided but is empty after trimming.
         // This is valid, results in requests with empty strings replacing markers.
         // Or return error? Let's allow it.
         // return nil, fmt.Errorf("sniper: payload list 1 is empty")
    }

	numberOfPayloadPositions := strings.Count(req.FormValue("requestData"), "§") / 2
	if numberOfPayloadPositions == 0 {
		// No markers found. Sniper technically doesn't *need* markers if it only targets one implicit thing,
		// but our § mechanism requires them. Return empty? Or error?
		// Let's assume markers are required for defining positions.
		return nil, fmt.Errorf("sniper: no payload markers (§...§) found in the request template")
	}

	return buildSniperMatrix(payloadList1, numberOfPayloadPositions), nil
}

// single set of payloads (uses the first list: payloadList)
// one or more payload positions
// it places the first payload in the first position only, keeping others empty
// then the second payload in the first position only...
// then the first payload in the second position only... etc.
// Essentially, it iterates through each position, applying all payloads from list 1 to it one by one.
func buildSniperMatrix(payloadList []string, numberOfPayloadPositions int) [][]string {
	var sniperMatrix [][]string

	if numberOfPayloadPositions <= 0 {
		return [][]string{} // No positions to attack
	}

	totalRequests := len(payloadList) * numberOfPayloadPositions
	sniperMatrix = make([][]string, 0, totalRequests) // Pre-allocate slice capacity

	// Iterate through each payload position marker (1 to N)
	for posIndex := 0; posIndex < numberOfPayloadPositions; posIndex++ {
		// Iterate through each payload in the list
		for _, payload := range payloadList {
			// Create a row for the matrix representing one request
			// This row will have N columns (one for each position marker)
			requestPayloads := make([]string, numberOfPayloadPositions)
			// requestPayloads is initialized with empty strings

			// Place the current payload into the current target position
			requestPayloads[posIndex] = payload

			// Add this combination to the final matrix
			sniperMatrix = append(sniperMatrix, requestPayloads)
		}
	}

	return sniperMatrix
}


// Original incorrect Sniper interpretation (like Burp's Pitchfork but with 1 list):
// func buildSniperMatrix_Old(payloadList []string, numberOfPayloadPositions int) [][]string {
//  var sniperMatrix [][]string
//  for _, payload := range payloadList {
//      subMatrix := generateSubMatrix(payload, numberOfPayloadPositions)
//      sniperMatrix = append(sniperMatrix, subMatrix...)
//  }
//  return sniperMatrix
// }
// func generateSubMatrix_Old(payload string, numberOfPayloadPositions int) [][]string {
//  subMatrix := make([][]string, numberOfPayloadPositions)
//  for i := range subMatrix {
//      subMatrix[i] = make([]string, numberOfPayloadPositions) // All empty initially
//      subMatrix[i][i] = payload // Place payload only at the i-th position
//  }
//  return subMatrix
// }