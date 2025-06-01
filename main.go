package main

import (
	"Compression_Upc/routes"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	// Define flags properly
	port := flag.String("p", "8080", "port to serve on")
	help := flag.Bool("h", false, "display help")

	// Parse flags
	flag.Parse()

	// Show help if requested
	if *help {
		fmt.Println("Uso: go run main.go [opciones]")
		fmt.Println("Opciones:")
		fmt.Println("  -p, --port <puerto>   Especifica el puerto en el que se ejecutará el servidor (por defecto es 8080)")
		fmt.Println("  -h, --help            Muestra esta ayuda")
		return
	}

	// Format port string
	cobraPort := ":" + *port

	// Create and configure server
	s := http.Server{
		Addr:    cobraPort,
		Handler: routes.MuxRoutes(),
	}

	fmt.Printf("Servidor iniciando en http://localhost%s\n", cobraPort)
	err := s.ListenAndServe()
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
		return
	}
	fmt.Println("Servidor detenido")
}
