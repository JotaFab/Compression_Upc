package huffman

import (
	"container/heap"
	"strings"
)

// Nodo del árbol de Huffman
type HuffmanNode struct {
	Char      rune
	Freq      int
	Left      *HuffmanNode
	Right     *HuffmanNode
}

// Cola de prioridad
type PriorityQueue []*HuffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Freq < pq[j].Freq }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*HuffmanNode))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

// Construye el árbol de Huffman desde un mapa de frecuencias
func BuildTree(freqMap map[rune]int) *HuffmanNode {
	pq := make(PriorityQueue, 0)
	for char, freq := range freqMap {
		pq = append(pq, &HuffmanNode{Char: char, Freq: freq})
	}
	heap.Init(&pq)

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*HuffmanNode)
		right := heap.Pop(&pq).(*HuffmanNode)
		newNode := &HuffmanNode{
			Char: 0,
			Freq: left.Freq + right.Freq,
			Left: left,
			Right: right,
		}
		heap.Push(&pq, newNode)
	}

	return heap.Pop(&pq).(*HuffmanNode)
}

// Genera el diccionario de codificación
func BuildCodes(node *HuffmanNode, prefix string, codeMap map[rune]string) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		codeMap[node.Char] = prefix
	}
	BuildCodes(node.Left, prefix+"0", codeMap)
	BuildCodes(node.Right, prefix+"1", codeMap)
}

// Codifica el texto
func Encode(text string) (string, map[rune]string) {
	// Paso 1: calcular frecuencias
	freqMap := make(map[rune]int)
	for _, ch := range text {
		freqMap[ch]++
	}

	// Paso 2: construir árbol de Huffman
	root := BuildTree(freqMap)

	// Paso 3: generar codificación
	codeMap := make(map[rune]string)
	BuildCodes(root, "", codeMap)

	// Paso 4: codificar el texto
	var encoded strings.Builder
	for _, ch := range text {
		encoded.WriteString(codeMap[ch])
	}

	return encoded.String(), codeMap
}

// Decodifica un texto binario con el árbol original
func Decode(encoded string, codeMap map[rune]string) string {
	// Invertir el mapa
	inverseMap := make(map[string]rune)
	for k, v := range codeMap {
		inverseMap[v] = k
	}

	var result strings.Builder
	var current strings.Builder
	for _, bit := range encoded {
		current.WriteRune(bit)
		if ch, ok := inverseMap[current.String()]; ok {
			result.WriteRune(ch)
			current.Reset()
		}
	}
	return result.String()
}
