package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TemplateGroup []Template

// Vars returns a map of all the unique vars assigned to the templates in the group
// if any vars conflict then an error will be returned
func (g TemplateGroup) Vars() (map[string]VarType, error) {
	vars := make(map[string]VarType)

	for _, t := range g {
		for k, v := range t.Meta.Vars {
			group, key := t.Keys()

			if !v.IsValid() {
				return nil, fmt.Errorf("%s:%s has invalid var type %s: %s", group, key, k, v)
			}

			if val, ok := vars[k]; ok && val != v {
				return nil, fmt.Errorf("%s:%s has type conflict with existing var %s: %s -> %s", group, key, k, v, val)
			}

			vars[k] = v
		}
	}

	return vars, nil
}

type Templates map[string]TemplateGroup

type Template struct {
	Meta     *MetaData
	Path     string
	Template *template.Template
}

// Keys returns the group and template key for the current template based on its filepath
func (t Template) Keys() (string, string) {
	filename := strings.TrimSuffix(filepath.Base(t.Path), filepath.Ext(t.Path))

	if strings.Contains(filename, ":") {
		parts := strings.SplitN(filename, ":", 2)
		return parts[0], parts[1]
	}

	return filename, filename
}

// Load loads the template files from the provided path and builds up the templates map
// along with preloading the template metadata
func Load(path string) (Templates, error) {
	tpls := make(Templates)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".tpl" {
			return nil
		}

		var (
			key      string
			data     []byte
			filename = filepath.Base(path)
			t        = Template{Path: path}
		)

		if data, err = os.ReadFile(path); err != nil {
			return err
		}

		if t.Meta, err = extractMetaData(path); err != nil {
			return err
		}

		if t.Template, err = template.New(filename).Funcs(tplFuncs).Parse(string(data)); err != nil {
			return err
		}

		if strings.Contains(filename, ":") {
			key = strings.Split(filename, ":")[0]
		} else {
			key = strings.TrimSuffix(filename, filepath.Ext(filename))
		}

		tpls[key] = append(tpls[key], t)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load template files: %w", err)
	}

	return tpls, nil
}
