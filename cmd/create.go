package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	cobra "github.com/spf13/cobra"

	"github.com/belieflab/dataview/cmd/ui/textInput"
)

var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(190)).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color(170)).Bold(true)
)

// bind our command to the existing cobra command
func init() {
	rootCmd.AddCommand(createCmd)
}

type listOptions struct {
	options []string
}

type Options struct {
	ProjectName *textInput.Output // * is a pointer
	ProjectType string
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  ".",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		options := Options{
			ProjectName: &textInput.Output{},
		}

		tprogram := tea.NewProgram(textInput.InitalTextInputModel(options.ProjectName, "Project Name:"))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}
	},
}
