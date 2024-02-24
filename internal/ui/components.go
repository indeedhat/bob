package ui

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type componentsModel struct {
	state *state
	form  *huh.Form
	next  func()
}

func newComponentsModel(s *state, next func()) *componentsModel {
	var (
		options []huh.Option[string]
		m       = &componentsModel{
			state: s,
			next:  next,
		}
	)

	for groupKey := range s.templates {
		options = append(options, huh.NewOption[string](groupKey, groupKey))
	}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key("components").
				Title("Choose your components").
				Options(options...).
				Value(&m.state.components).
				Validate(func(t []string) error {
					if len(t) == 0 {
						return errors.New("Select at least one group")
					}
					return nil
				}),
		),
	)

	return m
}

// Activate implements model.
func (*componentsModel) Activate() tea.Cmd {
	return nil
}

// Init implements model.
func (m *componentsModel) Init() tea.Cmd {
	return m.form.Init()
}

// Update implements model.
func (m *componentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (m *componentsModel) View() string {
	v := strings.TrimSuffix(m.form.View(), "\n\n")
	return v
}

var _ model = (*componentsModel)(nil)
