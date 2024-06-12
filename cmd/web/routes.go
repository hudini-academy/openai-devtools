package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", (http.HandlerFunc(app.Base)))
	mux.Get("/home", app.session.Enable(http.HandlerFunc(app.Home)))
	mux.Get("/", app.loadUser(http.HandlerFunc(app.Home)))

	mux.Get("/debugger", app.session.Enable(http.HandlerFunc(app.fetchDebugPage)))
	mux.Post("/debugger", app.session.Enable(http.HandlerFunc(app.fetchDebug)))

	// Login and Sign up
	mux.Get("/user/login", app.session.Enable(http.HandlerFunc(app.loginForm)))
	mux.Post("/user/login", app.session.Enable(http.HandlerFunc(app.login)))
	mux.Get("/user/signup", app.session.Enable(http.HandlerFunc(app.signUpForm)))
	mux.Post("/user/signup", app.session.Enable(http.HandlerFunc(app.signUp)))
	mux.Get("/user/logout", app.session.Enable(http.HandlerFunc(app.logoutUser)))




	// mux.Get("/handleQuery", http.HandlerFunc(app.handleQuery)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
