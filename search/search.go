package search

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	text    textinput.Model
	focused bool
	width   int
	height  int
}

func New() Model {
	ti := textinput.New()

	ti.Prompt = "gs://"

	ti.Placeholder = "<bucket>/<object>"
	ti.PlaceholderStyle = lipgloss.NewStyle()

	return Model{text: ti, focused: false}
}

func (m *Model) Focus() {
	m.text.Focus()
	m.focused = true
}

func (m *Model) Blur() {
	m.text.Blur()
	m.focused = false
}

func (m Model) GetFocus() bool {
	return m.focused
}

func (m *Model) SetDimension(width, height int) {
	m.text.Width = width
	m.width = width
	m.height = height
}

func (m Model) GetCurrentPath() string {
	return m.text.Value()
}

func (m *Model) UpdateCurrentPath(path string) {
	m.text.SetValue(path)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.GetFocus() {
		m.text, cmd = m.text.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.NewStyle().Render(m.text.View())
}
