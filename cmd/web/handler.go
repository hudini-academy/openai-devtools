package main

import (
	"OpenAIDevTools/pkg/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

var Response string
var AllCustomGPT []*models.CustomGPT

// Renders homepage
func (app *application) Base(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from session

	// Check if user id present in session
	// if not, show landing page
	isAuth, ok := app.session.Get(r, "Authenticated").(bool)
	if !ok {
		app.errorLog.Println("Error getting session")
		// http.Error(w, "Internal Server Error", 500)
	}

	// If user is not authenticated, show landing page
	if !isAuth {
		app.landingPage(w, r)
		return
	}

	userID := app.session.GetInt(r, "username")

	if userID == 0 {
		app.session.Put(r, "flash", "Unauthorized access")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch user details from userID
	userName, err := app.users.GetUsernameByID(userID)

	if err != nil {
		app.errorLog.Println("Error fetching user name:", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	AllCustomGPT, err := app.CustomGPT.GetFunction()
	if err != nil {
		return
	}

	// Render home page with user's name
	data := &TemplateData{
		Username:  userName,
		AllButton: AllCustomGPT,
	}

	// if preset, show the dashboard
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, data)
}

// Renders Landing page
func (app *application) landingPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/landing.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/hero.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, nil)
}

// Fetches the Degug response from OpenAI and redirects to home

func (app *application) renderCustomGPT(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/new.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, nil)
}

// Create a button by getting the name and message
func (app *application) createCustomGPT(w http.ResponseWriter, r *http.Request) {
	ButtonName := r.FormValue("FunctionName")
	newMessage := r.FormValue("FunctionMessage")
	err := app.CustomGPT.InsertFunction(ButtonName, newMessage)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	http.Redirect(w, r, "/rendercustomgpt", http.StatusSeeOther)
}

func (app *application) deleteCustomGPT(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	fmt.Println(id)
	err := app.CustomGPT.DeleteFunction(id)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// Render the page Prompt page
func (app *application) customGPTPage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	custom, err := app.CustomGPT.GetIndividualFunction(id)
	if err != nil {
		return
	}
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, &TemplateData{
		Response:  "",
		AllButton: AllCustomGPT,
		PageLayoutData: &models.CustomGPT{
			ID:         id,
			SystemName: custom.SystemName,
		},
		PromptMessage: "",
	})
}

// User function process the given buttons system message and give response
func (app *application) customGPTFunction(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	custom, err := app.CustomGPT.GetIndividualFunction(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	prompt := r.FormValue("Message")
	userSystem := &ChatSystem{
		SystemMessage: custom.SystemPrompt,
	}

	// Call completion API
	response, err := app.OpenAIClient.GetCompletionResponse(prompt, userSystem)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	Response := mdToHTML(response)
	// render the response
	render(w, files, &TemplateData{
		AllButton: AllCustomGPT,
		PageLayoutData: &models.CustomGPT{
			SystemName: custom.SystemName,
		},
		Response:      template.HTML(Response),
		PromptMessage: prompt,
	})
}


