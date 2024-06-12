package main

import (
	"github.com/bmizerany/pat"
	// "github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.Home))

	mux.Get("/debugger", http.HandlerFunc(app.fetchDebugPage))
	mux.Post("/debugger", http.HandlerFunc(app.fetchDebug))

	mux.Get("/tester", http.HandlerFunc(app.testerPage))
	mux.Post("/tester", http.HandlerFunc(app.tester))

	mux.Get("/formatter", http.HandlerFunc(app.formatterPage))
	mux.Post("/formatter", http.HandlerFunc(app.formatter))
	// mux.Get("/handleQuery", http.HandlerFunc(app.handleQuery)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
