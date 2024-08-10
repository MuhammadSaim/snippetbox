package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MuhammadSaim/snippetbox/internal/models"
)

// define a home handler function
func (app *applictaion) home(w http.ResponseWriter, r *http.Request){

	// Fetch the latest snippets from the DB
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Use the render func
	app.render(w, r, http.StatusOK, "home.tmpl", templateData{
		Snippets: snippets,
	})

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
		return
	}

	// Use the render func
	app.render(w, r, http.StatusOK, "view.tmpl", templateData{
		Snippet: snippet,
	})

}

// add snippetCreate hadnle function
func (app *applictaion) snippetCreate(w http.ResponseWriter, r *http.Request){
	app.render(w, r, http.StatusOK, "create.tmpl", templateData{})
}

// add snippetCreatePost handle to store the data into DB
func (app *applictaion) snippetCreatePost(w http.ResponseWriter, r *http.Request)  {

	// ParseForm which adds any data in POST request
	// If there is any error we will shoot clientErrors
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// use the Get func to chain with ParseForm and get the data
	title := r.FormValue("title")
	content := r.FormValue("content")

	// FormValue return data in string so we have to for the
	// expire values we have to convert it into Int
	expiresIn, err := strconv.Atoi(r.FormValue("expires"))
	if err !=  nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// pass this data to Insert method to store in the DB
	id, err := app.snippets.Insert(title, content, expiresIn)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// redirect the user to the snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
