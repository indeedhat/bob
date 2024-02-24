package templates

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"

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
	Vars map[string]VarType `yaml:"vars"`
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
func extractMetaData(path string) (*MetaData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	t := template.New("")
	t.Tree.Mode = parse.ParseComments

	tplName := filepath.Base(path)

	tree, err := t.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	if tree.Tree.Root.Nodes[0].Type() != parse.NodeComment {
		return nil, errors.New("Failed to find metadata for template.\n It must be contained within a comment node at the very top of the template file")
	}

	hcl := strings.TrimPrefix(strings.TrimSuffix(tree.Tree.Root.Nodes[0].String(), "*/}}"), "{{/*")

	return parseMetaData(tplName, []byte(hcl))
}
