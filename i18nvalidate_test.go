package i18nvalidate

import (
	"strings"
	"testing"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	ptTranslations "github.com/go-playground/validator/v10/translations/pt"
)

// User is a simple DTO used across tests.
type User struct {
	FirstName string `validate:"required"       trans:"en:First Name;pt:Primeiro Nome"`
	Email     string `validate:"required,email" trans:"en:Email;pt:E-mail"`
}

// helper to build a Validator configured with English (default) and Portuguese.
func newTestValidator(t *testing.T) *Validator {
	t.Helper()

	v, err := New(
		"en",
		Translator{
			Translator:           en.New(),
			RegisterTranslations: enTranslations.RegisterDefaultTranslations,
		},
		Translator{
			Translator:           pt.New(),
			RegisterTranslations: ptTranslations.RegisterDefaultTranslations,
		},
	)
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}
	return v
}

func TestValidate_Success(t *testing.T) {
	v := newTestValidator(t)

	u := User{FirstName: "Daniel", Email: "daniel@example.com"}

	if err := v.Validate(u); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidate_RequiredEnglish(t *testing.T) {
	v := newTestValidator(t)

	err := v.Validate(User{}, "en")
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	verrs, ok := err.(*ValidationErrors)
	if !ok {
		t.Fatalf("expected *ValidationErrors, got %T", err)
	}

	if len(verrs.TranslatedErrors) != 2 {
		t.Fatalf(
			"expected 2 validation errors, got %d",
			len(verrs.TranslatedErrors),
		)
	}

	msg, exists := verrs.TranslatedErrors["User.FirstName"]
	if !exists {
		t.Fatalf("missing translated error for FirstName")
	}
	if !strings.Contains(msg, "First Name") {
		t.Errorf(
			"expected translated field name 'First Name' in message, got %q",
			msg,
		)
	}
	if strings.Contains(msg, "FirstName") {
		t.Errorf("unexpected untranslated field name in message: %q", msg)
	}
}

func TestValidate_RequiredPortuguese(t *testing.T) {
	v := newTestValidator(t)

	err := v.Validate(User{}, "pt")
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	verrs := err.(*ValidationErrors)
	msg := verrs.TranslatedErrors["User.FirstName"]

	if !strings.Contains(msg, "Primeiro Nome") {
		t.Errorf(
			"expected Portuguese translated field name in message, got %q",
			msg,
		)
	}
	if strings.Contains(msg, "FirstName") {
		t.Errorf("unexpected untranslated field name in message: %q", msg)
	}
}

func TestValidate_FallbackToDefault(t *testing.T) {
	v := newTestValidator(t)

	err := v.Validate(
		User{},
		"es",
	) // Spanish not configured – should fall back to English.
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	verrs := err.(*ValidationErrors)
	msg := verrs.TranslatedErrors["User.FirstName"]
	if !strings.Contains(msg, "First Name") {
		t.Errorf("expected fallback to English, got %q", msg)
	}
}

func TestValidate_NilInput(t *testing.T) {
	v := newTestValidator(t)

	if err := v.Validate(nil); err != nil {
		t.Fatalf("expected nil error for nil input, got %v", err)
	}
}
