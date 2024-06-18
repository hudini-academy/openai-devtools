package main

import (
	"github.com/bmizerany/pat"
	// "github.com/justinas/alice"
	"net/http"
)

// routes sets up the application's routing system using the pat package.
// It returns an http.Handler that can be used to serve HTTP requests.
func (app *application) routes() http.Handler {
	// Create a new pat router.
	mux := pat.New()

	// Define a route for the root ("/") URL.
	// The Base method is called when a GET request is made to this route.
	mux.Get("/", app.session.Enable(http.HandlerFunc(app.Base)))

	// Define routes for handling custom GPT operations.
	mux.Get("/rendercustomgpt", http.HandlerFunc(app.renderCustomGPT))
	mux.Post("/createcustomgpt", http.HandlerFunc(app.createCustomGPT))
	mux.Post("/deletecustomgpt", http.HandlerFunc(app.deleteCustomGPT))

	// Define routes for the custom GPT page.
	mux.Get("/customgpt", http.HandlerFunc(app.customGPTPage))
	mux.Post("/customgpt", http.HandlerFunc(app.customGPTFunction))

	// Define routes for user authentication (login, signup, logout).
	mux.Get("/login", app.session.Enable(http.HandlerFunc(app.loginForm)))
	mux.Post("/login", app.session.Enable(http.HandlerFunc(app.login)))
	mux.Get("/signup", app.session.Enable(http.HandlerFunc(app.signUpForm)))
	mux.Post("/signup", app.session.Enable(http.HandlerFunc(app.signUp)))
	mux.Get("/logout", app.session.Enable(http.HandlerFunc(app.logoutUser)))

	// Serve static files from the "./ui/static/" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Return the configured router.
	return mux
}
