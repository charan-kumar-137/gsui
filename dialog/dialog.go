package dialog

import (
	"github.com/charan-kumar-137/gsui/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	text    string
	focused bool
	width   int
	height  int
}

func New() Model {
	return Model{text: "GS"}
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m Model) GetFocus() bool {
	return m.focused
}

func (m *Model) SetDimension(width, height int) {
	m.width = width
	m.height = height
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {
	case list.CurrentData:
		var data = msg.GetCurrentData()
		if data != nil {
			if data.IsBucket {
				var bucket = data.GetBucket(msg.GetCurrentCursor())
				if bucket != nil {
					m.text = bucket.DisplayString()
				} else {
					m.text = "Not Found Bucket " + msg.GetPath()
				}
			} else {
				var object = data.GetObject(msg.GetCurrentCursor())
				if object != nil {
					m.text = object.DisplayString()
				} else {
					m.text = "Not Found Object " + msg.GetPath()
				}
			}
		} else {
			m.text = "Not Found Data " + msg.GetPath()
		}
	}

	return m, nil
}

func (m Model) View() string {
	return lipgloss.NewStyle().Render(m.text)
}
