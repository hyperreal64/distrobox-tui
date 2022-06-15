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

type Model struct {
	tableModel table.Model
}

// getDataRows calls getDistroBoxItems() from distrobox.go
// and returns the info as a slice of table.Row
func getDataRows() (rows []table.Row) {
	distroBoxItems := getDistroboxItems()
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

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 15).WithStyle(idColumnStyle),
		table.NewColumn(columnKeyName, "Name", 30),
		table.NewColumn(columnKeyStatus, "Status", 30),
		table.NewColumn(columnKeyImage, "Image", 70).WithStyle(imageColumnStyle),
	}

	rows := getDataRows()

	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down")
	keys.RowUp.SetKeys("k", "up")

	model := Model{
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(headerStyle).
			SelectableRows(false).
			Focused(true).
			Border(customBorder).
			WithKeyMap(keys).
			WithPageSize(10).
			WithBaseStyle(baseStyle).
			HighlightStyle(highlightStyle).
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

	footerFormatStr := fmt.Sprintf(
		"[Pg. %d/%d] Current selection: ",
		m.tableModel.CurrentPage(),
		m.tableModel.MaxPages())

	var boxNameFmtStr string
	if highlightedRow.Data == nil {
		boxNameFmtStr = "No distroboxes available"
	} else {
		boxNameFmtStr = fmt.Sprintf("%s", highlightedRow.Data[columnKeyName])
	}

	footerText := lipgloss.JoinHorizontal(
		0, footerStyle.Render(footerFormatStr),
		footerBoxNameStyle.Render(boxNameFmtStr),
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
			cmds = append(cmds, clearScreen())

		case "enter":
			boxName := m.tableModel.HighlightedRow().Data[columnKeyName].(string)
			cmds = append(cmds, enterDistroBox(boxName))
			cmds = append(cmds, clearScreen())

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
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.MarginTop(1).MarginBottom(1).MarginLeft(1).Render("DistroBox TUI"),
		lipgloss.JoinHorizontal(
			0, subtleStyle.MarginRight(1).Render("[Left/right: change page]"),
			subtleStyle.MarginLeft(1).MarginRight(1).Render("[Up/down, j/k: move]"),
			subtleStyle.MarginLeft(1).MarginRight(1).Render("[Enter: enter distrobox]"),
			subtleStyle.MarginLeft(1).MarginRight(1).Render("[S: stop]"),
			subtleStyle.MarginLeft(1).MarginRight(1).Render("[X: remove]"),
			subtleStyle.MarginLeft(1).MarginRight(1).Render("[R: refresh]"),
			subtleStyle.MarginLeft(1).Render("[Q: quit]")),
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
