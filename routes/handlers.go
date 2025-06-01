package routes

import (
	"Compression_Upc/huffman"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func MuxRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files from the static directory
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/compress", compressHandler)
	mux.HandleFunc("/decompress", decompressHandler)
	mux.HandleFunc("/download", downloadHandler)

	return mux
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Sanitize filename to prevent path traversal
	fileName = filepath.Base(fileName)

	// Ensure the file is in the process directory
	filePath := filepath.Join("process", fileName)
	if !strings.HasPrefix(filePath, filepath.Join("process")) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Set headers with sanitized filename
	w.Header().Set("Content-Disposition", "attachment; filename="+sanitizeFilename(fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Copy file to response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error sending file", http.StatusInternalServerError)
		return
	}

	// Delete file after sending
	if err := os.Remove(filePath); err != nil {
		// Log the error but don't return it to the client
		log.Printf("Error removing file %s: %v", filePath, err)
	}
}

// Add this helper function
func sanitizeFilename(filename string) string {
	// Remove any non-alphanumeric characters except for common file extensions
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			return r
		}
		return -1
	}, filename)
}

func compressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Handle file upload
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get and validate filename
	fileName := strings.TrimSpace(r.FormValue("fileName"))
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	// Sanitize filename to prevent path traversal
	fileName = filepath.Base(fileName)

	// Create process directory if it doesn't exist
	if err := os.MkdirAll("process", 0755); err != nil {
		http.Error(w, "Error creating process directory", http.StatusInternalServerError)
		return
	}

	// Save file
	inputPath := filepath.Join("process", fileName)
	if err := saveFile(file, fileName); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(inputPath)

	// Compress file
	outputPath := inputPath + ".huff"
	if err := huffman.Compress(inputPath, outputPath); err != nil {
		http.Error(w, "Error compressing file", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fileName + ".huff"))
}

func decompressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle file upload
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Get filename from form
		fileName := r.FormValue("fileName")
		if fileName == "" {
			http.Error(w, "File name is required", http.StatusBadRequest)
			return
		}

		// Remove spaces and verify .huff extension
		fileName = strings.TrimSpace(fileName)
		if !strings.HasSuffix(fileName, ".huff") {
			http.Error(w, "File must have .huff extension", http.StatusBadRequest)
			return
		}

		// Save uploaded file
		err = saveFile(file, fileName)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Remove .huff extension for output file
		outputName := strings.TrimSuffix(fileName, ".huff")

		// Decompress file
		err = huffman.Decompress("process/"+fileName, "process/"+outputName)
		if err != nil {
			http.Error(w, "Error decompressing file", http.StatusInternalServerError)
			return
		}

		// Return the output filename
		w.Write([]byte(outputName))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func saveFile(file multipart.File, fileName string) error {
	outFile, err := os.Create("process/" + fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy the uploaded file to the new file
	if _, err := io.Copy(outFile, file); err != nil {
		return err
	}
	return nil
}
