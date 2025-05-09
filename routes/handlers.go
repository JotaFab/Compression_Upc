package routes

import (
	"Compression_Upc/compression"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const uploadDir = "./uploads"

type ProcessResult struct {
	Filename      string  `json:"filename"`
	OriginalSize  int64   `json:"originalSize"`
	ProcessedSize int64   `json:"processedSize"`
	Ratio         float64 `json:"ratio"`
	DownloadPath  string  `json:"downloadPath"`
}

func SetupRoutes() *http.ServeMux {
	os.MkdirAll(uploadDir, 0755)
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", serveHome)
	mux.HandleFunc("/upload", handleFileProcess)
	mux.HandleFunc("/download/", serveFile)

	return mux
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func handleFileProcess(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al procesar archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	action := r.FormValue("action")
	result, err := processFile(file, header.Filename, action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver partial HTML con HTMX
	w.Header().Set("Content-Type", "text/html")
	renderStats(w, result)
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.URL.Path)
	path := filepath.Join(uploadDir, filename)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}

// Helpers internos
func processFile(file io.Reader, filename, action string) (ProcessResult, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return ProcessResult{}, err
	}

	var result ProcessResult
	if action == "compress" {
		compressed, err := compression.CompressText(string(content))
		if err != nil {
			return ProcessResult{}, err
		}

		filename := filepath.Join(uploadDir, filename+".compressed")
		err = compression.SaveCompressedFile(filename, compressed)
		if err != nil {
			return ProcessResult{}, err
		}

		result = ProcessResult{
			Filename:      filename,
			OriginalSize:  compressed.OriginalSize,
			ProcessedSize: compressed.CompressedSize,
			Ratio:         (1 - float64(compressed.CompressedSize)/float64(compressed.OriginalSize)) * 100,
			DownloadPath:  "/download/" + filepath.Base(filename),
		}
	} else if action == "decompress" {
		tempPath := filepath.Join(uploadDir, filename)
		tempFile, err := os.Create(tempPath)
		if err != nil {
			return ProcessResult{}, err
		}
		defer tempFile.Close()
		defer os.Remove(tempPath)

		_, err = io.Copy(tempFile, file)
		if err != nil {
			return ProcessResult{}, err
		}

		compressed, err := compression.LoadCompressedFile(tempPath)
		if err != nil {
			return ProcessResult{}, err
		}

		decompressed := compression.DecompressData(compressed)
		originalFilename := strings.TrimSuffix(filename, ".compressed")
		outputPath := filepath.Join(uploadDir, originalFilename)

		err = os.WriteFile(outputPath, []byte(decompressed), 0644)
		if err != nil {
			return ProcessResult{}, err
		}

		result = ProcessResult{
			Filename:      originalFilename,
			OriginalSize:  compressed.CompressedSize,
			ProcessedSize: int64(len(decompressed)),
			Ratio:         float64(len(decompressed)) / float64(compressed.CompressedSize) * 100,
			DownloadPath:  "/download/" + originalFilename,
		}
	}

	return result, nil
}

func renderStats(w http.ResponseWriter, result ProcessResult) {
	tmpl := `
	<div class="results" hx-trigger="load delay:3s" hx-get="/download/{{.Filename}}">
		<h3>Resultados:</h3>
		<p>Tamaño original: {{.OriginalSize}} bytes</p>
		<p>Tamaño final: {{.ProcessedSize}} bytes</p>
		<p>Ratio: {{.Ratio}}%</p>
		<p class="download-info">La descarga comenzará en 3 segundos...</p>
	</div>
	`
	template.Must(template.New("stats").Parse(tmpl)).Execute(w, result)
}
