package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type distroboxFinishedMsg struct{ err error }

type distroboxItem struct {
	id     string
	name   string
	status string
	image  string
}

type OCICmdOutput []struct {
	Id     string   `json:"Id"`
	Image  string   `json:"Image"`
	Mounts []string `json:"Mounts"`
	Names  []string `json:"Names"`
	Status string   `json:"Status"`
}

func clearScreen() tea.Cmd {
	cmd := exec.Command("clear")
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return distroboxFinishedMsg{err}
	})
}

func enterDistroBox(name string) tea.Cmd {
	cmd := exec.Command("distrobox", "enter", name)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return distroboxFinishedMsg{err}
	})
}

func removeDistroBox(name string) tea.Cmd {
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command("distrobox", "rm", name, "--force")
	cmd.Stdout = devnull
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return distroboxFinishedMsg{err}
	})
}

func stopDistroBox(name string) tea.Cmd {
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command("distrobox", "stop", name, "--yes")
	cmd.Stdout = devnull
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return distroboxFinishedMsg{err}
	})
}

func getOCICmd() string {
	podmanExists := true
	dockerExists := true

	podmanCmd, err := exec.LookPath("podman")
	if err != nil {
		podmanExists = false
	}

	dockerCmd, err := exec.LookPath("docker")
	if err != nil {
		dockerExists = false
	}

	if podmanExists {
		return podmanCmd
	} else if dockerExists {
		return dockerCmd
	} else {
		return ""
	}
}

func getDistroboxItems() (items []distroboxItem) {
	ociCmd := getOCICmd()
	if ociCmd == "" {
		log.Fatalln("Missing dependency: we need a container manager. Please install one of podman or docker.")
	}

	rawOutput, _ := exec.Command(ociCmd, "ps", "-a", "--format", "json").Output()
	var ociCmdOutput OCICmdOutput
	if err := json.Unmarshal(rawOutput, &ociCmdOutput); err != nil {
		log.Fatalln(err)
	}

	if len(ociCmdOutput) > 0 {
		for _, cmdOutputJsonElem := range ociCmdOutput {
			for _, mount := range cmdOutputJsonElem.Mounts {
				if mount == "/usr/bin/distrobox-host-exec" {
					box := distroboxItem{
						id:     cmdOutputJsonElem.Id[:12],
						name:   cmdOutputJsonElem.Names[0],
						status: cmdOutputJsonElem.Status,
						image:  cmdOutputJsonElem.Image,
					}

					items = append(items, box)
				}
			}
		}
	}

	return items
}
