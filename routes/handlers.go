package routes

import (
	"net/http"
)

func MuxRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Define your routes here
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.Handle("/upload", uploadHandler())
	mux.HandleFunc("/compress", compressHandler)
	mux.HandleFunc("/decompress", decompressHandler)

	return mux
}

func compressHandler(w http.ResponseWriter, r *http.Request) {
	// Handle compression logic here
	w.Write([]byte("Compression logic goes here"))
}
func decompressHandler(w http.ResponseWriter, r *http.Request) {
	// Handle decompression logic here
	w.Write([]byte("Decompression logic goes here"))
}
