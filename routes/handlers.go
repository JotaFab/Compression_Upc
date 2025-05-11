package routes

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func MuxRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Define your routes here
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	mux.HandleFunc("/compress", compressHandler)
	mux.HandleFunc("/decompress", decompressHandler)
	mux.HandleFunc("/download", downloadHandler)

	return mux
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// el archivo a descargar parseado desde la URL
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}
	// Set the headers for the file download
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	// Open the file inside the process directory
	file, err := os.Open("process/" + fileName)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// Serve the file to the client
	// Set the content length header
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Error getting file info", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	http.ServeContent(w, r, fileName, time.Now(), file)
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
		err = saveFile(file, fileName)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// redirect to /
	} else {

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func decompressHandler(w http.ResponseWriter, r *http.Request) {
	// Handle decompression logic here
	w.Write([]byte("Decompression logic goes here"))
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
