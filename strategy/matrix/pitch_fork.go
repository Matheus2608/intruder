package matrix

import (
	"fmt"
	"net/http"
	"strings"
)

type PitchFork struct{}

func (pf *PitchFork) BuildMatrix(req *http.Request) ([][]string, error) {
	originalMatrix, err := makeOriginalMatrix(req)
	if err != nil {
		return nil, fmt.Errorf("pitch fork: failed to read payloads: %w", err)
	}

	numberOfPayloadPositions := strings.Count(req.FormValue("requestData"), "§") / 2
	if numberOfPayloadPositions == 0 {
		return nil, fmt.Errorf("pitch fork: no payload markers (§...§) found in the request template")
	}

	// Pitchfork requires exactly one payload list per payload position marker
	if len(originalMatrix) != numberOfPayloadPositions {
		return nil, fmt.Errorf("pitch fork: requires %d payload lists, but %d were provided",
			numberOfPayloadPositions, len(originalMatrix))
	}

	// Check if all lists have the same number of payloads
	if len(originalMatrix) > 0 {
		expectedLength := len(originalMatrix[0])
		for i := 1; i < len(originalMatrix); i++ {
			if len(originalMatrix[i]) != expectedLength {
				return nil, fmt.Errorf("pitch fork: payload lists must have the same number of items (list 1 has %d, list %d has %d)",
					expectedLength, i+1, len(originalMatrix[i]))
			}
		}
		// If expectedLength is 0, all lists are empty, which is valid.
        if expectedLength == 0 {
             return [][]string{}, nil // No requests will be made if lists are empty
        }
	} else {
        // No positions, no matrix (already handled by numberOfPayloadPositions check)
        return [][]string{}, nil
    }


	return buildPitchForkMatrix(originalMatrix), nil
}

// multiple set of payloads (one list per position)
// requires N lists for N positions
// requires all lists have the same length M
// Creates M requests.
// Request 1 uses payload 1 from list 1, payload 1 from list 2, ... payload 1 from list N
// Request 2 uses payload 2 from list 1, payload 2 from list 2, ... payload 2 from list N
// ...
// Request M uses payload M from list 1, payload M from list 2, ... payload M from list N
func buildPitchForkMatrix(originalMatrix [][]string) [][]string {
	if len(originalMatrix) == 0 || len(originalMatrix[0]) == 0 {
		return [][]string{} // No positions or no payloads in the lists
	}

	numberOfRows := len(originalMatrix[0]) // Number of requests = number of payloads in each list
	numberOfColumns := len(originalMatrix) // Number of payload positions

	pitchForkMatrix := make([][]string, numberOfRows)
	for i := range pitchForkMatrix {
		pitchForkMatrix[i] = make([]string, numberOfColumns)
	}

	// Transpose the original matrix logic
	for listIndex, payloadList := range originalMatrix { // Iterate through lists (columns in final matrix)
		for payloadIndex, payloadValue := range payloadList { // Iterate through payloads (rows in final matrix)
			if payloadIndex < numberOfRows { // Safety check
				pitchForkMatrix[payloadIndex][listIndex] = payloadValue
			}
		}
	}

	return pitchForkMatrix
}