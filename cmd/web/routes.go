package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)
// routes sets up the application's routing system using the pat package.
// It returns an http.Handler that can be used to serve HTTP requests.
func (app *application) routes() http.Handler {
    mux := pat.New()

    // Serve the base page at the root URL ("/")
    mux.Get("/", app.session.Enable(http.HandlerFunc(app.Base)))

    // Render the custom GPT page
    mux.Get("/rendercustomgpt", http.HandlerFunc(app.renderCustomGPT))

    // Create a new custom GPT
    mux.Post("/createcustomgpt", http.HandlerFunc(app.createCustomGPT))

    // Delete a custom GPT
    mux.Post("/deletecustomgpt", http.HandlerFunc(app.deleteCustomGPT))

    // Serve the custom GPT page
    mux.Get("/customgpt", http.HandlerFunc(app.customGPTPage))

    // Perform a function on the custom GPT
    mux.Post("/customgpt", http.HandlerFunc(app.customGPTFunction))

    // Serve the login form
    mux.Get("/login", app.session.Enable(http.HandlerFunc(app.loginForm)))

    // Handle login
    mux.Post("/login", app.session.Enable(http.HandlerFunc(app.login)))

    // Serve the sign up form
    mux.Get("/signup", app.session.Enable(http.HandlerFunc(app.signUpForm)))

    // Handle sign up
    mux.Post("/signup", app.session.Enable(http.HandlerFunc(app.signUp)))

    // Handle logout
    mux.Get("/logout", app.session.Enable(http.HandlerFunc(app.logoutUser)))

    // Serve static files from the "./ui/static/" directory
    fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Get("/static/", http.StripPrefix("/static", fileServer))

    return mux
}