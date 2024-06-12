package main

import (
	"html/template"
	"log"
	"net/http"
)

//TODO: Render function

type PageLayoutData struct {
	Title string
}

type TemplateData struct {
	Response       string
	PageLayoutData *PageLayoutData
	Username string
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
