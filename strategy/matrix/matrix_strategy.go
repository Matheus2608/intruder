package matrix

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MatrixStrategy interface {
	// BuildMatrix returns the payload combinations or an error if inputs are invalid.
	BuildMatrix(req *http.Request) ([][]string, error)
}

var attackTypeMap = map[string]MatrixStrategy{
	"sniper":        &Sniper{},
	"battering-ram": &BatteringRam{},
	"pitch-fork":    &PitchFork{},
	"cluster-bomb":  &ClusterBomb{},
}

// GetMatrixStrategy returns the strategy for the given attack type.
// It panics if the type is not found (consider returning an error instead).
func GetMatrixStrategy(attackType string) MatrixStrategy {
	strategy, ok := attackTypeMap[attackType]
	if !ok {
		// Or return nil, fmt.Errorf("Type of attack '%s' not implemented", attackType)
		panic(fmt.Sprintf("Type of attack '%s' not implemented", attackType))
	}
	return strategy
}

// makeOriginalMatrix extracts payload lists from the request form.
// Returns a 2D slice where each inner slice is a payload list for one position,
// or an error if parsing fails or inputs are missing.
func makeOriginalMatrix(req *http.Request) ([][]string, error) {
	var originalPayloads [][]string
	payloadCount := 0

	// Determine the number of payload inputs expected based on markers in requestData
	// This helps validate if enough payload textareas were provided/filled.
	expectedPayloadPositions := strings.Count(req.FormValue("requestData"), "§") / 2

	if expectedPayloadPositions == 0 {
		// If there are no markers, should we still process payload1?
		// Burp behavior: Depends on attack type. Sniper/Battering Ram might use payload1 even without markers.
		// For simplicity now, let's require markers for payloads to be used.
		// return nil, fmt.Errorf("no payload markers (§...§) found in the request template")
		// Or, allow payload1 to be used implicitly if needed by the strategy? Let's require markers.
	}

	basePayloadKey := "payload"
	for i := 1; ; i++ {
		payloadKey := basePayloadKey + strconv.Itoa(i)
		payloadValue := req.FormValue(payloadKey)

		// Stop if we expect N payloads and have processed N, or if the form value is missing
        // We need *at least* as many payload inputs as markers for some modes (Pitchfork, ClusterBomb)
		if i > expectedPayloadPositions && expectedPayloadPositions > 0 {
			// We have more payload inputs than markers, which is okay for Sniper/BatteringRam
            // but might be ignored by others. Let's process all provided inputs.
		}

		if payloadValue == "" {
			// If this payload is required (i <= expectedPayloadPositions), it's an error.
			// If it's an optional extra payload input, we just stop.
			if i <= expectedPayloadPositions {
				// A payload list corresponding to a marker is missing or empty.
				// Allow empty lists? Burp allows this. Let's treat "" as an empty list [].
				 originalPayloads = append(originalPayloads, []string{}) // Append empty list
                 payloadCount++
                 // Check if we have now processed all expected positions
                 if payloadCount == expectedPayloadPositions {
                     break
                 }
                 continue // Continue checking next payload number just in case (e.g., payload1="", payload2="val")

			} else {
				// No more payload form values found, and we've processed all expected ones.
				break
			}
		}


		// Split payload list by lines, trim whitespace, remove empty lines
		payloadList := []string{}
		rawList := strings.Split(strings.ReplaceAll(payloadValue, "\r\n", "\n"), "\n") // Normalize newlines
		for _, p := range rawList {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" { // Only add non-empty payloads
				payloadList = append(payloadList, trimmed)
			}
		}
        // If after trimming, the list is empty, treat it as explicitly empty
        if len(payloadList) == 0 && payloadValue != "" { // User provided whitespace lines only
             originalPayloads = append(originalPayloads, []string{})
        } else {
		    originalPayloads = append(originalPayloads, payloadList)
        }
		payloadCount++

        // Stop if we have processed all expected positions (relevant for cluster/pitchfork)
        if payloadCount == expectedPayloadPositions && expectedPayloadPositions > 0{
             // We could break here if we strictly enforce only expected number of payloads.
             // Let's allow extra payload lists for now, strategies can ignore them if needed.
        }
	}

	// Validation: Check if we have at least one payload list if markers were present
	if expectedPayloadPositions > 0 && payloadCount == 0 {
		return nil, fmt.Errorf("payload markers (§...§) found, but no corresponding payload lists provided")
	}
    // Validation: Check if required number of payloads were provided for specific modes
    // This check should ideally happen within the specific strategy's BuildMatrix method.


	return originalPayloads, nil
}