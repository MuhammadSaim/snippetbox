package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/MuhammadSaim/snippetbox/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML.
type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// Create a humanDate func which returns a nicely formatted
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template func and the func themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// A function to implement the cache
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a map to memic as a cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob func to get a slice of all filepaths
	// match the pattren "./ui/html/pages/*.tmpl". This will essentially gives
	// slice of all the filepaths for our applictaion page templates
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// loop through the pages filepaths one by one
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath
		name := filepath.Base(page)

		// Parse the base template file into a template set
		// for using the functions We have to create empty template set and register the
		// functions before the ParseFiles function
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/layouts/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Use ParseGlob on this template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Parse the files into a template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the tempate set to the map, using the name of the page
		// like 'home.tmpl' as the key
		cache[name] = ts
	}

	// return the map
	return cache, nil
}
