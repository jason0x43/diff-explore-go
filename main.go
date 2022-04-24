package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
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
	height         int
	width          int
	history        []string
	watcherReady   bool
	watcherLoading spinner.Model

	searching bool
	query     string

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
	if m.searching {
		return fmt.Sprintf("search: %s", m.query)
	} else {
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
	}

	return ""
}

func (m appModel) Init() tea.Cmd {
	return m.watcherLoading.Tick
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if v := m.currentView(); v != nil {
			// view height is 1 less than screen height
			v.setSize(msg.Width, msg.Height-1)
		}

	case tea.KeyMsg:
		if m.searching {
			switch msg.String() {
			case "esc":
				m.searching = false
			case "backspace":
				if len(m.query) > 0 {
					m.query = m.query[0 : len(m.query)-1]
				}
			case "enter":
				m.searching = false
				m.currentView().findNext(m.query)
			case "ctrl+c":
				return m, tea.Quit
			default:
				if msg.Type == tea.KeyRunes {
					m.query += msg.String()
				}
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case " ":
				if c := m.currentView(); c != nil {
					c.mark()
				}

			case "/":
				m.searching = true
				m.query = ""

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

			case "n":
				if len(m.query) > 0 {
					m.currentView().findNext(m.query)
				}

			case "N":
				if len(m.query) > 0 {
					m.currentView().findPrev(m.query)
				}

			case "enter":
				if m.currentViewName() == m.commits.name() {
					m.stats.setDiff(m.commits.getRange())
					m.stats.setSize(m.width, m.height-1)
					m.pushView("stats")
				} else if m.currentViewName() == m.stats.name() {
					if m.stats.cursor >= 0 {
						m.diff.setDiff(m.commits.getRange(), m.stats.selected())
						m.diff.setSize(m.width, m.height-1)
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

	default:
		var cmd tea.Cmd
		if !m.watcherReady {
			m.watcherLoading, cmd = m.watcherLoading.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m appModel) View() string {
	mainSection := ""

	if c := m.currentView(); c != nil {
		mainSection = c.render()
	}

	statusTwo := ""
	if !m.watcherReady {
		statusTwo += m.watcherLoading.View()
	}

	if m.diff.opts.ignoreWhitespace {
		statusTwo += "W"
	}

	statusTwoStyle.Width(len(statusTwo) + 2)

	statusThree := fmt.Sprintf(
		"%d/%d",
		m.currentView().getCursor(),
		m.currentView().getCount()-1,
	)

	statusThreeStyle.Width(len(statusThree) + 2)

	statusOneStyle.Width(
		m.width -
			statusTwoStyle.GetWidth() -
			statusThreeStyle.GetWidth(),
	)

	status := lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusOneStyle.Render(m.getStatus()),
		statusTwoStyle.Render(statusTwo),
		statusThreeStyle.Render(statusThree),
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

	s := spinner.New()
	s.Spinner = spinner.Dot

	m := appModel{
		history:        []string{"commits"},
		commits:        newCommitsModel(),
		stats:          newStatsModel(),
		diff:           newDiffModel(),
		status:         "",
		watcherLoading: s,
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