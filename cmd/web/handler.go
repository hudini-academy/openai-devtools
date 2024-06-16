package main

import (
	"OpenAIDevTools/pkg/models"
	"fmt"
	"html/template"
	"log"
	"strings"

	"net/http"
	"strconv"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var Response string
var AllCustomGPT []*models.CustomGPT

// Renders homepage
func (app *application) Base(w http.ResponseWriter, r *http.Request) {

	// Check if user id present in session
	// if not, show landing page
	app.infoLog.Println("Test")
	isAuth, ok := app.session.Get(r, "Authenticated").(bool)
	if !ok {
		app.errorLog.Println("Error getting session")
		http.Error(w, "Internal Server Error", 500)
	}

	// If user is not authenticated, show landing page
	if !isAuth {
		app.landingPage(w, r)
		return
	}

	// if preset, show the dashboard
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	var errr error
	AllCustomGPT, errr = app.CustomGPT.GetFunction()
	if errr != nil {
		return
	}

	render(w, files, nil)
}

// isValidPassword checks if the given password meets the required criteria
func isValidPassword(password string) error {

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if len(password) < 8 || !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return fmt.Errorf("password must contain 8 characters including an uppercase, lowercase, digit and special character")
	}

	return nil
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
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, nil)
}

// Create a button by getting the name and message
func (app *application) createCustomGPT(w http.ResponseWriter, r *http.Request) {
	ButtonName := r.FormValue("FunctionName")
	newMessage := r.FormValue("FunctionMessage")

	// fmt.Println(ButtonName, newMessage)
	err := app.CustomGPT.InsertFunction(ButtonName, newMessage)
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
		"./ui/html/base.layout.tmpl",
	}

	render(w, files, &TemplateData{
		Response:  "",
		AllButton: AllCustomGPT,
		PageLayoutData: &models.CustomGPT{
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

	// fmt.Println("Response:", Response)
}

// func (app *application) handleQuery(w http.ResponseWriter, r *http.Request) {

// }

func (app *application) signUpForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/signup.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}
	flash := app.session.PopString(r, "flash")
	flashCategory := app.session.PopString(r, "flashCategory")
	render(w, files, &TemplateData{
		Flash:         flash,
		FlashCategory: flashCategory,
	})
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("name")
	userEmail := r.FormValue("email")
	password := r.FormValue("password")
	hashedPassword, errHashing := bcrypt.GenerateFromPassword([]byte(password), 12)
	if errHashing != nil {
		log.Println(errHashing)
		return
	}

	if !isValidEmail(userEmail) {
		app.session.Put(r, "flash", "Invalid email address")
		app.session.Put(r, "flashCategory", "failure")
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
		return
	}
	if app.users.UserExist(userEmail) {
		app.session.Put(r, "flash", "User already exists")
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
		return
	}
	if err := isValidPassword(password); err != nil {
		app.session.Put(r, "flash", err.Error())
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
		return
	}

	if len(strings.TrimSpace(userName)) != 0 && len(strings.TrimSpace(userEmail)) != 0 && len(strings.TrimSpace(password)) != 0 {
		err := app.users.Insert(userName, userEmail, hashedPassword)
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
		}

		app.session.Put(r, "flash", "SignUp Successful")
		app.session.Put(r, "flashCategory", "success")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	} else {
		app.session.Put(r, "flash", " item field cant be empty!")
		app.session.Put(r, "flashCategory", "failure")
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
	}

}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/login.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	flash := app.session.PopString(r, "flash")
	flashCategory := app.session.PopString(r, "flashCategory")
	render(w, files, &TemplateData{
		Flash:         flash,
		FlashCategory: flashCategory,
	})
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	userEmail := r.FormValue("email")
	userPass := r.FormValue("password")

	if !isValidEmail(userEmail) {
		app.session.Put(r, "flash", "Invalid email address")
		app.session.Put(r, "flashCategory", "failure")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	isUser, err := app.users.Authenticate(userEmail, userPass)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	if isUser != 0 { // TODO: VISHWA: need to check if this is correct logic
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "username", isUser)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		app.session.Put(r, "flash", "Login failed")
		app.session.Put(r, "flashCategory", "failure")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		app.session.Put(r, "Authenticated", false)
	}
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// app.session.Put(r, "flash", "Logout Success")
	app.session.Put(r, "Authenticated", false)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
