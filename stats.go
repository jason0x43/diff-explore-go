package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type statsModel struct {
	listModel
	stats     []stat
	addsWidth int
	delsWidth int
	commits   commitRange
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

func (m statsModel) getCommitsStr() string {
	if m.commits.start == m.commits.end {
		return m.commits.start[:8]
	}
	if m.commits.end == "" {
		return fmt.Sprintf("%s..<index>", m.commits.start[:8])
	}
	return fmt.Sprintf("%s..%s", m.commits.start[:8], m.commits.end[:8])
}

func (m *statsModel) setDiff(c commitRange) {
	m.commits = c
	if c.start == c.end {
		m.stats = gitShow(c.start)
	} else {
		m.stats = gitDiffStat(c.start, c.end)
	}
	m.listModel.init(len(m.stats), false)
	m.addsWidth = 0
	m.delsWidth = 0
	for _, stat := range m.stats {
		addDigits := int(math.Floor(math.Log10(float64(stat.Adds)))) + 1
		delDigits := int(math.Floor(math.Log10(float64(stat.Dels)))) + 1
		m.addsWidth = max(m.addsWidth, addDigits)
		m.delsWidth = max(m.delsWidth, delDigits)
	}
}

func (m *statsModel) refresh() {
	m.stats = gitDiffStat(m.commits.start, m.commits.end)
	m.listModel.setCount(len(m.stats))
}

func (m statsModel) renderStat(index int) string {
	s := m.stats[index]
	parts := []string{}

	if index == m.cursor {
		statStyle.Background(cursorBg)
		statAddStyle.Background(cursorBg)
		statDelStyle.Background(cursorBg)
		statModStyle.Background(cursorBg)
	} else {
		statStyle.UnsetBackground()
		statAddStyle.UnsetBackground()
		statDelStyle.UnsetBackground()
		statModStyle.UnsetBackground()
	}

	statAddStyle.Width(m.addsWidth + 1).PaddingRight(1)
	parts = append(parts, statAddStyle.Render(fmt.Sprintf("%d", s.Adds)))
	statDelStyle.Width(m.delsWidth + 1).PaddingRight(1)
	parts = append(parts, statDelStyle.Render(fmt.Sprintf("%d", s.Dels)))

	statStyle.Width(m.width)

	path := s.Path
	if s.OldPath != "" {
		path = path + " ← " + s.OldPath
	}
	parts = append(parts, statStyle.Render(path))

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		parts...,
	)
}

func (m statsModel) render() string {
	var lines []string
	if m.end-m.start == 0 {
		for i := 0; i < m.height/4; i++ {
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