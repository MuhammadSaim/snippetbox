package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/MuhammadSaim/snippetbox/internal/models"
)

// define a home handler function
func (app *applictaion) home(w http.ResponseWriter, r *http.Request){

	// initialize a slice containing the paths
	// it is important our base file should add on top
	files := []string{
		"./ui/html/layouts/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// use the template.ParseFile function to read the template file into a
	// template set. If there's an error, we log the detailed error message, use
	// the http.Error function to send an internal server error reponse to the user
	// and the return from the handler so no subsequent code is executed.

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		app.serverError(w, r, err)
		return
	}

	// Now we use the Execute method on the template set to write the
	// template content as the response body. The last parameter to Execute
	// represent any dynamic data that we want to pass in.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		app.serverError(w, r, err)
	}
}

// add snippetView hadnle function
func (app *applictaion) snippetView(w http.ResponseWriter, r *http.Request){

	// extract the value of an ID from the path
	// check the ID is valid intgere
	// if it not converted into an integer or value is less then 1
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the SnippetModel's Get to find the snippet and send it to the response
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		}else{
			app.serverError(w, r, err)
		}
	}

	// Write the snippet data as a plain-text HTTP response
	fmt.Fprintf(w, "%+v", snippet)
}

// add snippetCreate hadnle function
func (app *applictaion) snippetCreate(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Display a form for creating a new snippet"))
}

// add snippetCreatePost handle to store the data into DB
func (app *applictaion) snippetCreatePost(w http.ResponseWriter, r *http.Request)  {
	// create some variable to memic post data
	title := "My First Snippets"
	content := `{{ define "title" }} Home {{ end }}
				{{ define "main" }}
				<h2>Latest Snippets</h2>
				<p>There's nothing to see here yet!</p>
				{{ end }}
				`
	expiresIn := 7

	// pass this data to Insert method to store in the DB
	id, err := app.snippets.Insert(title, content, expiresIn)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// redirect the user to the snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
