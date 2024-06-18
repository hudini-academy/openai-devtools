package main

import (

	"net/http"
)

// AuthenticateMiddleware is a middleware function that checks if a user is authenticated.
// If the user is not authenticated, it redirects the user to the login page and stores a flash message.
// If the user is authenticated, it allows the request to proceed to the next handler.
func (app *application) AuthenticateMiddleware(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        // Check if the user is authenticated using the session store
        if !app.session.GetBool(r, "Authentication") {
            // Store a flash message in the session
            app.session.Put(r, "flash", "Log In Before Accessing the resources")
            // Redirect the user to the login page
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }
        // If the user is authenticated, allow the request to proceed to the next handler
        next.ServeHTTP(w, r)
    }
    // Return the wrapped handler function as an http.Handler
    return http.HandlerFunc(fn)
}
