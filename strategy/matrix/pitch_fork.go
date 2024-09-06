package matrix

import "net/http"

type PitchFork struct{}

func (pf *PitchFork) BuildMatrix(req *http.Request) [][]string {
	return buildPitchForkMatrix(makeOriginalMatrix(req))
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
