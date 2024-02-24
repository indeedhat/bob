package ui

import (
	"strings"

	"github.com/indeedhat/bob/internal/templates"
)

type state struct {
	templates       templates.Templates
	activeTemplates templates.Templates
	components      []string
	data            map[string]*string
}

func (s *state) Data(shape map[string]templates.VarType) map[string]any {
	data := make(map[string]any)

	for key, varType := range shape {
		switch varType {
		case templates.VarTypeString:
			data[key] = *s.data[key]
		case templates.VarTypeArray:
			data[key] = osSafeSplit(*s.data[key])
		case templates.VarTypeMap:
			var (
				keys   = osSafeSplit(*s.data[key+".keys"])
				values = osSafeSplit(*s.data[key+".values"])
			)

			data[key] = make(map[string]string)
			for i, k := range keys {
				if i < len(values) {
					data[k] = values[i]
				} else {
					data[k] = ""
				}
			}
		}
	}

	return data

}

func newState(t templates.Templates) *state {
	return &state{
		templates: t,
		data:      make(map[string]*string),
	}
}

func osSafeSplit(s string) []string {
	return strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}
