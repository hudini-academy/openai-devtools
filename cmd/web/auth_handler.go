package main

import (
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// signUpForm renders the signup page with data from the session.
func (app *application) signUpForm(w http.ResponseWriter, r *http.Request) {
    files := []string{
        "./ui/html/signup.page.tmpl",
        "./ui/html/header.partial.tmpl",
        "./ui/html/base.layout.tmpl",
    }
    // PopString retrieves and deletes a string value from the session with the given key.
    flash := app.session.PopString(r, "flash")
    flashCategory := app.session.PopString(r, "flashCategory")
    // TemplateData is a struct used to pass data to the templates.
    render(w, files, &TemplateData{
        Flash:         flash,
        FlashCategory: flashCategory,
    })
}

// signUp handles the signup process. It collects user data, validates it, and inserts it into the database.
// It also handles error cases and redirects the user to the appropriate page.
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

    // Insert the user data into the database
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

// loginForm renders the login page with flash messages from the session.
func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
    files := []string{
        "./ui/html/login.page.tmpl",
        "./ui/html/header.partial.tmpl",
        "./ui/html/base.layout.tmpl",
    }

    // PopString retrieves and deletes a string value from the session with the given key.
    flash := app.session.PopString(r, "flash")
    flashCategory := app.session.PopString(r, "flashCategory")

    // TemplateData is a struct for passing data to templates, including flash messages and categories.
    render(w, files, &TemplateData{
        Flash:         flash,
        FlashCategory: flashCategory,
    })
}

// login handles user authentication, validation, and redirects based on success or failure.
func (app *application) login(w http.ResponseWriter, r *http.Request) {
    // Retrieve user email and password from the form data
    userEmail := r.FormValue("email")
    userPass := r.FormValue("password")

    // Validate the email address
    if !isValidEmail(userEmail) {
        // If email is invalid, set flash message and redirect to login page
        app.session.Put(r, "flash", "Invalid email address")
        app.session.Put(r, "flashCategory", "failure")
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Authenticate the user using the provided email and password
    isUser, err := app.users.Authenticate(userEmail, userPass)
    
    // If there is an error during authentication, log it and return an failure
	if err != nil {
        app.errorLog.Println(err.Error())
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        app.session.Put(r, "flash", "Incorrect email address or password")
        app.session.Put(r, "flashCategory", "failure")
        return
    }

    // If the user is authenticated, set session variables and redirect to home page
    if isUser != -1 {
        app.session.Put(r, "Authenticated", true)
        app.session.Put(r, "username", isUser)
        http.Redirect(w, r, "/", http.StatusSeeOther)
    } else {
        // If the user is not authenticated, set flash message and redirect to login page
        app.session.Put(r, "flash", "Login failed")
        app.session.Put(r, "flashCategory", "failure")
        http.Redirect(w, r, "/login", http.StatusSeeOther)

        // Clear the authenticated session variable
        app.session.Put(r, "Authenticated", false)
    }
}

// logoutUser clears the user's session and redirects to the home page.
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
    // Clear the "Authenticated" session variable
    app.session.Put(r, "Authenticated", false)

    // Redirect the user to the home page
    http.Redirect(w, r, "/", http.StatusSeeOther)
}