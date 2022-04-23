package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type watcherMessage struct {
	event string
	path  string
}

type chord struct {
	key       string
	startTime time.Time
}

func (c chord) isExpired() bool {
	now := time.Now()
	return now.Sub(c.startTime).Milliseconds() > 1000
}

func (c *chord) start(key string) {
	c.startTime = time.Now()
	c.key = key
}

func (c chord) getKey() string {
	if c.isExpired() {
		return ""
	}
	return c.key
}

type appModel struct {
	height       int
	width        int
	history      []string
	watcherReady bool

	chord chord

	commits commitsModel
	stats   statsModel
	diff    diffModel

	status string
}

func trunc(s string, size int) string {
	if len(s) > size {
		return s[:size]
	}
	return s
}

func (m *appModel) currentView() listView {
	v := m.history[len(m.history)-1]
	switch v {
	case m.commits.name():
		return &m.commits
	case m.stats.name():
		return &m.stats
	case m.diff.name():
		return &m.diff
	}
	return nil
}

func (m appModel) currentViewName() string {
	if c := m.currentView(); c != nil {
		return c.name()
	}
	return ""
}

func (m *appModel) pushView(view string) {
	m.history = append(m.history, view)
}

func (m *appModel) popView() {
	m.history = m.history[:len(m.history)-1]
}

func (m appModel) getStatus() string {
	switch m.currentView().name() {
	case m.commits.name():
		return m.commits.getRangeStr()
	case m.stats.name():
		return m.commits.getRangeStr()
	case m.diff.name():
		r := m.commits.getRangeStr()
		path := m.stats.stat(m.stats.cursor).Path
		return fmt.Sprintf("%s: %s", r, path)
	}

	return ""
}

func (m appModel) Init() tea.Cmd {
	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if v := m.currentView(); v != nil {
			v.setSize(msg.Width, msg.Height)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case " ":
			if c := m.currentView(); c != nil {
				c.mark()
			}

		case "1":
			m.chord.start("1")

		case "G":
			if c := m.currentView(); c != nil {
				if m.chord.getKey() == "1" {
					c.scrollToTop()
				} else {
					c.scrollToBottom()
				}
			}

		case "ctrl+f":
			if c := m.currentView(); c != nil {
				c.nextPage()
			}

		case "ctrl+u":
			if c := m.currentView(); c != nil {
				c.prevPage()
			}

		case "j", "down":
			if c := m.currentView(); c != nil {
				c.nextItem()
			}

		case "J":
			if m.currentView().name() == m.diff.name() {
				m.stats.nextItem()
				m.diff.setDiffStat(m.stats.selected())
			}

		case "k", "up":
			if c := m.currentView(); c != nil {
				c.prevItem()
			}

		case "K":
			if m.currentView().name() == m.diff.name() {
				m.stats.prevItem()
				m.diff.setDiffStat(m.stats.selected())
			}

		case "enter":
			if m.currentViewName() == m.commits.name() {
				m.stats.setDiff(m.commits.getRange())
				m.stats.setSize(m.width, m.height)
				m.pushView("stats")
			} else if m.currentViewName() == m.stats.name() {
				if m.stats.cursor >= 0 {
					m.diff.setDiff(m.commits.getRange(), m.stats.selected())
					m.diff.setSize(m.width, m.height)
					m.pushView("diff")
				}
			}

		case "esc", "q":
			if len(m.history) == 1 {
				return m, tea.Quit
			}
			m.popView()

		case "w":
			if c := m.currentView(); c != nil && c.name() == "diff" {
				m.diff.opts.ignoreWhitespace = !m.diff.opts.ignoreWhitespace
				m.diff.refresh()
			}
		}

	case watcherMessage:
		switch msg.event {
		case "ready":
			m.watcherReady = true
		case "filechange":
			if m.currentViewName() == m.diff.name() && m.diff.path == msg.path {
				m.diff.refresh()
			}
			m.stats.refresh()
		}
	}

	m.status = m.getStatus()

	return m, nil
}

func (m appModel) View() string {
	mainSection := ""

	if c := m.currentView(); c != nil {
		mainSection = c.render()
	}

	statusRightStyle.Width(5)
	statusLeftStyle.Width(m.width - statusRightStyle.GetWidth())

	statusRight := ""
	if !m.watcherReady {
		// TODO use a spinner for this
		statusRight += "-"
	}
	if !m.chord.isExpired() && m.chord.key != "" {
		statusRight += m.chord.key
	}

	if m.diff.opts.ignoreWhitespace {
		statusRight += "W"
	}

	status := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusLeftStyle.Render(m.status),
		statusRightStyle.Render(statusRight),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.PlaceVertical(m.height-1, lipgloss.Top, mainSection),
		status,
	)
}

func main() {
	repoPath := "."
	if len(os.Args) > 1 {
		repoPath = os.Args[1]
	}
	os.Chdir(repoPath)

	m := appModel{
		history: []string{"commits"},
		commits: newCommitsModel(),
		stats:   newStatsModel(),
		diff:    newDiffModel(),
		status:  "",
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	onNotify := func(event, path string) {
		if event == "ready" {
			p.Send(watcherMessage{event: "ready", path: ""})
		} else {
			p.Send(watcherMessage{event: "filechange", path: path})
		}
	}
	watcher := watchRepo(".", onNotify)
	defer watcher.Close()

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}