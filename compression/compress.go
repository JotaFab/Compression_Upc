package compression

import (
	"Compression_Upc/huffman"
	"encoding/binary"
	"encoding/gob"
	"os"
)

type HuffmanData struct {
	Encoded        string
	CodeMap        map[rune]string
	OriginalSize   int64
	CompressedSize int64
}

func CompressText(text string) (HuffmanData, error) {
	huffEncoded, codeMap := huffman.Encode(text)
	// Convertir la cadena codificada a una secuencia de bits empaquetados
	bits := packBits(huffEncoded)

	return HuffmanData{
		Encoded:        huffEncoded,
		CodeMap:        codeMap,
		OriginalSize:   int64(len([]byte(text))),
		CompressedSize: int64(len(bits)),
	}, nil
}

// Convierte la cadena de '0's y '1's a bytes reales
func packBits(binStr string) []byte {
	packed := make([]byte, (len(binStr)+7)/8)
	for i := 0; i < len(binStr); i++ {
		if binStr[i] == '1' {
			packed[i/8] |= 1 << uint(7-i%8)
		}
	}
	return packed
}

func DecompressData(data HuffmanData) string {
	return huffman.Decode(data.Encoded, data.CodeMap)
}

func SaveCompressedFile(path string, data HuffmanData) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir el tamaño del mapa
	mapSize := uint32(len(data.CodeMap))
	binary.Write(file, binary.LittleEndian, mapSize)

	// Escribir el mapa de códigos
	for r, code := range data.CodeMap {
		binary.Write(file, binary.LittleEndian, uint32(r))
		codeLen := uint32(len(code))
		binary.Write(file, binary.LittleEndian, codeLen)
		file.Write([]byte(code))
	}

	// Escribir los datos comprimidos
	packed := packBits(data.Encoded)
	encodedLen := uint32(len(packed))
	binary.Write(file, binary.LittleEndian, encodedLen)
	file.Write(packed)

	return nil
}

func LoadCompressedFile(path string) (HuffmanData, error) {
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
