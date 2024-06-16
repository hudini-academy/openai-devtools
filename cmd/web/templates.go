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

func render(w http.ResponseWriter, files []string, data *TemplateData) {
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}
