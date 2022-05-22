package main

import (
	"os/exec"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type distroboxFinishedMsg struct{ err error }

func enterDistroBox(name string) tea.Cmd {
	cmd := exec.Command("distrobox", "enter", name)
	return tea.Exec(tea.WrapExecCommand(cmd), func(err error) tea.Msg {
		return distroboxFinishedMsg{err}
	})
}

func removeDistroBox(name string) tea.Cmd {
	cmd := exec.Command("distrobox", "rm", name, "--force")
	return tea.Exec(tea.WrapExecCommand(cmd), func(err error) tea.Msg {
		return distroboxFinishedMsg{err}
	})
}

func stopDistroBox(name string) tea.Cmd {
	cmd := exec.Command("distrobox", "stop", name)
	return tea.Exec(tea.WrapExecCommand(cmd), func(err error) tea.Msg {
		return distroboxFinishedMsg{err}
	})
}

type distroboxItem struct {
	id     string
	name   string
	status string
	image  string
}

func getDistroBoxItems() (items []distroboxItem) {
	output, _ := exec.Command("distrobox", "list").Output()

	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)

	outputString := re.ReplaceAllString(string(output), "")
	outputSlice := strings.Split(outputString, "\n")

	fieldSlice := [][]string{}
	for i := 1; i < len(outputSlice)-1; i++ {
		fieldSlice = append(fieldSlice, strings.Split(outputSlice[i], "|"))
	}

	for _, v := range fieldSlice {
		name := strings.TrimSpace(v[1])
		id := strings.TrimSpace(v[0])
		status := strings.TrimSpace(v[2])
		image := strings.TrimSpace(v[3])

		box := distroboxItem{id: id, name: name, status: status, image: image}
		items = append(items, box)
	}

	return items
}
