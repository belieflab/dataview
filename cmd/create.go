package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cobra "github.com/spf13/cobra"

	"jspsych/cmd/ui/multiInput"
	"jspsych/cmd/ui/textInput"
)

var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#077bff")).Bold(true)
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
	Short: "Create a new jsPsych Experiment",
	Long:  `This command sets up a new jsPsych project, allowing the user to choose between different jsPsych versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the setup options
		options := Options{
			ProjectName: &textInput.Output{},
			ProjectType: &multiInput.Selection{}, // Make sure this is not nil!
		}

		// Run textInput program to get the project name
		tprogram := tea.NewProgram(textInput.InitalTextInputModel(options.ProjectName, "What would you like to name your new experiment? (Please use camelCase):"))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}

		// Set up and run multiInput program to choose jsPsych version
		listOfOptions := listOptions{
			options: []string{"jsPsych 6.3", "jsPsych 7.x"},
		}
		options.ProjectType = &multiInput.Selection{} // Ensure this is initialized
		tprogram = tea.NewProgram(multiInput.InitalModelMulti(listOfOptions.options, options.ProjectType, "Which version of jsPsych would you like to use?"))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}

		// Output the results
		// fmt.Println("Project Name:", options.ProjectName.Output)
		// fmt.Println("Project Type:", options.ProjectType.Choice) // Make sure this is initialized before use

		// Convert project type "jsPsych 6.3" or "jsPsych 7.x" to just "6.3" or "7.x"
		version := strings.Trim(options.ProjectType.Choice, "jsPsych ")

		// Now, use regular expression to find just the numbers before the dot
		re := regexp.MustCompile(`^\d+`) // This matches one or more digits at the beginning of the string
		trim := re.FindString(version)   // Find the match

		// Initialize Git and clone template files
		gitCommands := []string{
			"git init",
			fmt.Sprintf("git submodule add git@github.com:belieflab/jsPsychWrapper-v7.x.git wrap"),
		}
		fileCommands := []string{
			"mkdir -p ./css",
			"echo \"/* add local styling here */\" >> ./css/exp.css",
			"mkdir -p ./exp",
			fmt.Sprintf("cp -rf ./wrap/tmp/v%s/index.php ./index.php", trim),
			fmt.Sprintf("cp -rf ./wrap/tmp/v%s/timeline.js ./exp/timeline.js", trim),
			fmt.Sprintf("cp -rf ./wrap/tmp/v%s/main.js ./exp/main.js", trim),
			"cp -rf ./wrap/tmp/conf.js ./exp/conf.js",
			"cp -rf ./wrap/tmp/lang.js ./exp/lang.js",
			"cp -rf ./wrap/tmp/var.php ./exp/var.php",
			"echo \"// add local functions here \" >> ./exp/fn.js",
			// Add more file operations as needed
		}
		linkCommands := []string{
			"ln -s ./wrap/link/data.php ./data.php",
			"ln -s ./wrap/link/redirect.php./ redirect.php",
			"ln -s ./wrap/link/sync.sh ./sync.sh",
		}

		// Execute all commands
		for _, cmd := range gitCommands {
			exec.Command("bash", "-c", cmd).Run()
		}
		for _, cmd := range fileCommands {
			exec.Command("bash", "-c", cmd).Run()
		}
		for _, cmd := range linkCommands {
			exec.Command("bash", "-c", cmd).Run()
		}

		// Create data folder and initialize Git repository
		exec.Command("mkdir", "-p", "./data").Run()
		exec.Command("touch", "./data/.gitkeep").Run()

		// Rename project directory
		newName := options.ProjectName.Output // Ensure this is your validated project name
		err := os.Rename("../createExperiment", "../"+newName)
		if err != nil {
			fmt.Printf("WARNING: Failed to rename experiment. Directory already exists with the name %s.\n", newName)
		} else {
			fmt.Println("Experiment renamed successfully to " + newName)
		}

		// After all setup is done
		// Assuming 'jsPsychBinary' is the path to the binary, replace it with the actual path if different
		err = os.Remove("./jspsych") // Modify as necessary
		if err != nil {
			// Handle the error, maybe the file didn't exist or there were permissions issues
			fmt.Printf("WARNING: Failed to remove jsPsych binary: %v.\n", err)
		} else {
			fmt.Println("jsPsych binary removed successfully.")
		}

		// Change working directory to the new project name
		err = os.Chdir("../" + newName)
		if err != nil {
			// Handle the error, maybe the directory doesn't exist or there are permissions issues
			fmt.Printf("WARNING: Failed to change working directory to the new project: %v.\n", err)
		} else {
			fmt.Println("Current working directory changed to the new project.")
		}

		// Git operations
		// gitSetupAndPush(newName, version) // Implement this function for git add, commit, branch, remote add, and push operations
	},
}
