package i18nvalidate

import (
	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// RegisterTranslationsFunc is a helper signature used to register custom
// translations for a specific validator instance and translator.
// It should call `v.RegisterTranslation` internally.
//
// Example:
//
//	func(v *validator.Validate, trans ut.Translator) error {
//	    return v.RegisterTranslation("required", trans, ...)
//	}
//
// The returned error should be propagated back to the caller of New so that
// initialisation fails fast if a translation cannot be registered.
//
// This matches the contract used by the go-playground/validator documentation.
type RegisterTranslationsFunc func(
	v *validator.Validate,
	trans ut.Translator,
) error

// Translator bundles together a locales.Translator implementation and the
// function responsible for registering all validation rule translations for
// that locale.
//
// The Translator.Translator field is the concrete locale implementation from
// github.com/go-playground/locales (e.g. en, pt, es, ...). All of the
// supported locales must be provided up-front when creating a new Validator
// via New().
//
// RegisterTranslations is invoked by New for each supplied locale and should
// register every validation rule that can be encountered by your application.
// Missing rule translations will cause the underlying validator.Validate call
// to fallback to the untranslated default error message.
//
// Example usage can be found in the README.
type Translator struct {
	Translator           locales.Translator
	RegisterTranslations RegisterTranslationsFunc
}
