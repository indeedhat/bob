package ui

import (
	"bytes"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	green = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	red   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	checkMark = green.SetString("✓")
	crossMark = red.SetString("✗")
)

type buildModel struct {
	state       *state
	next        func()
	tasksOutput bytes.Buffer
	hasError    bool
}

func newBuildModel(s *state, next func()) *buildModel {
	return &buildModel{
		state: s,
		next:  next,
	}
}

// Activate implements model.
func (m *buildModel) Activate() tea.Cmd {
	activeTemplates := m.state.templates.SelectGroups(m.state.components)
	shape, _ := activeTemplates.DataShape()
	log := activeTemplates.MakeFiles(m.state.Data(shape))

	for _, task := range *log.Tasks {
		mark := checkMark
		if task.Error != nil {
			m.hasError = true
			mark = crossMark
		}

		m.tasksOutput.WriteString(fmt.Sprintf("%s %s\n", mark, task.Task))
	}

	if m.hasError {
		log.Rollback()
	}

	return nil
}

// Init implements model.
func (*buildModel) Init() tea.Cmd {
	return nil
}

// Update implements model.
func (m *buildModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "esc", "q", "enter":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View implements model.
func (m *buildModel) View() string {
	message := m.tasksOutput.String()

	if m.hasError {
		message += "\n\n" + red.SetString("Changes reverted").String()
	}

	return message
}

var _ model = (*buildModel)(nil)
