package routes

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// Add this helper function to calculate compression ratio
func calculateCompressionRatio(originalSize, compressedSize int64) float64 {
	return float64(compressedSize) / float64(originalSize) * 100
}

func createTestFile(t *testing.T, name, content string) {
	t.Helper()
	err := os.WriteFile(name, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
}

func TestCompressHandler(t *testing.T) {
	// Create test file with more substantial content for better analysis
	content := strings.Repeat("hello world ", 10000)
	createTestFile(t, "test.txt", content)

	// Get original file size
	originalFileInfo, err := os.Stat("test.txt")
	if err != nil {
		t.Fatalf("Failed to get original file info: %v", err)
	}
	originalSize := originalFileInfo.Size()

	// Start timing
	start := time.Now()

	file, err := os.Open("test.txt")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to copy file: %v", err)
	}
	writer.WriteField("fileName", "test.txt")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/compress", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	compressHandler(rr, req)

	// Calculate execution time
	executionTime := time.Since(start)

	// Get compressed file size
	compressedFileInfo, err := os.Stat("process/test.txt.huff")
	if err != nil {
		t.Fatalf("Failed to get compressed file info: %v", err)
	}
	compressedSize := compressedFileInfo.Size()

	// Calculate and print metrics
	compressionRatio := calculateCompressionRatio(originalSize, compressedSize)

	t.Logf("\nCompression Analysis:")
	t.Logf("Original Size: %d bytes", originalSize)
	t.Logf("Compressed Size: %d bytes", compressedSize)
	t.Logf("Compression Ratio: %.2f%%", compressionRatio)
	t.Logf("Execution Time: %v", executionTime)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), ".huff") {
		t.Errorf("Expected .huff in response, got %s", rr.Body.String())
	}
}

func TestDecompressHandler(t *testing.T) {
	// First, compress a file so we have a .huff file to decompress
	TestCompressHandler(t)

	// Start timing
	start := time.Now()

	file, err := os.Open("process/test.txt.huff")
	if err != nil {
		t.Fatalf("Failed to open compressed file: %v", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "test.txt.huff")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Failed to copy file: %v", err)
	}
	writer.WriteField("fileName", "test.txt.huff")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/decompress", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	decompressHandler(rr, req)

	// Calculate execution time
	executionTime := time.Since(start)

	// Get decompressed file size
	decompressedFileInfo, err := os.Stat("process/test.txt")
	if err != nil {
		t.Fatalf("Failed to get decompressed file info: %v", err)
	}
	decompressedSize := decompressedFileInfo.Size()

	t.Logf("\nDecompression Analysis:")
	t.Logf("Decompressed Size: %d bytes", decompressedSize)
	t.Logf("Execution Time: %v", executionTime)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "test.txt") {
		t.Errorf("Expected test.txt in response, got %s", rr.Body.String())
	}
}
