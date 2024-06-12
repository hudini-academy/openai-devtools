package main

import (
	"context"
	"net/http"
)

func (app *application) authenticateMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !app.session.GetBool(r, "Authenticated") {
			app.session.Put(r, "flash", "Log In Before Accessing the resources")
			http.Redirect(w, r, "/user/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
func (app *application) loadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := app.session.Get(r, "user").(string)
		if ok && user != "" {
			ctx := context.WithValue(r.Context(), "user", user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}