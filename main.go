package main

import (
	"intruder/handlers"
	"log"
	"net/http"
)

func main() {

	// Criando o roteador com o novo servidor
	router := http.NewServeMux()

	// Configurando as rotas RESTful
	router.HandleFunc("POST /attack", handlers.AttackHandler)
	router.HandleFunc("GET /", handlers.GetRootHandler)
	// Iniciando o servidor
	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
