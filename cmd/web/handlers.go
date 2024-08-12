package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MuhammadSaim/snippetbox/internal/models"
	"github.com/MuhammadSaim/snippetbox/internal/validator"
)

// Define a snippetCreateForm struct to represent the form data and validtaion
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// define a home handler function
func (app *applictaion) home(w http.ResponseWriter, r *http.Request) {

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
func (app *applictaion) snippetView(w http.ResponseWriter, r *http.Request) {

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
		} else {
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
func (app *applictaion) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// add snippetCreatePost handle to store the data into DB
func (app *applictaion) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// ParseForm which adds any data in POST request
	// If there is any error we will shoot clientErrors
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Declare a new empty instance of the snippetCreateForm struct
	var form snippetCreateForm

	// Call the Decode method of the form decoder, passing in the current
	// request and a pointer to our struct with the relavent fields.
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Because the Validator struct is embedded by the snippetCreateForm
	// We can call checkField directly on it tyo execute our validations
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank.")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long.")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank.")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Use the valid method to see if any of the checks failed.
	if !form.Valid() {
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
