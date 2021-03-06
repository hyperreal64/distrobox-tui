# distrobox-tui

[![Go Report Card](https://goreportcard.com/badge/github.com/hyperreal64/distrobox-tui)](https://goreportcard.com/report/github.com/hyperreal64/distrobox-tui)

![screenshot.png](/screenshot.png)

A minimal TUI for [Distrobox](https://github.com/89luca89/distrobox) using [Bubbletea](https://github.com/charmbracelet/bubbletea).

Features [Catppuccin](https://github.com/catppuccin/catppuccin) color palette. Support for theme selection coming in future release.

My intention is to learn the Bubbletea framework by creating something (sort of?) useful.
## Install

### Requirements
* Go 1.18+
* Distrobox

```bash
go install github.com/hyperreal64/distrobox-tui@latest
```

## Usage

* Must be run from the host OS
* Ensure `$GOPATH/bin` is in your shell's $PATH

```bash
distrobox-tui
```

Currently it is not possible to *create* Distroboxes in the TUI, but this might be added in the future.
