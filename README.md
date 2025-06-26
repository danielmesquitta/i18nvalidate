# i18nvalidate

Internationalized, struct-based input validation for Go ✨  
Built on top of [go-playground/validator](https://github.com/go-playground/validator) and [go-playground/universal-translator](https://github.com/go-playground/universal-translator), `i18nvalidate` makes it trivial to return user-friendly validation errors in multiple languages.

---

## Features

- 🔠 **Multi-locale** – translate error messages and field names with a single function call.
- 🏷️ **`trans` struct tag** – describe localised field names right where you define your DTOs.
- 🛠️ **Drop-in replacement** for `validator.Validate`: keep all existing validation tags.
- 🚦 **Graceful fallback** – if the requested locale is missing the package falls back to the default language.

---

## Installation

```bash
go get github.com/danielmesquitta/i18nvalidate
```

The library requires **Go 1.24+** (the same minimum version as `validator/v10`).

---

## Quick-start

```go
package main

import (
    "fmt"

    "github.com/danielmesquitta/i18nvalidate"
    "github.com/go-playground/locales/en"
    "github.com/go-playground/locales/pt"
    enTranslations "github.com/go-playground/validator/v10/translations/en"
    ptTranslations "github.com/go-playground/validator/v10/translations/pt"
)

// User is the data we want to validate.
// The normal `validate` tags still apply, plus the optional `trans` tag
// where you can supply human-readable field names per language.
type User struct {
    FirstName string `validate:"required" trans:"en:First Name;pt:Primeiro Nome"`
    Email     string `validate:"required,email" trans:"pt:E-mail"`
}

func main() {
    v, _ := i18nvalidate.New(
      "en",
      i18nvalidate.Translator{Translator: en.New(), RegisterTranslations: enTranslations.RegisterDefaultTranslations},
      i18nvalidate.Translator{Translator: pt.New(), RegisterTranslations: ptTranslations.RegisterDefaultTranslations},
    )

    err := v.Validate(User{}, "en") // request English messages
    if err != nil {
        fmt.Println(err.Error())
        // Output: "First Name is a required field; Email is a required field"
    }
}
```

---

## API

### `New(default string, translators ...Translator) (*Validator, error)`

Use this constructor when you want full control over the locales you bundle:

```go
import (
    "fmt"

    "github.com/danielmesquitta/i18nvalidate"
    "github.com/go-playground/locales/fr"
    "github.com/go-playground/locales/de"
    frTranslations "github.com/go-playground/validator/v10/translations/fr"
    deTranslations "github.com/go-playground/validator/v10/translations/de"
)

v, err := i18nvalidate.New("fr",
    i18nvalidate.Translator{Translator: fr.New(), RegisterTranslations: frTranslations.RegisterDefaultTranslations},
    i18nvalidate.Translator{Translator: de.New(), RegisterTranslations: deTranslations.RegisterDefaultTranslations},
)
```

### `Validate(data any, lang ...string) error`

Validates `data` (a struct pointer or value). If validation fails, an `error` whose message contains the translated errors joined by `; ` is returned. Pass the IETF BCP-47 language tag you want, e.g. `"es"`. If omitted it falls back to the default language supplied to `New`.

---

## `trans` tag syntax

The optional `trans` struct tag lets you provide readable field names per locale:

```go
//        ┌─language code
//        │      ┌──translated name
//        │      │
// trans:"en:User Name;pt:Nome do Usuário;es:Nombre de Usuario"
```

Rules:

1. Multiple translations are separated by `;`.
2. Each translation is `<lang>:<value>`.
3. Whitespace around tokens is ignored.
4. Missing languages are simply skipped; the field name stays unchanged for those locales.

---

## License

Distributed under the MIT License. See `LICENSE` for more information.
