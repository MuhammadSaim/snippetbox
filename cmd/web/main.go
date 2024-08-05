package main

import (
	"log"
	"net/http"
)

func main(){

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
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)

	// printing the log message
	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
