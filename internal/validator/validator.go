package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// Use the regexp.MustCompile func to parse the regex for email
var EmailRegX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Define a new validator struct contains a map of validation errors
type Validator struct {
	FieldErrors map[string]string
}

// This will returns true if the FieldErrors map not contain any errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds an error message to our map
func (v *Validator) AddFieldError(key, message string) {
	// We have to initialize the map first if its not initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds an error message to the FieldErrors map
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank returns true if a value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars return true if a value contains no more then given chars
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue returns true if a value is in a list
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// MinChars return true if a value contains at least n chars
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches returns true if a value matches a provided regex
func Matches(value string, regEx *regexp.Regexp) bool {
	return regEx.MatchString(value)
}
