package templates

import (
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var tplFuncs = template.FuncMap{
	"ucfirst": ucfirst,
}

func ucfirst(s string) string {
	return cases.Title(language.English).String(s)
}
