package main

import (
	"github.com/bmizerany/pat"
	// "github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := pat.New()

	mux.Get("/", app.session.Enable(http.HandlerFunc(app.Base)))

	mux.Get("/rendercustomgpt", http.HandlerFunc(app.renderCustomGPT))
	mux.Post("/createcustomgpt", http.HandlerFunc(app.createCustomGPT))

	mux.Post("/deletecustomgpt", http.HandlerFunc(app.deleteCustomGPT))

	mux.Get("/customgpt", http.HandlerFunc(app.customGPTPage))
	mux.Post("/customgpt", http.HandlerFunc(app.customGPTFunction))

	// Login and Sign up
	mux.Get("/login", app.session.Enable(http.HandlerFunc(app.loginForm)))
	mux.Post("/login", app.session.Enable(http.HandlerFunc(app.login)))
	mux.Get("/signup", app.session.Enable(http.HandlerFunc(app.signUpForm)))
	mux.Post("/signup", app.session.Enable(http.HandlerFunc(app.signUp)))
	mux.Get("/logout", app.session.Enable(http.HandlerFunc(app.logoutUser)))

	// mux.Get("/handleQuery", http.HandlerFunc(app.handleQuery)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
