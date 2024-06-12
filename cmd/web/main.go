package main

import (
	"log"
	"net/http"
)

var openAIClient *OpenAIClient

func createConfig() *application {
	// Create OpenAI client
	openAIClient = NewOpenAIClient(ApiKey)

	return &application{OpenAIClient: openAIClient}
}

type application struct {
	Response     string
	OpenAIClient *OpenAIClient
}

func main() {
	app := createConfig()
	// Address of the HTTP server in Default is http://localhost:4000
	addr := ":4000"

	srv := &http.Server{
		Addr:    addr,
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
