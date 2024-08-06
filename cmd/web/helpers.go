package main

import (
	"net/http"
)

// The serverError helper writes a log entry at Error level
// then sends a generic 500 internal server error response to the user
func (app *applictaion) serverError(w http.ResponseWriter, r *http.Request, err error){
	var (
		method = r.Method
		uri = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
func (app *applictaion) clientError(w http.ResponseWriter, status int){
	http.Error(w, http.StatusText(status), status)
}
