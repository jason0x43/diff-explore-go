package main

import "github.com/charmbracelet/lipgloss"

var cursorBg = lipgloss.Color("0")
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
var statusStyle = lipgloss.NewStyle().
	Inline(true).
	Background(lipgloss.Color("8")).
	Foreground(lipgloss.Color("15"))
var statStyle = lipgloss.NewStyle().Inline(true)
var diffNormalStyle = lipgloss.NewStyle().
	Inline(true)
var diffAddStyle = lipgloss.NewStyle().
	Inline(true).
	Foreground(lipgloss.Color("2"))
var diffDelStyle = lipgloss.NewStyle().
	Inline(true).
	Foreground(lipgloss.Color("1"))
var diffModStyle = lipgloss.NewStyle().
	Inline(true).
	Foreground(lipgloss.Color("18"))
var diffSepStyle = lipgloss.NewStyle().
	Inline(true).
	Foreground(lipgloss.Color("6"))
var popupStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("13")).
	Padding(1, 0).
	BorderTop(true).
	BorderLeft(true).
	BorderRight(true).
	BorderBottom(true)
