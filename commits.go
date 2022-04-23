package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type commitRange struct {
	start string
	end string
}

type commitsModel struct {
	listModel
	commits []commit
}

func newCommitsModel() commitsModel {
	commits := gitLog()
	m := commitsModel{commits: commits}
	m.listModel.init(len(commits), true)
	return m
}

func (m commitsModel) name() string {
	return "commits"
}

func (m commitsModel) commit(index int) commit {
	return m.commits[index]
}

func (m commitsModel) selected() commit {
	return m.commits[m.cursor]
}

func (m commitsModel) renderCommit(index int) string {
	c := m.commit(index)

	ctime := time.Unix(c.Timestamp, 0)
	var age string
	years, months, days, hours, mins, secs, _ := Elapsed(ctime, time.Now())
	if years > 0 {
		age = fmt.Sprintf("%dY", years)
	} else if months > 0 {
		age = fmt.Sprintf("%dM", months)
	} else if days > 0 {
		age = fmt.Sprintf("%dD", days)
	} else if hours > 0 {
		age = fmt.Sprintf("%dh", hours)
	} else if mins > 0 {
		age = fmt.Sprintf("%dm", mins)
	} else {
		age = fmt.Sprintf("%ds", secs)
	}

	name := c.AuthorName
	if len(name) > 20 {
		parts := strings.Split(name, " ")
		if len(parts) >= 3 {
			name = fmt.Sprintf("%s ", parts[0])
			for i := 0; i < len(parts)-2; i++ {
				name += fmt.Sprintf("%c", parts[i][0])
			}
			name += fmt.Sprintf(" %s", parts[len(parts)-1])
		} else if len(parts) == 2 {
			name = fmt.Sprintf("%c %s", parts[0][0], parts[1])
		} else {
			name = name[0:20]
		}
	}

	marker := ""
	if index == m.marked {
		marker = "▶"
	}

	subjectStyle.Width(m.width -
		markerStyle.GetWidth() -
		hashStyle.GetWidth() -
		ageStyle.GetWidth() -
		nameStyle.GetWidth())

	if index == m.cursor {
		markerStyle.Background(cursorBg)
		hashStyle.Background(cursorBg)
		ageStyle.Background(cursorBg)
		nameStyle.Background(cursorBg)
		branchStyle.Background(cursorBg)
		tagStyle.Background(cursorBg)
		refStyle.Background(cursorBg)
		subjectStyle.Background(cursorBg)
	} else {
		markerStyle.UnsetBackground()
		hashStyle.UnsetBackground()
		ageStyle.UnsetBackground()
		nameStyle.UnsetBackground()
		branchStyle.UnsetBackground()
		tagStyle.UnsetBackground()
		refStyle.UnsetBackground()
		subjectStyle.UnsetBackground()
	}

	branches := ""
	tags := ""
	refs := ""

	if c.Decoration != "" {
		info := parseDecoration(c.Decoration)
		for _, b := range info.branches {
			branches += fmt.Sprintf("[%s] ", b)
		}
		if branches != "" {
			branches = branchStyle.Render(branches)
		}

		for _, t := range info.tags {
			tags += fmt.Sprintf("<%s> ", t)
		}
		if tags != "" {
			tags = tagStyle.Render(tags)
		}

		for _, r := range info.refs {
			refs += fmt.Sprintf("{%s} ", r)
		}
		if refs != "" {
			refs = refStyle.Render(refs)
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		markerStyle.Render(marker),
		hashStyle.Render(c.Commit[0:8]),
		ageStyle.Render(age),
		nameStyle.Render(name),
		branches,
		tags,
		refs,
		subjectStyle.Render(c.Subject),
	)
}

func (m commitsModel) render() string {
	var lines []string
	for i := m.first; i < m.last; i++ {
		lines = append(lines, m.renderCommit(i))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m commitsModel) getRange() (r commitRange) {
	r.start = m.commits[m.cursor].Commit
	r.end = ""

	if m.marked >= 0 {
		if m.marked > m.cursor {
			r.end = r.start
			r.start = m.commits[m.marked].Commit
		} else {
			r.end = m.commits[m.marked].Commit
		}
	}

	return
}

func (m commitsModel) getRangeStr() string {
	r := m.getRange()
	if r.end == "" {
		r.end = "<index>"
	}
	return fmt.Sprintf("%s..%s", trunc(r.start, 8), trunc(r.end, 8))
}