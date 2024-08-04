package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// define a home handler function
func home(w http.ResponseWriter, r *http.Request){

	// use the template.ParseFile function to read the template file into a
	// template set. If there's an error, we log the detailed error message, use
	// the http.Error function to send an internal server error reponse to the user
	// and the return from the handler so no subsequent code is executed.

	ts, err := template.ParseFiles("./ui/html/pages/home.tmpl")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Now we use the Execute method on the template set to write the
	// template content as the response body. The last parameter to Execute
	// represent any dynamic data that we want to pass in.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// add snippetView hadnle function
func snippetView(w http.ResponseWriter, r *http.Request){

	// extract the value of an ID from the path
	// check the ID is valid intgere
	// if it not converted into an integer or value is less then 1
	// return a 404 page not found response.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	w.Write([]byte(msg))
}

// add snippetCreate hadnle function
func snippetCreate(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Display a form for creating a new snippet"))
}
