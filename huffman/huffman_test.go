package huffman

import (
	"bytes"
	"os"
	"testing"
)

func TestHuffman(t *testing.T) {

	testFile := "test.txt"
	// 1. Comprimir el archivo usando la función Compress
	compressedFile := "test.huff"
	err := Compress(testFile, compressedFile)
	if err != nil {
		t.Fatalf("Compression failed: %v", err)
	}

	// 3. Verificar que el archivo comprimido se haya creado
	_, err = os.Stat(compressedFile)
	if os.IsNotExist(err) {
		t.Fatalf("Compressed file was not created: %v", err)
	}

	originalData, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Error reading original file: %v", err)
	}
	// 4. Comparar el tamaño del archivo comprimido con el original
	originalSize := len(originalData)
	compressedFileInfo, err := os.Stat(compressedFile)
	if err != nil {
		t.Fatalf("Error getting compressed file size: %v", err)
	}
	compressedSize := compressedFileInfo.Size()

	t.Logf("Original size: %d bytes, Compressed size: %d bytes", originalSize, compressedSize)
	if compressedSize >= int64(originalSize) {
		t.Logf("Compression did not reduce file size.  This is OK for small files or files with uniform distribution.")
	}

	// 5. Descomprimir el archivo comprimido
	decompressedFile := "test_decompressed.txt"
	err = Decompress(compressedFile, decompressedFile)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}

	// 6. Verificar que el archivo descomprimido se haya creado
	_, err = os.Stat(decompressedFile)
	if os.IsNotExist(err) {
		t.Fatalf("Decompressed file was not created: %v", err)
	}

	// 7. Comparar el contenido del archivo descomprimido con el original
	decompressedData, err := os.ReadFile(decompressedFile)
	if err != nil {
		t.Fatalf("Error reading decompressed file: %v", err)
	}

	if !bytes.Equal(originalData, decompressedData) {
		t.Errorf("Decompressed data does not match original data.\nOriginal:\n%s\nDecompressed:\n%s", originalData, decompressedData)
	} else {
		t.Logf("Integrity check passed: Decompressed data matches original data.")
	}
}
