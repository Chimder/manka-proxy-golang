package main

import (
	"fmt"
	"mankaproxy/routes"
	"net/http"
	"os"
)

func main() {
	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8080"
	}
	server := http.Server{
		Addr:    ":" + PORT,
		Handler: routes.Routes(),
	}

	fmt.Println("Server listening on port:", PORT)
	server.ListenAndServe()
}
