package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
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

var (
	catppuccinLatte = map[string]string{
		"idColumnStyle":      "#8839ef", // mauve
		"imageColumnStyle":   "#8839ef", // mauve
		"headerStyle":        "#40a02b", // green
		"footerBoxNameStyle": "#e64553", // maroon
		"baseStyleBorderFg":  "#181825", // mantle (mocha)
		"baseStyleFg":        "#4c4f69", // text
		"highlightStyleBg":   "#6c6f85", // subtext0
		"highlightStyleFg":   "#e6e9ef", // mantle
		"titleStyleFg":       "#1e66f5", // blue
		"subtleStyleFg":      "#fe640b", // peach
	}

	catppuccinMocha = map[string]string{
		"idColumnStyle":      "#cba6f7", // mauve
		"imageColumnStyle":   "#cba6f7", // mauve
		"headerStyle":        "#a6e3a1", // green
		"footerBoxNameStyle": "#eba0ac", // maroon
		"baseStyleBorderFg":  "#313244", // mantle
		"baseStyleFg":        "#cdd6f4", // text
		"highlightStyleBg":   "#313244", // subtext0
		"highlightStyleFg":   "#b4befe", // mantle
		"titleStyleFg":       "#89b4fa", // blue
		"subtleStyleFg":      "#fab387", // peach
	}
)

var (
	idColumnStyle = lipgloss.NewStyle().
			Faint(true).
			Foreground(lipgloss.AdaptiveColor{Light: catppuccinLatte["idColumnStyle"], Dark: catppuccinMocha["idColumnStyle"]})

	imageColumnStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: catppuccinLatte["imageColumnStyle"], Dark: catppuccinMocha["imageColumnStyle"]}).
				Faint(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: catppuccinLatte["headerStyle"], Dark: catppuccinMocha["headerStyle"]}).
			Bold(true)

	footerStyle = lipgloss.NewStyle()

	footerBoxNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: catppuccinLatte["footerBoxNameStyle"], Dark: catppuccinMocha["footerBoxNameStyle"]}).
				Bold(true)

	baseStyle = lipgloss.NewStyle().
			BorderForeground(lipgloss.AdaptiveColor{Light: catppuccinLatte["baseStyleBorderFg"], Dark: catppuccinMocha["baseStyleBorderFg"]}).
			Foreground(lipgloss.AdaptiveColor{Light: catppuccinLatte["baseStyleFg"], Dark: catppuccinMocha["baseStyleFg"]}).
			Align(lipgloss.Left)

	highlightStyle = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: catppuccinLatte["highlightStyleBg"], Dark: catppuccinMocha["highlightStyleBg"]}).
			Foreground(lipgloss.AdaptiveColor{Light: catppuccinLatte["highlightStyleFg"], Dark: catppuccinMocha["highlightStyleFg"]})

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: catppuccinLatte["titleStyleFg"], Dark: catppuccinMocha["titleStyleFg"]}).
			Bold(true)

	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: catppuccinLatte["subtleStyleFg"], Dark: catppuccinMocha["subtleStyleFg"]})
)
