package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
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

// Render the template but first find out in the cache
func (app *applictaion) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData)  {
	// Retrive the appropriate template set from the cache based on the page
	// If there is no entry in the cache will through an error
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the %s template does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer without any errors, we are safe to go
	// http.ResponseWritter. If there's an error, call our server error helper
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter.
	// Note: This is another time where we pass our http.ResponseWriter to a function that takes an io.Writer
	buf.WriteTo(w)

}

// Create a newTemplateData helper
func (app *applictaion) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}
