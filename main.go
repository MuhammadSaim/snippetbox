package main

import (
	"log"
	"net/http"
)

// define a home handler function
func home(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello from Snippetbox"))
}

// add snippetView hadnle function
func snippetView(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Display a specific snippet"))
}

// add snippetCreate hadnle function
func snippetCreate(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Display a form for creating a new snippet"))
}


func main(){
	// use the http.NewServerMux() for creating the servermux
	// register the home function as the handler for the root url
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// printing the log message
	log.Print("Starting server on :4000")

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)

}
