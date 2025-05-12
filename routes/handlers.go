package routes

import (
	"Compression_Upc/huffman"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := "process/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Copy file to response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error sending file", http.StatusInternalServerError)
		return
	}

	// Eliminar el archivo después de enviarlo
	defer os.Remove(filePath)
}

func compressHandler(w http.ResponseWriter, r *http.Request) {
	var fileName string
	if r.Method == http.MethodPost {
		// Handle file upload
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save the uploaded file to the process directory
		fileName = r.FormValue("fileName")
		if fileName == "" {
			http.Error(w, "File name is required", http.StatusBadRequest)
			return
		}

		// quitar espacios en blanco
		fileName = strings.TrimSpace(fileName)

		err = saveFile(file, fileName)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		huffman.Compress("process/"+fileName, "process/"+fileName+".huff")

		// responde con el nombre del archivo comprimido
		w.Write([]byte(fileName + ".huff"))
		return

	} else {

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
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
