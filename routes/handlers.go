package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

func MuxRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Define your routes here
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
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
	if r.Method == http.MethodPost {
		// Handle file upload
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Process the uploaded file (e.g., save it to disk)

		// redirect to /
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {

	}
	// Handle compression logic here
	w.Write([]byte("Compression logic goes here"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func decompressHandler(w http.ResponseWriter, r *http.Request) {
	// Handle decompression logic here
	w.Write([]byte("Decompression logic goes here"))
}
