package matrix

import (
	"net/http"
	"strconv"
	"strings"
)

type MatrixStrategy interface {
	BuildMatrix(req *http.Request) [][]string
}

var attackTypeMap = map[string]MatrixStrategy{
	"sniper":        &Sniper{},
	"battering-ram": &BatteringRam{},
	"pitch-fork":    &PitchFork{},
	"cluster-bomb":  &ClusterBomb{},
}

func GetMatrixStrategy(attackType string) MatrixStrategy {
	strategy, ok := attackTypeMap[attackType]
	if !ok {
		panic("Type of attack not implemented")
	}
	return strategy
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
