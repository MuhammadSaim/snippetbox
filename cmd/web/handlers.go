package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/MuhammadSaim/snippetbox/internal/models"
)

// Define a snippetCreateForm struct to represent the form data and validtaion
type snippetCreateForm struct {
	Title string
	Content string
	Expires int
	FieldErrors map[string]string
}

// define a home handler function
func (app *applictaion) home(w http.ResponseWriter, r *http.Request){

	// Fetch the latest snippets from the DB
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the render func
	app.render(w, r, http.StatusOK, "home.tmpl", data)

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

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// Use the render func
	app.render(w, r, http.StatusOK, "view.tmpl", data)

}

// add snippetCreate hadnle function
func (app *applictaion) snippetCreate(w http.ResponseWriter, r *http.Request){
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
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

	// FormValue return data in string so we have to for the
	// expire values we have to convert it into Int
	expiresIn, err := strconv.Atoi(r.FormValue("expires"))
	if err !=  nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title: r.FormValue("title"),
		Content: r.FormValue("content"),
		Expires: expiresIn,
		FieldErrors: map[string]string{},
	}

	// check the title value is not blank and is note more then 100 chars
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank."
	}else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long."
	}

	// Check the content value isn't blank
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank."
	}

	// Check the expires value matches one of these values 1, 7 or 365
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "The field must equal 1, 7 or 365"
	}

	// If there is any validation error then we have to re-render the view
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// pass this data to Insert method to store in the DB
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// redirect the user to the snippet page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
