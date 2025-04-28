package main

import (
	"embed"
	"encoding/json"
	"log"
	"strings"

	"github.com/jeandeaual/go-locale"
)

//go:embed i18n
var locales embed.FS

// all locales are, at this time, stored in a map
// TODO: maybe improve this
var translations = map[string]map[string]string{}

func init() {
	dirs, err := locales.ReadDir("i18n")
	if err != nil {
		panic(err)
	}
	for _, f := range dirs {
		if f.IsDir() {
			continue
		}

		source, err := locales.Open("i18n/" + f.Name())
		if err != nil {
			panic(err)
		}
		defer source.Close()
		decoder := json.NewDecoder(source)
		translation := map[string]string{}
		if err := decoder.Decode(&translation); err != nil {
			log.Println("Error decoding JSON:", err)
		}
		translations[f.Name()] = translation
	}
}

// Translate translates a message using the current locale. If the locale is not
// supported, it will fall back to the default locale (en-US).
func (a *App) Translate(m string) string {
	locale, err := locale.GetLanguage()
	if err != nil {
		locale = "en-US"
	}
	result := a.T(m, locale, false)
	return result
}

// T translates a message using the given language. If "md" is true, it will
// compute the markdown to HTML.
func (a *App) T(m, lang string, md bool) string {
	parts := strings.Split(lang, "-")
	lang = parts[0] + ".json"
	if tr, ok := translations[lang]; ok {
		if t, ok := tr[m]; ok {
			if md {
				log.Println("Markdown:", t)
				t = string(MDtoHTML(t))
			}
			return t
		}
	}
	if lang != "en" {
		return a.T(m, "en", md)
	}
	return m
}
