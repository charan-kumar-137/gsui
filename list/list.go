package list

import (
	"strings"

	"github.com/charan-kumar-137/gsui/gcs"
	"github.com/charan-kumar-137/gsui/keys"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	tableUnUsedHeight = 7
)

type CurrentData struct {
	path   string
	data   *gcs.Data
	cursor int
}

func (cr CurrentData) GetCurrentData() *gcs.Data {
	return cr.data
}

func (cr CurrentData) GetCurrentCursor() int {
	return cr.cursor
}

func (cr CurrentData) GetPath() string {
	return cr.path
}

type Model struct {
	table       table.Model
	currentPath string
	data        *gcs.Data
	focused     bool
	width       int
	height      int
}

func getTableKeyMap() table.KeyMap {
	return table.KeyMap{
		LineUp:   keys.Keys.Up,
		LineDown: keys.Keys.Down,
		PageUp:   keys.Keys.PageUp,
		PageDown: keys.Keys.PageDown,
	}
}

func getTable(data *gcs.Data) table.Model {
	cols, rows := data.GetTableData()
	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithKeyMap(getTableKeyMap()),
		// table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#3367D6"))

	t.SetStyles(s)

	return t
}

func (m Model) GetCurrentPath() string {
	return m.currentPath
}

func (m *Model) UpdateCurrentPath(path string) {

	var data = gcs.GetData(path)
	if data != nil {

		m.table = getTable(data)
	} else {
		m.table = table.New()
	}
	m.currentPath = path
	m.data = data

}

func (m Model) GetSelectedRow() CurrentData {
	return CurrentData{data: m.GetData(), cursor: m.GetCursor(), path: m.GetCurrentPath()}
}

func (m Model) GetData() *gcs.Data {
	return m.data
}

func (m Model) GetCursor() int {
	return m.table.Cursor()
}

func (m Model) getSelectedName() string {

	if m.GetCursor() < 0 || m.data == nil {
		return ""
	}

	if m.data.IsBucket {
		var bucket = m.data.GetBucket(m.GetCursor())
		if bucket != nil {
			return bucket.Name
		}
	} else {
		var object = m.data.GetObject(m.GetCursor())
		if object != nil {
			return object.Name
		}
	}

	return ""
}

func New() Model {
	var data *gcs.Data = gcs.GetData("")

	return Model{table: getTable(data), data: data}
}

func (m *Model) Focus() {
	m.table.Focus()
	m.focused = true
}

func (m *Model) Blur() {
	m.table.Blur()
	m.focused = false
}

func (m Model) GetFocus() bool {
	return m.focused
}

func (m *Model) SetDimension(width, height int) {
	m.table.SetHeight(height - tableUnUsedHeight)
	m.table.SetWidth(width)
	m.width = width
	m.height = height
}

func (m Model) Init() tea.Cmd {

	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	if !m.GetFocus() {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Right):
			var path = m.currentPath
			if len(m.currentPath) == 0 {
				path = m.getSelectedName()
			} else {
				path = m.currentPath + "/" + m.getSelectedName()
			}
			m.UpdateCurrentPath(path)
			m.Focus()
		case key.Matches(msg, keys.Keys.Left):
			var path = m.currentPath
			if strings.LastIndex(m.currentPath, "/") == -1 {
				path = ""
			} else {
				path = m.currentPath[:strings.LastIndex(m.currentPath, "/")]
			}
			m.UpdateCurrentPath(path)
			m.Focus()
		}
	}

	if m.GetFocus() {
		m.table, cmd = m.table.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {

	return lipgloss.NewStyle().Render(m.table.View())
}
