package main

// Import the Required packages
import (
	"OpenAIDevTools/pkg/models"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Initialize the Response as a string
var Response string

// Initialize the customGPT Struct slice which has all the CustomGPT fields
var AllCustomGPT []*models.CustomGPT

// Render homepage, base function checks user authentication, retrieves details,
// fetches custom GPT functions, renders home with user's name and functions;
// redirects to landing if not authenticated.
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

// landingPage renders the landing page with a welcome message and links to login/signup,
// accessible to all users. Params: w http.ResponseWriter, r *http.Request
func (app *application) landingPage(w http.ResponseWriter, r *http.Request) {
	// Define the list of template files to be rendered
	files := []string{
		"./ui/html/landing.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/hero.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	// Call the render function to render the landing page template
	// Pass nil as the data for the template, as the landing page does not require any dynamic data
	render(w, files, nil)
}

// renderCustomGPT renders page for new custom GPT function creation;
// authenticated users only. Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) renderCustomGPT(w http.ResponseWriter, r *http.Request) {
	// Define the list of template files to be rendered
	files := []string{
		"./ui/html/new.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	// Call the render function to render the new custom GPT function page template
	// Pass nil as the data for the template, as the page does not require any dynamic data
	render(w, files, nil)
}

// createCustomGPT handles creation of new custom GPT function, retrieves form data,
// prints debug info, inserts function via CustomGPT service;
// logs errors on failure, returns HTTP 500 on error, redirects to home with HTTP 303 otherwise.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) createCustomGPT(w http.ResponseWriter, r *http.Request) {
	ButtonName := r.FormValue("FunctionName")
	newMessage := r.FormValue("FunctionMessage")

	fmt.Println(ButtonName, newMessage)

	// fmt.Println(ButtonName, newMessage)
	err := app.CustomGPT.InsertFunction(ButtonName, newMessage)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// deleteCustomGPT handles deletion of custom GPT function, retrieves function ID from request data,
// deletes function via CustomGPT service; logs errors on failure,
// returns HTTP 500 on error, redirects to home with HTTP 303 otherwise.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) deleteCustomGPT(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ID of the function to be deleted from the form data
	id, _ := strconv.Atoi(r.FormValue("id"))

	// Print the ID for debugging purposes
	fmt.Println(id)

	// Call the DeleteFunction method of the CustomGPT service to delete the function
	err := app.CustomGPT.DeleteFunction(id)

	// If an error occurs during the deletion process, log the error and return an HTTP 500 Internal Server Error response
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Redirect the user to the home page with a HTTP 303 See Other status code
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// customGPTPage handles rendering of custom GPT function page, retrieves function ID,
// fetches details via CustomGPT service; renders template on success,
// passes TemplateData with response, functions, ID, system name, and prompt.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) customGPTPage(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ID of the function from the form data
	id, _ := strconv.Atoi(r.FormValue("id"))

	// Fetch the details of the function from the database
	custom, err := app.CustomGPT.GetIndividualFunction(id)
	if err != nil {
		return
	}

	// Define the list of template files to be rendered
	files := []string{
		"./ui/html/response.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	// Render the custom GPT function page template
	// Pass a TemplateData struct containing the response, all custom GPT functions,
	// the ID and system name of the selected function, and an empty prompt message
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

// customGPTFunction handles execution of custom GPT function, retrieves function ID,
// fetches details via CustomGPT service; logs error on retrieval failure, gets prompt message,
// creates ChatSystem with function prompt, calls OpenAIClient for completion response, logs API errors,
// renders response template on success with TemplateData containing functions, system name, response, and prompt.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
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

// signUpForm handles user signup form rendering, retrieves flash messages from session,
// renders signup form template with flash messages.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) signUpForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/signup.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	// Retrieve flash messages from session
	flash := app.session.PopString(r, "flash")
	flashCategory := app.session.PopString(r, "flashCategory")

	// Render the signup form template with flash messages
	render(w, files, &TemplateData{
		Flash:         flash,
		FlashCategory: flashCategory,
	})
}

// signUp handles user signup, retrieves form data, validates email,
// hashes password, checks if user already exists, validates password,
// inserts user into Users service; logs error on insertion failure,
// returns HTTP 500 on error, redirects to login with HTTP 303 on successful signup.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("name")
	userEmail := r.FormValue("email")
	password := r.FormValue("password")

	// Hash the password using bcrypt
	hashedPassword, errHashing := bcrypt.GenerateFromPassword([]byte(password), 12)
	if errHashing != nil {
		log.Println(errHashing)
		return
	}

	// Validate the email address
	if !isValidEmail(userEmail) {
		app.session.Put(r, "flash", "Invalid email address")
		app.session.Put(r, "flashCategory", "failure")
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	// Check if the user already exists in the database
	if app.users.UserExist(userEmail) {
		app.session.Put(r, "flash", "User already exists")
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	// Validate the password
	if err := isValidPassword(password); err != nil {
		app.session.Put(r, "flash", err.Error())
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	// Insert the user into the database
	if len(strings.TrimSpace(userName)) != 0 && len(strings.TrimSpace(userEmail)) != 0 && len(strings.TrimSpace(password)) != 0 {
		err := app.users.Insert(userName, userEmail, hashedPassword)
		if err != nil {
			app.errorLog.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
		}
		app.session.Put(r, "flash", "SignUp Successful")
		app.session.Put(r, "flashCategory", "success")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		app.session.Put(r, "flash", " item field cant be empty!")
		app.session.Put(r, "flashCategory", "failure")
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
	}
}

// loginForm handles user login form rendering, retrieves flash messages from session,
// renders login form template with flash messages.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/login.page.tmpl",
		"./ui/html/header.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	// Retrieve flash messages from session
	flash := app.session.PopString(r, "flash")
	flashCategory := app.session.PopString(r, "flashCategory")

	// Render the login form template with flash messages
	render(w, files, &TemplateData{
		Flash:         flash,
		FlashCategory: flashCategory,
	})
}

// login handles user login, retrieves form data, validates email,
// authenticates user via Users service; logs error on authentication failure,
// returns HTTP 500 on error, redirects to home with HTTP 303 on successful login.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	userEmail := r.FormValue("email")
	userPass := r.FormValue("password")

	// Validate the email address
	if !isValidEmail(userEmail) {
		app.session.Put(r, "flash", "Invalid email address")
		app.session.Put(r, "flashCategory", "failure")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Authenticate the user using the Users service
	isUser, err := app.users.Authenticate(userEmail, userPass)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	// If the user is authenticated, set the "Authenticated" flag to true and redirect to the home page
	if isUser != -1 {
		app.session.Put(r, "Authenticated", true)
		app.session.Put(r, "username", isUser)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// If the user is not authenticated, set the "flash" message for login failure and redirect to the login page
		app.session.Put(r, "flash", "Login failed")
		app.session.Put(r, "flashCategory", "failure")
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		// Clear the "Authenticated" flag
		app.session.Put(r, "Authenticated", false)
	}
}

// logoutUser handles user logout, invalidates session, clears "Authenticated" flag, redirects to home page.
// Params: w http.ResponseWriter, r *http.Request. Returns: None
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Clear the "flash" message for logout success
	// app.session.Put(r, "flash", "Logout Success")

	// Invalidate the user's session by setting the "Authenticated" flag to false
	app.session.Put(r, "Authenticated", false)

	// Redirect the user to the home page with a HTTP 303 See Other status code
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
