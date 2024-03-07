package multiInput

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8")).Background(lipgloss.Color("#030303")).Bold(true)
	titleStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#00ADD8")).Foreground(lipgloss.Color("#030303")).Bold(true)
)

type Selection struct {
	Choice string
}

func (s *Selection) Update(value string) {
	s.Choice = value
}

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
	choice   *Selection
	header   string
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return nil
}

// constructor
func InitalModelMulti(choices []string, selection *Selection, header string) model {

	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
		choice:   selection,
		header:   titleStyle.Render(header),
	}
}

// realtime callback
// realtime callback
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ": // Confirm selection immediately.
			// Update the choice with the current selection.
			if m.choice != nil { // Check that m.choice is not nil
				m.choice.Update(m.choices[m.cursor])
			}
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	return m, nil
}

// renders logic from textinput component to the screen
// renders logic from textinput component to the screen
func (m model) View() string {
	s := m.header + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
		}
		// No need for 'checked' since selection is immediate upon 'enter'
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nPress Enter to select your choice.\n" // Updated instructions
	return s
}
