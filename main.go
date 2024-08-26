package main

import (
	"intruder/handlers"
	"log"
	"net/http"
)

func main() {

	router := http.NewServeMux()
	// Define static directory archives handler
	fs := http.FileServer(http.Dir("./static"))

	// Define the routes
	router.Handle("GET /", fs)
	router.HandleFunc("POST /attack", handlers.AttackHandler)

	// Initialize the server
	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
