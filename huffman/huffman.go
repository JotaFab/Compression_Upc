package huffman

import (
	"bytes"
	"container/heap"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Huffman es un paquete que implementa la compresión y descompresión de archivos usando el algoritmo de Huffman.
// El algoritmo de Huffman es un método de compresión sin pérdida que utiliza un árbol binario para codificar los datos.
// El paquete proporciona dos funciones principales: Compress y Decompress.
// Compress comprime un archivo de entrada y guarda el resultado en un archivo de salida con extensión .huff.
// Decompress descomprime un archivo de entrada .huff y guarda el resultado en un archivo de salida.


// Estructura para el nodo del árbol de Huffman
type huffmanNode struct {
	Frequency int          `json:"frequency"`
	Char      byte         `json:"char"`
	Left      *huffmanNode `json:"left,omitempty"`
	Right     *huffmanNode `json:"right,omitempty"`
}

// Definición del min-heap para los nodos de Huffman
type priorityQueue []*huffmanNode

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].Frequency < pq[j].Frequency
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x any) {
	node := x.(*huffmanNode)
	*pq = append(*pq, node)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func buildHuffmanTree(frequencies map[byte]int) *huffmanNode {
	pq := make(priorityQueue, 0, len(frequencies))
	for char, freq := range frequencies {
		pq = append(pq, &huffmanNode{Frequency: freq, Char: char})
	}
	heap.Init(&pq)

	for pq.Len() > 1 {
		// Extraer los dos nodos con menor frecuencia
		node1 := heap.Pop(&pq).(*huffmanNode)
		node2 := heap.Pop(&pq).(*huffmanNode)

		// Crear un nuevo nodo interno con la suma de las frecuencias
		mergedNode := &huffmanNode{
			Frequency: node1.Frequency + node2.Frequency,
			Left:      node1,
			Right:     node2,
		}

		// Insertar el nuevo nodo en la cola de prioridad
		heap.Push(&pq, mergedNode)
	}

	// El nodo restante en la cola de prioridad es la raíz del árbol de Huffman
	if pq.Len() == 1 {
		return pq[0]
	}
	return nil // En caso de que el archivo esté vacío
}

// Compress comprime el archivo de entrada y guarda el resultado en el archivo de salida con extensión .huff.
func Compress(inputFile string, outputFile string) error {
	// 1. Leer el contenido del archivo de entrada.
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}
	// 2. Contar la frecuencia de cada byte en un mapa o diccionario.
	frequencies := countFrequencies(content)
	// 3. Crear el árbol de Huffman.
	root := buildHuffmanTree(frequencies)
	// PrintTreeAsJSON(root)
	// PrintLeafs(root)
	// 4. Crear la tabla de códigos Huffman.
	codes := generateCodes(root)
	fmt.Println("Códigos Huffman generados:")
	for char, code := range codes {
		fmt.Printf("Carácter: '%c', Código: %s\n", char, code)
	}
	// 5. Codificar el contenido del archivo usando la tabla de códigos.
	encodedData := encode(content, codes)
	// 6. Guardar los datos comprimidos y la tabla de códigos en el archivo de salida.
	err = saveCompressedFile(outputFile, encodedData, codes)
	if err != nil {
		return err
	}
	return nil
}

