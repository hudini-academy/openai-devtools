package main

import (

	"net/http"
)

func (app *application) AuthenticateMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !app.session.GetBool(r, "Authentication") {
			app.session.Put(r, "flash", "Log In Before Accessing the resources")
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
