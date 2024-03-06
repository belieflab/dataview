package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cobra "github.com/spf13/cobra"

	"dataview/cmd/ui/multiInput"
	"dataview/cmd/ui/textInput"
)

var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
)

// bind our command to the existing cobra command
func init() {
	rootCmd.AddCommand(createCmd)
}

type listOptions struct {
	options []string
}

type Options struct {
	ProjectName *textInput.Output     // * is a pointer
	ProjectType *multiInput.Selection // * is a pointer
}

// bare bones implementation of the cobra command
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

		listOfOptions := listOptions{
			options: []string{"Python", "Go", "Node", "Java", "C++", "C#"},
		}

		tprogram = tea.NewProgram(multiInput.InitalModelMulti(listOfOptions.options, options.ProjectType, "Project Type:"))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}
	},
}
