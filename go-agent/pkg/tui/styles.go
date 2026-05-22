package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	senderStyle = lipgloss.NewStyle().Foreground(special)
	botStyle    = lipgloss.NewStyle().Foreground(highlight)

	textStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"})

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
)
