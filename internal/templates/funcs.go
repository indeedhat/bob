package templates

import (
	"fmt"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var builtinFuncsAlias = template.FuncMap{
	"and":      func() {},
	"call":     func() {},
	"html":     func() {},
	"index":    func() {},
	"slice":    func() {},
	"js":       func() {},
	"len":      func() {},
	"not":      func() {},
	"or":       func() {},
	"print":    fmt.Sprint,
	"printf":   fmt.Sprintf,
	"println":  fmt.Sprintln,
	"urlquery": func() {},

	// Comparisons
	"eq": func() {}, // ==
	"ge": func() {}, // >=
	"gt": func() {}, // >
	"le": func() {}, // <=
	"lt": func() {}, // <
	"ne": func() {}, // !=
}

var tplFuncs = template.FuncMap{
	"ucfirst": func(s string) string {
		return cases.Title(language.English).String(s)
	},
}
