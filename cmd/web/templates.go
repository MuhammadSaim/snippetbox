package main

import "github.com/MuhammadSaim/snippetbox/internal/models"

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML.
type templateData struct {
	Snippet models.Snippet
	Snippets []models.Snippet
}
