package utils

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	err      error
	done     chan struct{}
	spinner  spinner.Model
	quitting bool
}

func initialModel(done chan struct{}) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s, done: done}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case <-m.done:
		m.quitting = true
		return m, tea.Quit
	default:
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s Sending %s to Kindle.. press q to quit\n\n", m.spinner.View(), filepath.Base(os.Args[1]))
	if m.quitting {
		return "\n"
	}
	return str
}

func WithSpinner(fn func() error) error {
	done := make(chan struct{})
	p := tea.NewProgram(initialModel(done))
	errs := make(chan error, 1)

	go func() {
		errs <- fn()
		close(done)
	}()

	if _, err := p.Run(); err != nil {
		return err
	}

	return <-errs
}
