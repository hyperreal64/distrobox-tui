package main

import (
	"encoding/json"
	"fmt"
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
	Id     string `json:"Id"`
	Image  string `json:"Image"`
	Labels struct {
		Manager string `json:"manager"`
	} `json:"Labels"`
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
	dbCmd := fmt.Sprintf("clear && distrobox enter %s", name)
	cmd := exec.Command("/bin/sh", "-c", dbCmd)
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

		for _, jsonElem := range ociCmdOutput {
			for _, mount := range jsonElem.Mounts {
				if mount == "/usr/bin/distrobox-export" {
					box := distroboxItem{
						id:     jsonElem.Id[:12],
						name:   jsonElem.Names[0],
						status: jsonElem.Status,
						image:  jsonElem.Image,
					}

					items = append(items, box)
				}
			}
		}

		if len(items) == 0 {
			for _, jsonElem := range ociCmdOutput {
				if jsonElem.Labels.Manager == "distrobox" {
					box := distroboxItem{
						id:     jsonElem.Id[:12],
						name:   jsonElem.Names[0],
						status: jsonElem.Status,
						image:  jsonElem.Image,
					}

					items = append(items, box)
				}
			}
		}
	}

	return items
}
