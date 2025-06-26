package i18nvalidate

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationErrors wraps validator.ValidationErrors and provides a joined
// string representation where each field error is already translated.
// The TranslatedErrors map contains an entry for every struct field namespace
// with its corresponding translated message.
type ValidationErrors struct {
	validator.ValidationErrors
	TranslatedErrors map[string]string
}

// Error implements the error interface, returning a semicolon-separated list
// of translated validation error messages.
func (v *ValidationErrors) Error() string {
	msgs := make([]string, 0, len(v.TranslatedErrors))
	for _, msg := range v.TranslatedErrors {
		msgs = append(msgs, msg)
	}
	return strings.Join(msgs, "; ")
}
