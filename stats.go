package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type statsModel struct {
	listModel
	stats   []stat
	commits commitRange
}

func newStatsModel() statsModel {
	m := statsModel{}
	m.listModel.init(0, true)
	return m
}

func (m statsModel) name() string {
	return "stats"
}

func (m statsModel) stat(index int) stat {
	return m.stats[index]
}

func (m statsModel) selected() stat {
	return m.stats[m.cursor]
}

func (m *statsModel) setDiff(c commitRange) {
	m.commits = c
	m.stats = gitDiffStat(c.start, c.end)
	m.listModel.init(len(m.stats), false)
}

func (m *statsModel) refresh() {
	m.stats = gitDiffStat(m.commits.start, m.commits.end)
	m.listModel.setCount(len(m.stats))
}

func (m statsModel) renderStat(index int) string {
	s := m.stats[index]

	statStyle.Width(m.width)

	statTypeStyle := statModStyle
	if s.Change[0] == 'A' {
		statTypeStyle = statAddStyle
	} else if s.Change[0] == 'D' {
		statTypeStyle = statRemStyle
	}

	if index == m.cursor {
		statStyle.Background(cursorBg)
		statTypeStyle.Background(cursorBg)
	} else {
		statStyle.UnsetBackground()
		statTypeStyle.UnsetBackground()
	}

	path := s.Path
	if s.OldPath != "" {
		path = path + " ← " + s.OldPath
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		statTypeStyle.Render(string(s.Change[0])),
		statStyle.Render(path),
	)
}

func (m statsModel) render() string {
	var lines []string
	if m.end - m.start == 0 {
		for i := 0; i < m.height / 4; i++ {
			lines = append(lines, "")
		}
		centerStyle := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(m.width)
		lines = append(lines, centerStyle.Render("No changes"))
		return lipgloss.JoinVertical(lipgloss.Center, lines...)
	} else {
		for i := m.start; i < m.end; i++ {
			lines = append(lines, m.renderStat(i))
		}
		return lipgloss.JoinVertical(lipgloss.Left, lines...)
	}
}

func (m *statsModel) findNext(query string) {
	q := strings.ToLower(query)
	for i := m.cursor + 1; i < m.count; i++ {
		c := strings.ToLower(m.renderStat(i))
		if strings.Contains(c, q) {
			m.cursor = i
			break
		}
	}
}

func (m *statsModel) findPrev(query string) {
	q := strings.ToLower(query)
	for i := m.cursor - 1; i >= 0; i++ {
		c := strings.ToLower(m.renderStat(i))
		if strings.Contains(c, q) {
			m.cursor = i
			break
		}
	}
}