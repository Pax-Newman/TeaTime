// Provides a simple system clock based timer
package main

import (
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Globals
var (
	mtx    sync.Mutex
	lastId int
)

// Functions

func NewWithInterval(interval time.Duration) Model {
	mtx.Lock()
	defer mtx.Unlock()

	lastId++

	return Model{
		interval: interval,
		id:       lastId,
	}
}

func New() Model {
	return NewWithInterval(time.Second)
}

// Messages
type TickMsg struct {
	clickTime time.Time
	id        int
}
type StartStopMsg struct {
	running bool
	id      int
}
type ResetMsg struct{ id int }

// Commands
func clockTick(interval time.Duration, id int) tea.Cmd {
	return tea.Every(interval, func(t time.Time) tea.Msg {
		return TickMsg{t, id}
	})
}

// Model
type Model struct {
	start   time.Time
	current time.Time

	id int

	running  bool
	interval time.Duration
}

// Model Methods
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.Reset(),
		m.Start(),
	)
}

func (m Model) Start() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return StartStopMsg{true, m.id} },
		clockTick(m.interval, m.id),
	)
}

func (m Model) Stop() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return StartStopMsg{false, m.id} },
		clockTick(m.interval, m.id),
	)
}

func (m Model) Reset() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return ResetMsg{m.id} },
		clockTick(m.interval, m.id),
	)
}

func (m Model) Toggle() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return StartStopMsg{!m.running, m.id} },
		clockTick(m.interval, m.id),
	)
}

func (m Model) Running() bool {
	return m.running
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if msg.id == m.id && m.running {
			m.current = time.Time(msg.clickTime)
			return m, clockTick(m.interval, m.id)
		}
	case StartStopMsg:
		if msg.id == m.id {
			// keep timer consistent after pauses
			if msg.running && !m.running && !m.start.IsZero() {
				m.start = m.start.Add(time.Since(m.current))
				m.current = time.Now()
			}
			m.running = msg.running
		}
	case ResetMsg:
		if msg.id == m.id {
			m.start = time.Now()
			m.current = time.Now()
		}

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

func (m Model) View() string {
	elapsed := m.current.Sub(m.start).Round(m.interval).String()
	return elapsed
}

func main() {
	tea.NewProgram(NewWithInterval(time.Millisecond)).Start()
}
