package main

import "github.com/charmbracelet/lipgloss"

// colors
var cursorBg = lipgloss.Color("0")
var addFg = lipgloss.Color("2")
var remFg = lipgloss.Color("1")
var modFg = lipgloss.Color("18")

// styles
var markerStyle = lipgloss.NewStyle().Width(2)
var hashStyle = lipgloss.NewStyle().
	Width(9).
	PaddingRight(1).
	Foreground(lipgloss.Color("5"))
var ageStyle = lipgloss.NewStyle().
	Align(lipgloss.Right).
	Width(4).
	PaddingRight(1).
	Foreground(lipgloss.Color("4"))
var nameStyle = lipgloss.NewStyle().
	Width(21).
	PaddingRight(1).
	Foreground(lipgloss.Color("2"))
var branchStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("6"))
var tagStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("5"))
var refStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3"))
var subjectStyle = lipgloss.NewStyle().Inline(true)
var statusOneStyle = lipgloss.NewStyle().
	Inline(true).
	Background(lipgloss.Color("8")).
	Foreground(lipgloss.Color("15"))
var statusTwoStyle = lipgloss.NewStyle().
	Inline(true).
	Width(5).
	Align(lipgloss.Center).
	Background(lipgloss.Color("7")).
	Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"})
var statusThreeStyle = lipgloss.NewStyle().
	Inline(true).
	Width(3).
	Align(lipgloss.Center).
	Background(lipgloss.Color("12")).
	Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"})
var statStyle = lipgloss.NewStyle().Inline(true)
var statAddStyle = lipgloss.NewStyle().
	Align(lipgloss.Right).
	Width(2).
	Foreground(addFg)
var statDelStyle = lipgloss.NewStyle().
	Align(lipgloss.Right).
	Width(2).
	Foreground(remFg)
var statModStyle = lipgloss.NewStyle().
	Width(2).
	Foreground(modFg)
var diffNormalStyle = lipgloss.NewStyle().
	Inline(true)
var diffAddStyle = lipgloss.NewStyle().
	Inline(true).
	Foreground(addFg)
var diffRemStyle = lipgloss.NewStyle().
	Inline(true).
	Foreground(remFg)
var diffModStyle = lipgloss.NewStyle().
	Inline(true).
	Foreground(modFg)
var diffSepStyle = lipgloss.NewStyle().
	Inline(true).
	Foreground(lipgloss.Color("6"))