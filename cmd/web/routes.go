package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *applictaion) routes() http.Handler {
	// use the http.NewServerMux() for creating the servermux
	// register the home function as the handler for the root url
	mux := http.NewServeMux()

	// create a file server which serves files out of the "static" dir.
	// path should be the relative to the project path or root directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle func to register the file server
	// which can handle all the routes starts with "/static/"
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Create a new middleware chain containing the middleware chain with the
	// handler func
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Registered the other application routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

	// authentication routes
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// protected routes
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create a middleware chain containing our middleware
	// which will be used for every request in our applictaion
	standard := alice.New(app.logRequest, commonHeaders)

	// Return the standard middleware chain followed by the servermux
	return standard.Then(mux)
}
