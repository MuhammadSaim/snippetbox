package main

import (
	"html/template"
	"path/filepath"

	"github.com/MuhammadSaim/snippetbox/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML.
type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
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
		ts, err := template.ParseFiles("./ui/html/layouts/base.tmpl")
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
