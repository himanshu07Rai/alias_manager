package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type createModel struct {
	nameInput  textinput.Model
	cmdInput   textinput.Model
	focusIndex int
	err        error
	submitted  bool
}

func initialCreateModel() createModel {
	ni := textinput.New()
	ni.Placeholder = "e.g. gs"
	ni.Prompt = "Alias name: "
	ni.CharLimit = 50
	ni.Width = 40

	ci := textinput.New()
	ci.Placeholder = "e.g. git status"
	ci.Prompt = "Command:    "
	ci.CharLimit = 200
	ci.Width = 60

	return createModel{
		nameInput: ni,
		cmdInput:  ci,
	}
}

func (m createModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m createModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.submitted = false
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 1
			}

			if m.focusIndex == 0 {
				m.nameInput.Focus()
				m.cmdInput.Blur()
			} else {
				m.nameInput.Blur()
				m.cmdInput.Focus()
			}

			return m, nil

		case "enter":
			name := m.nameInput.Value()
			cmd := m.cmdInput.Value()

			if name == "" {
				m.err = fmt.Errorf("alias name cannot be empty")
				return m, nil
			}
			if cmd == "" {
				m.err = fmt.Errorf("command cannot be empty")
				return m, nil
			}

			m.submitted = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	if m.focusIndex == 0 {
		m.nameInput, cmd = m.nameInput.Update(msg)
	} else {
		m.cmdInput, cmd = m.cmdInput.Update(msg)
	}

	return m, cmd
}

func (m createModel) View() string {
	if m.submitted {
		return ""
	}

	s := "\n  Create Alias\n\n"
	s += "  " + m.nameInput.View() + "\n"
	s += "  " + m.cmdInput.View() + "\n"

	if m.err != nil {
		s += "\n  ✗ " + m.err.Error() + "\n"
	}

	s += "\n  tab/↑↓ navigate  •  enter confirm  •  esc cancel\n\n"

	return s
}

func runCreateTUI() (name, cmd string, submitted bool) {
	p := tea.NewProgram(initialCreateModel())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ TUI error: %v\n", err)
		os.Exit(1)
	}

	model := m.(createModel)
	if !model.submitted {
		return "", "", false
	}
	return model.nameInput.Value(), model.cmdInput.Value(), true
}

type editModel struct {
	name       string
	currentCmd string
	cmdInput   textinput.Model
	err        error
	submitted  bool
}

func initialEditModel(name, currentCmd string) editModel {
	ci := textinput.New()
	ci.Placeholder = currentCmd
	ci.Prompt = "Command: "
	ci.CharLimit = 200
	ci.Width = 60
	ci.SetValue(currentCmd)
	ci.Focus()
	ci.CursorEnd()

	return editModel{
		name:       name,
		currentCmd: currentCmd,
		cmdInput:   ci,
	}
}

func (m editModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.submitted = false
			return m, tea.Quit

		case "enter":
			cmd := m.cmdInput.Value()
			if cmd == "" {
				m.err = fmt.Errorf("command cannot be empty")
				return m, nil
			}
			m.submitted = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.cmdInput, cmd = m.cmdInput.Update(msg)
	return m, cmd
}

func (m editModel) View() string {
	if m.submitted {
		return ""
	}

	s := fmt.Sprintf("\n  Edit Alias: %s\n\n", m.name)
	s += "  " + m.cmdInput.View() + "\n"

	if m.err != nil {
		s += "\n  ✗ " + m.err.Error() + "\n"
	}

	s += "\n  enter confirm  •  esc cancel\n\n"

	return s
}

func runEditTUI(name, currentCmd string) (newCmd string, submitted bool) {
	p := tea.NewProgram(initialEditModel(name, currentCmd))
	m, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ TUI error: %v\n", err)
		os.Exit(1)
	}

	model := m.(editModel)
	if !model.submitted {
		return "", false
	}
	return model.cmdInput.Value(), true
}

func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
