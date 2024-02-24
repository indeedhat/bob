package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/indeedhat/bob/internal/templates"
)

type model interface {
	tea.Model
	Activate() tea.Cmd
}

type UiModel struct {
	stages       []model
	currentStage int
	stageChanged bool
}

func New(t templates.Templates) *UiModel {
	var (
		m = &UiModel{stageChanged: true}
		s = newState(t)
	)

	m.stages = []model{
		newComponentsModel(s, m.next),
		newDataModel(s, m.next),
		newBuildModel(s, m.next),
	}

	return m
}

// Init implements model.
func (m *UiModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, s := range m.stages {
		if cmd := s.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// Update implements model.
func (m *UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	_, cmd := m.stages[m.currentStage].Update(msg)
	cmds = append(cmds, cmd)

	if m.currentStage < 0 || m.currentStage >= len(m.stages) {
		cmds = append(cmds, tea.Quit)
	}

	if m.stageChanged {
		m.stageChanged = false
		if cmd := m.stages[m.currentStage].Activate(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// View implements model.
func (m *UiModel) View() string {
	if m.currentStage < 0 || m.currentStage >= len(m.stages) {
		return ""
	}

	return m.stages[m.currentStage].View()
}

func (u *UiModel) next() {
	u.currentStage++
	u.stageChanged = u.currentStage < len(u.stages)
}

var _ tea.Model = (*UiModel)(nil)
