package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// Column field labels
const (
	columnKeyID     = "id"
	columnKeyName   = "name"
	columnKeyStatus = "status"
	columnKeyImage  = "image"
)

// Border definition
var (
	customBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		InnerJunction:  "┼",

		InnerDivider: "│",
	}
)

type Model struct {
	tableModel table.Model
}

// getDataRows calls getDistroBoxItems() from distrobox.go
// and returns the info as a slice of table.Row
func getDataRows() (rows []table.Row) {
	distroBoxItems := getDistroBoxItems()
	for _, v := range distroBoxItems {
		row := table.NewRow(table.RowData{
			columnKeyID:     v.id,
			columnKeyName:   v.name,
			columnKeyStatus: v.status,
			columnKeyImage:  v.image,
		})
		rows = append(rows, row)
	}

	return rows
}

// Create a new Model object
func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 15).WithStyle(
			lipgloss.NewStyle().
				Faint(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#7b59c0", Dark: "#ABE9B3"}),
		),
		table.NewColumn(columnKeyName, "Name", 30),
		table.NewColumn(columnKeyStatus, "Status", 30),
		table.NewColumn(columnKeyImage, "Image", 70).WithStyle(
			lipgloss.NewStyle().
				Faint(true).
				Foreground(lipgloss.AdaptiveColor{Light: "#ca402b", Dark: "#ABE9B3"}),
		),
	}

	rows := getDataRows()

	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down")
	keys.RowUp.SetKeys("k", "up")

	model := Model{
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#379a37", Dark: "#DDB6F2"}).Bold(true)).
			SelectableRows(false).
			Focused(true).
			Border(customBorder).
			WithKeyMap(keys).
			WithStaticFooter("Footer!").
			WithPageSize(5).
			WithBaseStyle(
				lipgloss.NewStyle().
					BorderForeground(lipgloss.AdaptiveColor{Light: "#1b181b", Dark: "#575268"}).
					Foreground(lipgloss.AdaptiveColor{Light: "#695d69", Dark: "#D9E0EE"}).
					Align(lipgloss.Left),
			).
			HighlightStyle(
				lipgloss.NewStyle().
					Background(lipgloss.AdaptiveColor{Light: "#ab9bab", Dark: "#575268"}).
					Foreground(lipgloss.AdaptiveColor{Light: "#1b181b", Dark: "#C9CBFF"}),
			).
			SortByAsc(columnKeyName),
	}

	model.updateFooter()

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

// updateFooter() updates the View's footer when a row is selected
// by the user. The columnKeyName of the data in that row will be
// displayed in the footer.
func (m *Model) updateFooter() {
	highlightedRow := m.tableModel.HighlightedRow()

	footerText := fmt.Sprintf(
		"Pg. %d/%d - Currently looking at: %s",
		m.tableModel.CurrentPage(),
		m.tableModel.MaxPages(),
		highlightedRow.Data[columnKeyName],
	)

	m.tableModel = m.tableModel.WithStaticFooter(footerText)
}

// Update() listens and processes input from the user and executes
// actions according to what keybindings the user inputs
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	// We control the footer text, so make sure to update it
	m.updateFooter()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "Q":
			cmds = append(cmds, tea.Quit)

		case "enter":
			boxName := m.tableModel.HighlightedRow().Data[columnKeyName].(string)
			cmds = append(cmds, enterDistroBox(boxName))

		case "S":
			boxName := m.tableModel.HighlightedRow().Data[columnKeyName].(string)
			cmds = append(cmds, stopDistroBox(boxName))

		case "X", "delete":
			boxName := m.tableModel.HighlightedRow().Data[columnKeyName].(string)
			cmds = append(cmds, removeDistroBox(boxName))

		case "R":
			m = NewModel()
		}
	}

	return m, tea.Batch(cmds...)
}

// View() defines what is displayed in the user interface
func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#516aec", Dark: "#96CDFB"}).
		Bold(true)

	bodyStyle := lipgloss.NewStyle()

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.MarginTop(1).MarginBottom(1).Render("DistroBox TUI"),
		lipgloss.JoinHorizontal(0, bodyStyle.Width(12).Render("Left/Right"), bodyStyle.MarginLeft(5).Render("change page")),
		lipgloss.JoinHorizontal(0, bodyStyle.Width(12).Render("Enter"), bodyStyle.MarginLeft(5).Render("enter selected distrobox")),
		lipgloss.JoinHorizontal(0, bodyStyle.Width(12).Render("S"), bodyStyle.MarginLeft(5).Render("stop selected distrobox")),
		lipgloss.JoinHorizontal(0, bodyStyle.Width(12).Render("X"), bodyStyle.MarginLeft(5).Render("remove selected distrobox")),
		lipgloss.JoinHorizontal(0, bodyStyle.Width(12).Render("R"), bodyStyle.MarginLeft(5).Render("refresh view")),
		lipgloss.JoinHorizontal(0, bodyStyle.MarginBottom(1).Width(12).Render("Q"), bodyStyle.MarginLeft(5).Render("quit")),
		m.tableModel.View(),
	)

	return lipgloss.NewStyle().MarginLeft(1).Render(view)
}

func main() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
