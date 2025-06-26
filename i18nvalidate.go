package i18nvalidate

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Validator is the main entry-point for performing struct validation with
// translated error messages.
//
// A Validator instance is safe for concurrent use by multiple goroutines.
type Validator struct {
	Validator       *validator.Validate
	Uni             *ut.UniversalTranslator
	translators     map[string]ut.Translator
	defaultLang     string
	registeredTypes sync.Map // map[reflect.Type]struct{}
}

// New constructs a new Validator configured with the provided locales.
//
// At least one Translator must be supplied, and one of them must match the
// defaultLang parameter.
func New(
	defaultLang string,
	translators ...Translator,
) (*Validator, error) {
	if len(translators) == 0 {
		return nil, fmt.Errorf("at least one locale must be supplied")
	}

	var fallback Translator
	for _, t := range translators {
		if t.Translator.Locale() == defaultLang {
			fallback = t
			break
		}
	}
	if fallback.Translator == nil {
		return nil, fmt.Errorf(
			"default language '%s' not found in supplied locales",
			defaultLang,
		)
	}

	args := make([]locales.Translator, 0, len(translators)+1)
	args = append(args, fallback.Translator)
	for _, t := range translators {
		args = append(args, t.Translator)
	}
	uni := ut.New(args[0], args...)

	v := validator.New()

	translatorMap := make(map[string]ut.Translator, len(translators))
	for _, t := range translators {
		trans, found := uni.GetTranslator(t.Translator.Locale())
		if !found {
			continue
		}
		translatorMap[t.Translator.Locale()] = trans

		if err := t.RegisterTranslations(v, trans); err != nil {
			return nil, err
		}
	}

	return &Validator{
		Validator:   v,
		Uni:         uni,
		translators: translatorMap,
		defaultLang: defaultLang,
	}, nil
}

// Validate performs struct validation and returns translated error messages.
//
// If the provided data is nil or not a struct/pointer to struct, the method
// exits early doing nothing. When validation errors occur they are returned as
// *ValidationErrors so callers can inspect individual field messages.
func (v *Validator) Validate(data any, lang ...string) error {
	if data == nil {
		return nil
	}

	if err := v.registerFieldTranslations(data); err != nil {
		return err
	}

	targetLang := v.defaultLang
	if len(lang) > 0 && lang[0] != "" {
		targetLang = lang[0]
	}

	trans, ok := v.translators[targetLang]
	if !ok {
		trans = v.translators[v.defaultLang]
	}

	if err := v.Validator.Struct(data); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			msgs := make(map[string]string, len(verrs))
			for _, fe := range verrs {
				msg := fe.Translate(trans)
				fieldTranslated, err := trans.T(fe.StructNamespace())
				if err == nil && fieldTranslated != "" {
					msg = strings.Replace(
						msg,
						fe.StructField(),
						fieldTranslated,
						1,
					)
				}
				msgs[fe.StructNamespace()] = msg
			}
			return &ValidationErrors{
				ValidationErrors: verrs,
				TranslatedErrors: msgs,
			}
		}
		return err
	}

	return nil
}

// registerFieldTranslations stores field name translations found in struct
// tags so that they can be used when building the final error messages.
func (v *Validator) registerFieldTranslations(data any) error {
	t := reflect.TypeOf(data)

	// Dereference pointers so we always deal with the concrete struct type.
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}

	if _, loaded := v.registeredTypes.LoadOrStore(t, struct{}{}); loaded {
		return nil
	}

	return v.registerFieldTranslationsType(t)
}

func (v *Validator) registerFieldTranslationsType(
	t reflect.Type,
) error {
	if t == nil {
		return nil
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if tagValue := field.Tag.Get("trans"); tagValue != "" {
			parts := strings.Split(tagValue, ";")
			for _, part := range parts {
				kv := strings.SplitN(strings.TrimSpace(part), ":", 2)
				if len(kv) != 2 {
					continue
				}

				langCode := strings.TrimSpace(kv[0])
				translatedName := strings.TrimSpace(kv[1])

				if trans, ok := v.translators[langCode]; ok {
					key := t.Name() + "." + field.Name
					if err := trans.Add(key, translatedName, true); err != nil {
						return err
					}
				}
			}
		}

		nestedType := field.Type
		if nestedType.Kind() == reflect.Pointer {
			nestedType = nestedType.Elem()
		}
		if nestedType.Kind() == reflect.Struct {
			if err := v.registerFieldTranslationsType(nestedType); err != nil {
				return err
			}
		}
	}

	return nil
}
