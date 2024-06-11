package main

import (
	"fmt"
	"net/http"
)

var Response string

// Renders homepage
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, Response)
}

// Fetches the Degug response from OpenAI and redirects to home
func (app *application) fetchDebugPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, Response)
}
func (app *application) fetchDebug(w http.ResponseWriter, r *http.Request) {
	prompt := r.FormValue("Message")

	DebuggerSystem := &ChatSystem{
		SystemMessage: debuggerSystemMessage,
	}

	// Call completion API
	response, err := client.CompleteText(prompt, DebuggerSystem)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	Response = string(response)
	render(w, files, Response)

	fmt.Println("Response:", Response)
}