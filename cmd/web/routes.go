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

	// Registered the other application routes
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing our middleware
	// which will be used for every request in our applictaion
	standard := alice.New(app.logRequest, commonHeaders)

	// Return the standard middleware chain followed by the servermux
	return standard.Then(mux)
}
