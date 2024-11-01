package main

import (
	"fmt"
	"mankaproxy/routes"
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: routes.Routes(),
	}

	fmt.Println("Server listening on port :8080")
	server.ListenAndServe()
}