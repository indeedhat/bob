package templates

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

type VarType string

// IsValid checksi if the value of the VarType is of an known value
func (t VarType) IsValid() bool {
	return t == VarTypeString ||
		t == VarTypeArray ||
		t == VarTypeMap
}

const (
	VarTypeString VarType = "String"
	VarTypeArray  VarType = "Array"
	VarTypeMap    VarType = "Map"
)

type MetaData struct {
	Path struct {
		BasePath string `yaml:"base"`
		FileName string `yaml:"fileName"`
	} `yaml:"path"`
	Data map[string]VarType `yaml:"data"`
}

// SavePath builds up the final save path from its component parts then passes it through the
// template engine to have any markers filled
func (m MetaData) SavePath(vars map[string]any) (string, error) {
	t, err := template.New("meta-path").Parse(path.Join(m.Path.BasePath, m.Path.FileName))
	if err != nil {
		return "", fmt.Errorf("failed to parse save path template: %w", err)
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to build save path: %w", err)
	}

	return buf.String(), nil
}

// parseMetaData decodes the HCL template metadata into the MetaData struct
func parseMetaData(name string, data []byte) (*MetaData, error) {
	var m MetaData

	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("faild to parse metadata :%w", err)
	}

	return &m, nil
}

// extractMetaData loads the template file itself and locates then parses the metadata block
func extractMetaData(path string, t *template.Template) (*MetaData, error) {
	m := t.Lookup("meta")
	if m == nil {
		return nil, errors.New("meta block not found in template file")
	}

	tplName := filepath.Base(path)
	return parseMetaData(tplName, []byte(m.Tree.Root.String()))
}
