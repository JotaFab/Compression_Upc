package main

import (
	"Compression_Upc/routes"
	"fmt"
	"net/http"
)

func main() {

	s := http.Server{
		Addr:    ":8080",
		Handler: routes.MuxRoutes(),
	}

	fmt.Println("Servidor iniciando en http://localhost:8080")
	s.ListenAndServe()
}
