package main

import (
	"OpenAIDevTools/pkg/models"
	"html/template"
	"log"
	"net/http"
)

//TODO: Render function
// TemplateData is a placeholder for any data structure you may need for your templates
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

// render parses and executes a template with the provided data, handling errors and returning HTTP 500 if necessary.
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
