package textInput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#077bff")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
)

type (
	errMsg error
)

type Output struct {
	Output string
}

func (o *Output) update(val string) {
	o.Output = val
}

type model struct {
	textInput textinput.Model
	err       error
	output    *Output
	header    string
}

// constructor
func InitalTextInputModel(output *Output, header string) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		output:    output,
		header:    titleStyle.Render(header),
	}
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// realtime callback
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if len(m.textInput.Value()) > 1 {
				m.output.update(m.textInput.Value())
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// we handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// renders logic from textinput component to the screen
func (m model) View() string {
	if m.output != nil && m.output.Output != "" { // Check if an output has been set
		// If a selection has been made, display the header and the output without showing the input field
		return fmt.Sprintf("%s\n\n%s\n\n", m.header, m.output.Output)
	} else {
		// If no selection has been made, display the input field as normal
		return fmt.Sprintf("%s\n\n%s\n\n", m.header, m.textInput.View())
	}
}
