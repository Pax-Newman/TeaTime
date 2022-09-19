package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Messages
type TickMsg time.Time
type StartStopMsg struct{ running bool }
type ResetMsg struct{}

// Commands
func clockTick(interval time.Duration) tea.Cmd {
	return tea.Every(interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Model
type model struct {
	start   time.Time
	current time.Time

	duration time.Duration

	running  bool
	interval time.Duration
}

// Model Methods
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.Reset(),
		m.Start(),
	)
}

func (m model) Start() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return StartStopMsg{running: true} },
		clockTick(m.interval),
	)
}

func (m model) Stop() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return StartStopMsg{running: false} },
		clockTick(m.interval),
	)
}

func (m model) Reset() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return ResetMsg{} },
		clockTick(m.interval),
	)
}

func (m model) Toggle() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return StartStopMsg{!m.running} },
		clockTick(m.interval),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if m.running {
			m.current = time.Time(msg)
			m.duration += m.interval
			return m, clockTick(m.interval)
		}
	case StartStopMsg:
		// keep timer consistent after pauses
		if msg.running && !m.running && !m.start.IsZero() {
			m.start = m.start.Add(time.Since(m.current))
			m.current = time.Now()
		}
		m.running = msg.running
		return m, nil
	case ResetMsg:
		m.start = time.Now()
		m.current = time.Now()
		return m, nil

	// -------------------------- TEMP DEBUGGING CASES --------------------------
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "r":
			return m, m.Reset()
		case "s":
			return m, m.Toggle()
		}
	}
	return m, nil
}

func (m model) View() string {
	elapsed := m.current.Sub(m.start).Round(m.interval).String()

	return elapsed
}

func main() {
	tea.NewProgram(model{
		interval: time.Millisecond,
	}).Start()
}
