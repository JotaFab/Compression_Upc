package routes

import (
    "bytes"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func createTestFile(t *testing.T, name, content string) {
    t.Helper()
    err := os.WriteFile(name, []byte(content), 0644)
    if err != nil {
        t.Fatalf("Failed to create test file: %v", err)
    }
}

func TestCompressHandler(t *testing.T) {
    createTestFile(t, "test.txt", "hello hello hello world")

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

    if rr.Code != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", rr.Code)
    }
    if !strings.Contains(rr.Body.String(), "test.txt") {
        t.Errorf("Expected test.txt in response, got %s", rr.Body.String())
    }
}