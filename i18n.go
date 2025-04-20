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
		decoder.Decode(&translation)
		translations[f.Name()] = translation
	}
}

func (a *App) Translate(m string) string {
	locale, err := locale.GetLanguage()
	if err != nil {
		locale = "en-US"
	}
	return a.T(m, locale)
}

func (a *App) T(m, lang string) string {
	parts := strings.Split(lang, "-")
	lang = parts[0] + ".json"
	if tr, ok := translations[lang]; ok {
		if t, ok := tr[m]; ok {
			return t
		}
	}
	log.Println("Translation not found, falling back to en")
	return m
}
