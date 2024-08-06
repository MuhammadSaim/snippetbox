package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type applictaion struct {
	logger *slog.Logger
}

func main(){

	// Define a new command-line flag with the name 'addr', a default
	// value of ":4000"
	addr := flag.String("addr", ":4000", "HTTP network address")


	// Importantly, we use flag parse func to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr variable
	flag.Parse()

	// use the slog.New func to initialize a new structured logger, which
	// writes to the standard out stream and use the default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize a new instance of our application struct, containing the
	// dependencies for the time being just adding our logger

	app := &applictaion{
		logger: logger,
	}

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

	// use the Info() method to log the starting server message at info
	logger.Info("Starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)

	// And we also use the Error() method to lag any error message returned by
	// http.ListenAndServe() at Error. End of that terminate the application with os.Exit
	logger.Error(err.Error())

	// exit the application
	os.Exit(1)
}
