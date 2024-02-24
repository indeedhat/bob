package ui

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/indeedhat/bob/internal/templates"
)

type dataModel struct {
	state *state
	form  *huh.Form
	next  func()

	// testing
	group *huh.Group
}

func newDataModel(s *state, next func()) *dataModel {
	var (
		options []huh.Option[string]
		m       = &dataModel{
			state: s,
			next:  next,
		}
	)

	for groupKey := range s.templates {
		options = append(options, huh.NewOption[string](groupKey, groupKey))
	}

	return m
}

// Activate implements model.
func (m *dataModel) Activate() tea.Cmd {
	shape, err := m.state.templates.SelectGroups(m.state.components).DataShape()
	if err != nil {
		tea.Println(err.Error())
		return tea.Quit
	}

	var (
		fields   []huh.Field
		validate = func(s string) error {
			if s == "" {
				return errors.New("Required")
			}
			return nil
		}
	)

	for key, varType := range shape {
		switch varType {
		case templates.VarTypeString:
			var tmp string
			m.state.data[key] = &tmp
			fields = append(
				fields,
				huh.NewInput().Title(key).Value(m.state.data[key]).Validate(validate),
			)
		case templates.VarTypeArray:
			var tmp string
			m.state.data[key] = &tmp
			fields = append(
				fields,
				huh.NewText().Title(key).Value(m.state.data[key]).Validate(validate),
			)
		case templates.VarTypeMap:
			var tmp, tmp1 string
			m.state.data[key+".keys"] = &tmp
			m.state.data[key+".values"] = &tmp1
			fields = append(
				fields,
				huh.NewText().Title(key+" (Keys)").Value(m.state.data[key+".keys"]).Validate(validate),
				huh.NewText().Title(key+" (Values)").Value(m.state.data[key+".values"]).Validate(validate),
			)
		}
	}

	m.form = huh.NewForm(
		huh.NewGroup(fields...),
	)

	return m.form.Init()
}

// Init implements model.
func (m *dataModel) Init() tea.Cmd {
	return nil
}

// Update implements model.
func (m *dataModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		m.next()
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// View implements model.
func (m *dataModel) View() string {
	if m.form == nil {
		return ""
	}
	v := strings.TrimSuffix(m.form.View(), "\n\n")
	return v
}

var _ model = (*dataModel)(nil)
