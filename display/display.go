package display

import (
	// "fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	dialog "github.com/charan-kumar-137/gsui/dialog"
	"github.com/charan-kumar-137/gsui/keys"

	"github.com/charan-kumar-137/gsui/list"
	"github.com/charan-kumar-137/gsui/search"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Enum for Selecting Active Display
type ActiveDisplay int

const (
	// Search View
	SEARCH ActiveDisplay = iota
	// List View
	LIST ActiveDisplay = iota
	// Dialog View
	DIALOG ActiveDisplay = iota
	// None
	NONE ActiveDisplay = iota
)

var (
	// Total No. of Display
	TotalDisplay int = 4

	// To Avoid Overflow of Content
	unUsedWidth  int = 2
	unUsedHeight int = 2

	// Default Border Style
	borderStyle lipgloss.Style = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#3367D6"))

	// Active Border Style
	activeBorderStyle lipgloss.Style = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#3367D6"))
)

// Display Model
type Model struct {
	// Search View
	searchView search.Model
	// List View
	listView list.Model
	// Dialog View
	dialogView dialog.Model
	// Active Display
	active ActiveDisplay
}

// Toggle toggleFocus between displays
func (m *Model) toggleFocus() {

	m.active = (m.active + 1) % ActiveDisplay(TotalDisplay)

	switch m.active {
	case SEARCH:
		m.searchView.Focus()
		m.listView.Blur()
		m.dialogView.Blur()
	case LIST:
		m.listView.Focus()
		m.searchView.Blur()
		m.dialogView.Blur()
	case DIALOG, NONE:
		m.blur()
	}
}

func (m *Model) blur() {
	m.active = NONE
	m.searchView.Blur()
	m.listView.Blur()
	m.dialogView.Blur()
}

// Get Border
func (m Model) getBorder(display ActiveDisplay) lipgloss.Style {
	if m.active == display {
		return activeBorderStyle
	}
	return borderStyle
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Escape):
			m.blur()
		case key.Matches(msg, keys.Keys.Tab):
			m.toggleFocus()
		case key.Matches(msg, keys.Keys.Quit):
			if m.active == NONE {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.searchView.SetDimension(msg.Width, msg.Height)
		m.listView.SetDimension(msg.Width, msg.Height)
		m.dialogView.SetDimension(msg.Width, msg.Height)
	}

	// Update Search View Text
	searchViewUpdate, searchViewUpdateCmd := m.searchView.Update(msg)
	m.searchView = searchViewUpdate
	cmds = append(cmds, searchViewUpdateCmd)

	// Update Actions of List View
	listViewUpdateAction, listViewUpdateActionCmd := m.listView.Update(msg)
	m.listView = listViewUpdateAction
	cmds = append(cmds, listViewUpdateActionCmd)

	// Update based on active display
	switch m.active {
	case SEARCH:
		m.listView.UpdateCurrentPath(m.searchView.GetCurrentPath())
	case LIST:
		m.searchView.UpdateCurrentPath(m.listView.GetCurrentPath())
	}

	// Update Dialog View
	dialogViewUpdate, dialogViewUpdateCmd := m.dialogView.Update(m.listView.GetSelectedRow())
	m.dialogView = dialogViewUpdate
	cmds = append(cmds, dialogViewUpdateCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	// Get Actual Width From Terminal
	actualWidth, actualHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	var usableWidth = actualWidth - unUsedWidth
	var usableHeight = actualHeight - unUsedHeight

	// Search View
	var searchView = m.getBorder(SEARCH).Width(usableWidth).Render(m.searchView.View())
	var searchViewHeight = lipgloss.Height(searchView)

	// List View
	var listWidth = int(0.7 * float64(usableWidth))
	var listHeight = usableHeight - (searchViewHeight)
	var listView = m.getBorder(LIST).Width(listWidth).Height(listHeight).Render(m.listView.View())

	// Dialog View
	var dialogWidth = usableWidth - lipgloss.Width(listView)
	var dialogHeight = usableHeight - (searchViewHeight)
	var dialogView = m.getBorder(DIALOG).Width(dialogWidth).Height(dialogHeight).Render(m.dialogView.View())

	// Final View to be Displayed
	var view = lipgloss.JoinVertical(lipgloss.Top,
		searchView,
		lipgloss.JoinHorizontal(lipgloss.Top, listView, dialogView),
	)

	return view
}

func Run() {

	var searchView = search.New()
	var listView = list.New()
	var dialogView = dialog.New()

	var m = Model{
		searchView: searchView,
		listView:   listView,
		dialogView: dialogView,
		active:     NONE,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal("Failed to start gsui", err)
		os.Exit(1)
	}
}
