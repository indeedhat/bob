package templates

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/indeedhat/bob/internal/transaction"
)

type TemplateGroup []Template

// DataShape returns a map of all the unique vars assigned to the templates in the group
// if any vars conflict then an error will be returned
func (g TemplateGroup) DataShape() (map[string]VarType, error) {
	vars := make(map[string]VarType)

	for _, t := range g {
		group, groupKey := t.Keys()

		for k, v := range t.Meta.Data {
			if !v.IsValid() {
				return nil, fmt.Errorf("%s:%s has invalid var type %s: %s", group, groupKey, k, v)
			}

			varKey := ucfirst(k)
			if val, ok := vars[varKey]; ok && val != v {
				return nil, fmt.Errorf("%s:%s has type conflict with existing var %s: %s -> %s",
					group,
					groupKey,
					k,
					v,
					val,
				)
			}

			vars[varKey] = v
		}
	}

	return vars, nil
}

type Templates map[string]TemplateGroup

func (t Templates) SelectGroups(keys []string) Templates {
	groups := make(Templates)

	for _, key := range keys {
		if group, ok := t[key]; ok {
			groups[key] = group
		}
	}

	return groups
}

func (t Templates) DataShape() (map[string]VarType, error) {
	vars := make(map[string]VarType)

	for groupKey, group := range t {
		shape, err := group.DataShape()
		if err != nil {
			return nil, err
		}

		for k, v := range shape {
			if val, ok := vars[k]; ok && val != v {
				return nil, fmt.Errorf("%s has type conflict with existing var %s: %s -> %s",
					groupKey,
					k,
					v,
					val,
				)
			}

			vars[k] = v
		}
	}

	return vars, nil
}

func (t Templates) MakeFiles(data map[string]any) *transaction.FileMaker {
	maker := transaction.NewFileMaker(func(log *transaction.TaskLog[string]) {
		for _, group := range t {
			for _, template := range group {
				savePath, err := template.Meta.SavePath(data)
				if err != nil {
					log.Add(path.Join(template.Meta.Path.BasePath, template.Meta.Path.FileName))
					log.Failed(err)
					continue
				}

				log.Do(savePath, func() error {
					fh, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						return err
					}
					defer fh.Close()

					return template.Template.Execute(fh, data)
				})
			}
		}
	})

	maker.Run()

	return maker
}

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

		if t.Template, err = template.New(filename).Funcs(tplFuncs).Parse(string(data)); err != nil {
			return err
		}

		if t.Meta, err = extractMetaData(filename, t.Template); err != nil {
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
