package main

import (
	"Compression_Upc/routes"
	"fmt"
	"net/http"
)

func main() {

	handler := routes.SetupRoutes()
	fmt.Println("Servidor iniciado en http://localhost:8080")

	s := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	s.ListenAndServe()
}
