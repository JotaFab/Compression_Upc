package main

import (
	"fmt"
	"os"
	"strings"
	"Compression_Upc/rle"
	"Compression_Upc/huffman"
	"encoding/gob"
)

// Estructura para guardar la codificación Huffman en disco
type HuffmanData struct {
	Encoded  string
	CodeMap  map[rune]string
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func saveHuffmanData(path string, data HuffmanData) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(data)
}

func loadHuffmanData(path string) (HuffmanData, error) {
	file, err := os.Open(path)
	if err != nil {
		return HuffmanData{}, err
	}
	defer file.Close()

	var data HuffmanData
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

func main() {
	// Paso 1: leer archivo original
	original, err := readFile("original.txt")
	if err != nil {
		fmt.Println("Error leyendo archivo:", err)
		return
	}
	fmt.Println("Original:", original)

	// Paso 2: compresión con RLE
	rleCompressed := rle.RLECompress(original)
	fmt.Println("RLE Comprimido:", rleCompressed)

	// Paso 3: compresión con Huffman sobre el resultado RLE
	huffEncoded, codeMap := huffman.Encode(rleCompressed)
	fmt.Println("Huffman Encoded:", huffEncoded)

	// Guardar codificación Huffman para descompresión
	huffData := HuffmanData{
		Encoded: huffEncoded,
		CodeMap: codeMap,
	}
	err = saveHuffmanData("compressed.gob", huffData)
	if err != nil {
		fmt.Println("Error guardando Huffman:", err)
		return
	}

	// Paso 4: Descompresión
	loadedData, err := loadHuffmanData("compressed.gob")
	if err != nil {
		fmt.Println("Error cargando Huffman:", err)
		return
	}
	huffDecoded := huffman.Decode(loadedData.Encoded, loadedData.CodeMap)
	rleDecoded := rle.RLEDecompress(huffDecoded)
	fmt.Println("Descomprimido Final:", rleDecoded)

	if strings.TrimSpace(original) == strings.TrimSpace(rleDecoded) {
		fmt.Println("✅ La descompresión es correcta.")
	} else {
		fmt.Println("❌ Error en la descompresión.")
	}
}
