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

	// use the Info() method to log the starting server message at info
	logger.Info("Starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, app.routes())

	// And we also use the Error() method to lag any error message returned by
	// http.ListenAndServe() at Error. End of that terminate the application with os.Exit
	logger.Error(err.Error())

	// exit the application
	os.Exit(1)
}
