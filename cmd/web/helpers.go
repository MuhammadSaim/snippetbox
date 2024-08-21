package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/MuhammadSaim/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// The serverError helper writes a log entry at Error level
// then sends a generic 500 internal server error response to the user
func (app *applictaion) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
func (app *applictaion) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Render the template but first find out in the cache
func (app *applictaion) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
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
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.IsAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

// Create a new decodePostForm helper method. The second param is dst
// is the target destination that we want to decode the form data into
func (app *applictaion) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// snippetCreateForm handler
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode on our decoder instance, passing the target destination as the
	// first param
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decoder method
		// will return an error
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal
		return err
	}

	return nil
}

// Return true if the user is authenticated
func (app *applictaion) IsAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

// generate unique id for the snippets
func (app *applictaion) GenerateUniqueID() (string, error) {
	const maxAttempts = 1000
	for i := 1; i <= maxAttempts; i++ {
		uniqueID := GetUniqueBase64ID()
		_, err := IDExists(app, uniqueID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return uniqueID, nil
			}
			return "", err
		}
	}
	return "", fmt.Errorf("failed to generate the unique id")
}

// check the unique id in DB
func IDExists(app *applictaion, ID string) (bool, error) {
	_, err := app.snippets.Get(ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Get the unique id
func GetUniqueBase64ID() string {
	base64_chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

	// seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	sb.Grow(len(base64_chars))

	for i := 0; i < 11; i++ {
		rand_char := base64_chars[r.Intn(len(base64_chars))]
		sb.WriteByte(rand_char)
	}

	return sb.String()
}
