package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type diffModel struct {
	listModel
	diff        []string
	commits		commitRange
	path        string
	oldPath     string
	opts        diffOptions
}

func newDiffModel() diffModel {
	m := diffModel{}
	m.listModel.init(0, false)
	return m
}

func (m diffModel) name() string {
	return "diff"
}

func (m *diffModel) setDiff(c commitRange, s stat) {
	m.commits = c
	m.setDiffStat(s)
}

func (m *diffModel) setDiffStat(s stat) {
	m.path = s.Path
	m.oldPath = s.OldPath
	m.refresh()
}

func (m *diffModel) refresh() {
	m.diff = gitDiff(m.commits.start, m.commits.end, m.path, m.oldPath, m.opts)
	m.listModel.setCount(len(m.diff))
}

func (m diffModel) renderDiffLine(index int) string {
	d := m.diff[index]
	d = strings.ReplaceAll(d, "\t", "    ")

	if len(d) > 0 {
		switch d[0] {
		case '-':
			return diffRemStyle.Render(d)
		case '+':
			return diffAddStyle.Render(d)
		case '@':
			return diffSepStyle.Render(d)
		}
	}

	return diffNormalStyle.Render(d)
}

func (m diffModel) render() string {
	var lines []string
	for i := m.first; i < m.last; i++ {
		lines = append(lines, m.renderDiffLine(i))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m *diffModel) findNext(query string) {
	q := strings.ToLower(query)
	for i := m.cursor + 1; i < m.count; i++ {
		c := strings.ToLower(m.renderDiffLine(i))
		if strings.Contains(c, q) {
			m.cursor = i
			break
		}
	}
}

func (m *diffModel) findPrev(query string) {
	q := strings.ToLower(query)
	for i := m.cursor - 1; i >= 0; i++ {
		c := strings.ToLower(m.renderDiffLine(i))
		if strings.Contains(c, q) {
			m.cursor = i
			break
		}
	}
}