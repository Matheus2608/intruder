package matrix

import "net/http"

type ClusterBomb struct{}

func (cb *ClusterBomb) BuildMatrix(req *http.Request) [][]string {
	return buildClusterBombMatrix(makeOriginalMatrix(req))
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
