package main

import (
	"log"
	"net/http"
	"flag"
)

var client *OpenAIClient

func createConfig() {
	// Create OpenAI client
	client = NewOpenAIClient(ApiKey)
}


type application struct{
	Response string
}

func main() {
	createConfig()
// Address of the HTTP server in Default is http://localhost:4000
    addr := flag.String("addr", ":4000", "HTTP network address")

	app := &application{}
	srv:=&http.Server{
		Addr: *addr,
		Handler: app.routes(),
	}
	err:=srv.ListenAndServe()
	if err!=nil{
        log.Fatal(err)
    }
}
