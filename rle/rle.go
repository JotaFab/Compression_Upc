package rle

import (
	"strconv"
	"strings"
)

// RLECompress comprime una cadena utilizando Run-Length Encoding.
func RLECompress(input string) string {
	if len(input) == 0 {
		return ""
	}

	var result strings.Builder
	count := 1
	for i := 1; i < len(input); i++ {
		if input[i] == input[i-1] {
			count++
		} else {
			result.WriteByte(input[i-1])
			result.WriteString(strconv.Itoa(count))
			count = 1
		}
	}
	result.WriteByte(input[len(input)-1])
	result.WriteString(strconv.Itoa(count))

	return result.String()
}

// RLEDecompress descomprime una cadena codificada con RLE.
func RLEDecompress(input string) string {
	var result strings.Builder
	i := 0
	for i < len(input) {
		char := input[i]
		i++
		countStr := strings.Builder{}
		for i < len(input) && input[i] >= '0' && input[i] <= '9' {
			countStr.WriteByte(input[i])
			i++
		}
		count, err := strconv.Atoi(countStr.String())
		if err != nil {
			continue // manejar errores o ignorar valores inválidos
		}
		result.WriteString(strings.Repeat(string(char), count))
	}
	return result.String()
}
