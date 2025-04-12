package matrix

import (
	"fmt"
	"math"
	"net/http"
	"strings"
)

type ClusterBomb struct{}

const maxClusterBombRequests = 10000 // Safety limit for combinations

func (cb *ClusterBomb) BuildMatrix(req *http.Request) ([][]string, error) {
	originalMatrix, err := makeOriginalMatrix(req)
	if err != nil {
		return nil, fmt.Errorf("cluster bomb: failed to read payloads: %w", err)
	}

	numberOfPayloadPositions := strings.Count(req.FormValue("requestData"), "§") / 2
	if numberOfPayloadPositions == 0 {
		return nil, fmt.Errorf("cluster bomb: no payload markers (§...§) found in the request template")
	}

	// Cluster bomb requires exactly one payload list per payload position marker
	if len(originalMatrix) != numberOfPayloadPositions {
		return nil, fmt.Errorf("cluster bomb: requires %d payload lists, but %d were provided",
			numberOfPayloadPositions, len(originalMatrix))
	}

    // Check for empty lists - they contribute 1 to the permutation count (empty string)
    // Calculate total combinations and check against limit
    totalCombinations := 1.0
    hasNonEmptyList := false
    for _, list := range originalMatrix {
        listLen := len(list)
        if listLen == 0 {
            listLen = 1 // Treat empty list as one possibility (empty string)
        } else {
             hasNonEmptyList = true
        }
        // Check for potential overflow before multiplying
        if float64(listLen) > math.MaxFloat64 / totalCombinations {
             return nil, fmt.Errorf("cluster bomb: potential number of requests exceeds limits (overflow)")
        }
        totalCombinations *= float64(listLen)
    }

     if !hasNonEmptyList && numberOfPayloadPositions > 0{
         // All lists are empty, results in one request with all empty payloads
     } else if totalCombinations == 0 && numberOfPayloadPositions > 0 {
          // This case shouldn't happen with the len=0 check above
          return [][]string{}, nil // No combinations
     }


	if totalCombinations > maxClusterBombRequests {
		return nil, fmt.Errorf("cluster bomb: number of requests (%.0f) exceeds maximum limit (%d). Reduce payload list sizes",
			totalCombinations, maxClusterBombRequests)
	}
    if totalCombinations == 0 {
         return [][]string{}, nil
    }


	return buildClusterBombMatrix(originalMatrix), nil
}


// multiple set of payloads (one list per position)
// Creates N1 * N2 * ... * Nk requests, where Ni is the length of the i-th payload list.
// Tests all possible combinations of payloads.
func buildClusterBombMatrix(originalMatrix [][]string) [][]string {
	if len(originalMatrix) == 0 {
		return [][]string{}
	}

	// Handle empty lists correctly: treat them as a list containing one empty string ""
    processedMatrix := make([][]string, len(originalMatrix))
    for i, list := range originalMatrix {
        if len(list) == 0 {
            processedMatrix[i] = []string{""} // Replace empty list with [""]
        } else {
            processedMatrix[i] = list
        }
    }


	// Start permutation process
	result := [][]string{{}} // Start with an empty combination

	for _, payloadList := range processedMatrix {
		var temp [][]string
		for _, existingCombo := range result {
			for _, payload := range payloadList {
				// Create new combination by appending current payload
				newCombo := make([]string, len(existingCombo)+1)
				copy(newCombo, existingCombo)
				newCombo[len(existingCombo)] = payload
				temp = append(temp, newCombo)
			}
		}
		result = temp // Update result with new combinations
	}

	return result
}