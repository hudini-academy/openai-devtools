package main

import (
	"OpenAIDevTools/pkg/models"
	"html/template"
	"log"
	"net/http"
)

//TODO: Render function

type TemplateData struct {
	Response       template.HTML
	PageLayoutData *models.CustomGPT
	AllButton      []*models.CustomGPT
	PromptMessage  string
	Username       string
	Flash          string
	PromptID       int
	FlashCategory  string
}

// render function parses and executes templates using provided data, 
// handles errors, logs, and returns HTTP 500 if needed. 
// Params: w http.ResponseWriter, files []string, data *TemplateData.
func render(w http.ResponseWriter, files []string, data *TemplateData) {
	// Parse the templates from the provided file paths.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// Log the error and return an HTTP 500 Internal Server Error.
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Execute the parsed templates with the provided data.
	err = ts.Execute(w, data)
	if err != nil {
		// Log the error and return an HTTP 500 Internal Server Error.
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