// Decompress descomprime el archivo de entrada .huff y guarda el resultado en el archivo de salida.
func Decompress(inputFile string, outputFile string) error {
	// 1. Leer el archivo comprimido, incluyendo los códigos Huffman.
	encodedData, codes, err := readCompressedFile(inputFile)
	if err != nil {
		return err
	}

	// 2. Decodificar los datos comprimidos utilizando los códigos Huffman.
	decodedData := decode(encodedData, codes)

	// 3. Guardar los datos descomprimidos en el archivo de salida.
	err = os.WriteFile(outputFile, decodedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Funciones auxiliares que implementaremos a continuación:
func countFrequencies(data []byte) map[byte]int {
	frequencies := make(map[byte]int)
	for _, b := range data {
		frequencies[b]++
	}
	return frequencies
}

func generateCodes(root *huffmanNode) map[byte]string {
	codes := make(map[byte]string)
	generateCodesRecursive(root, "", codes)
	return codes
}

func generateCodesRecursive(node *huffmanNode, currentCode string, codes map[byte]string) {
	if node == nil {
		return
	}

	// Si es un nodo hoja, hemos llegado a un carácter
	if node.Left == nil && node.Right == nil {
		codes[node.Char] = currentCode
		return
	}

	// Recorrer el subárbol izquierdo añadiendo '0' al código
	generateCodesRecursive(node.Left, currentCode+"0", codes)

	// Recorrer el subárbol derecho añadiendo '1' al código
	generateCodesRecursive(node.Right, currentCode+"1", codes)
}
func encode(data []byte, codes map[byte]string) []byte {
	var encodedData bytes.Buffer
	for _, b := range data {
		if code, ok := codes[b]; ok {
			encodedData.WriteString(code)
		}
	}

	// Convertir la secuencia de bits (string) a un slice de bytes
	return bitsToBytes(encodedData.String())
}

func bitsToBytes(bits string) []byte {
	var result bytes.Buffer
	var currentByte byte
	bitCount := 0

	for _, bit := range bits {
		currentByte <<= 1
		if bit == '1' {
			currentByte |= 1
		}
		bitCount++

		if bitCount == 8 {
			result.WriteByte(currentByte)
			currentByte = 0
			bitCount = 0
		}
	}

	// Si quedan bits incompletos al final, los rellenamos con ceros
	if bitCount > 0 {
		currentByte <<= (8 - bitCount)
		result.WriteByte(currentByte)
	}

	return result.Bytes()
}

func saveCompressedFile(outputFile string, encodedData []byte, codes map[byte]string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 1. Escribir el número de códigos
	numCodes := uint32(len(codes))
	err = binary.Write(file, binary.BigEndian, &numCodes)
	if err != nil {
		return err
	}

	// 2. Escribir cada código en el encabezado
	for char, code := range codes {
		// Escribir el byte
		err = binary.Write(file, binary.BigEndian, &char)
		if err != nil {
			return err
		}

		// Escribir la longitud del código
		codeLen := uint8(len(code))
		err = binary.Write(file, binary.BigEndian, &codeLen)
		if err != nil {
			return err
		}

		// Escribir el código
		_, err = file.WriteString(code)
		if err != nil {
			return err
		}
	}

	// 3. Escribir los datos comprimidos
	_, err = file.Write(encodedData)
	if err != nil {
		return err
	}

	return nil
}

func readCompressedFile(inputFile string) ([]byte, map[byte]string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// 1. Leer el número de códigos
	var numCodes uint32
	err = binary.Read(file, binary.BigEndian, &numCodes)
	if err != nil {
		return nil, nil, err
	}

	// 2. Leer la tabla de códigos del encabezado
	codes := make(map[byte]string)
	for i := uint32(0); i < numCodes; i++ {
		// Leer el byte
		var char byte
		err = binary.Read(file, binary.BigEndian, &char)
		if err != nil {
			return nil, nil, err
		}

		// Leer la longitud del código
		var codeLen uint8
		err = binary.Read(file, binary.BigEndian, &codeLen)
		if err != nil {
			return nil, nil, err
		}

		// Leer el código
		codeBytes := make([]byte, codeLen)
		_, err = io.ReadFull(file, codeBytes)
		if err != nil {
			return nil, nil, err
		}
		codes[char] = string(codeBytes)
	}
	// 3. Leer los datos comprimidos restantes del archivo
	encodedData, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	return encodedData, codes, nil
}

func decode(encodedData []byte, codes map[byte]string) []byte {
	// Crear la tabla de decodificación inversa (código -> byte)
	decodingTable := make(map[string]byte)
	for char, code := range codes {
		decodingTable[code] = char
	}

	var decodedData bytes.Buffer
	var currentCode string

	// Convertir los datos comprimidos (bytes) a una secuencia de bits (string)
	bitString := bytesToBits(encodedData)

	for _, bit := range bitString {
		currentCode += string(bit)
		if originalByte, ok := decodingTable[currentCode]; ok {
			decodedData.WriteByte(originalByte)
			currentCode = "" // Reiniciar el código actual
		}
	}

	return decodedData.Bytes()
}
func bytesToBits(data []byte) string {
	var bitsBuffer bytes.Buffer
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1
			if bit == 1 {
				bitsBuffer.WriteString("1")
			} else {
				bitsBuffer.WriteString("0")
			}
		}
	}
	return bitsBuffer.String()
}

// Add this new function to print the tree as JSON
func PrintTreeAsJSON(node *huffmanNode) error {
	jsonData, err := json.MarshalIndent(node, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

// PrintLeafs prints the leaf nodes of the Huffman tree
func PrintLeafs(node *huffmanNode) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		fmt.Printf("Character: '%c', Frequency: %d\n", node.Char, node.Frequency)
		return
	}
	PrintLeafs(node.Left)
	PrintLeafs(node.Right)
}
