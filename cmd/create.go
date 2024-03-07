package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	cobra "github.com/spf13/cobra"

	"jspsych/cmd/ui/multiInput"
	"jspsych/cmd/ui/textInput"
)

// bind our command to the existing cobra command
func init() {
	rootCmd.AddCommand(createCmd)
}

type listOptions struct {
	options []string
}

type Options struct {
	ExperimentName *textInput.Output     // * is a pointer
	GitHubAccount  *textInput.Output     // * is a pointer
	jsPsychVersion *multiInput.Selection // * is a pointer
}

// validateCamelCase checks if a string is in camelCase format.
func validateCamelCase(input string) bool {
	// Regex for matching camelCase. This pattern matches strings that start with a lowercase letter
	// followed by any combination of lowercase and uppercase letters.
	matchCamelCase := regexp.MustCompile(`^[a-z]+(?:[A-Z][a-z]*)*$`)
	return matchCamelCase.MatchString(input)
}

// bare bones implementation of the cobra command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new jsPsych Experiment",
	Long:  `This command sets up a new jsPsych project, allowing the user to choose between different jsPsych versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the setup options
		options := Options{
			ExperimentName: &textInput.Output{},
			GitHubAccount:  &textInput.Output{},
			jsPsychVersion: &multiInput.Selection{}, // Make sure this is not nil!
		}

		// Run textInput program to get the project name with camelCase validation
		for {
			// Set up and run textInput program to get the project name
			tprogram := tea.NewProgram(textInput.InitalTextInputModel(options.ExperimentName, "What would you like to name your new experiment? (Please use camelCase):"))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(err) // Handle non-validation related errors
			}

			// Validate the camelCase format for the project name
			if validateCamelCase(options.ExperimentName.Output) {
				break // Break out of the loop if the project name is in camelCase
			} else {
				fmt.Println("Project name is not in camelCase format. Please try again using camelCase.")
				options.ExperimentName.Output = "" // Reset the output to ensure the user can input again
			}
		}

		// Loop until an organization name is entered (since your validation might be different for organization names, adjust as necessary)
		for {
			// Run textInput program to get the organization name
			tprogram := tea.NewProgram(textInput.InitalTextInputModel(options.GitHubAccount, "What is the name of your GitHub organization? (default 'belieflab'):"))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(err)
			}

			// Check if an organization name has been provided; if not, set to 'belieflab'.
			if options.GitHubAccount.Output != "" {
				break // Valid input; break out of the loop
			} else {
				fmt.Println("No GitHub organization name entered. Try Again.")
				options.GitHubAccount.Output = "" // Reset the output to ensure the user can input again
			}
		}

		// Set up and run multiInput program to choose jsPsych version
		listOfOptions := listOptions{
			options: []string{"jsPsych 6.3", "jsPsych 7.x"},
		}

		options.jsPsychVersion = &multiInput.Selection{} // Ensure this is initialized
		tprogram := tea.NewProgram(multiInput.InitalModelMulti(listOfOptions.options, options.jsPsychVersion, "Which version of jsPsych would you like to use?"))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}

		// Convert project type "jsPsych 6.3" or "jsPsych 7.x" to just "6.3" or "7.x"
		re := regexp.MustCompile(`^\d+`) // This matches one or more digits at the beginning of the string
		// Now, use regular expression to find just the numbers before the dot
		version := re.FindString(strings.Trim(options.jsPsychVersion.Choice, "jsPsych "))

		// Initialize Git and clone template files
		gitFork := []string{
			"rm -rf .git",
			"git init",
		}

		fileCommands := []string{
			// Create data folder and initialize .gitkeep
			"mkdir -p ./data",
			"touch ./data/.gitkeep",
			// Style
			"mkdir -p ./css",
			"echo \"/* add local styling here */\" >> ./css/exp.css",
			// Index
			fmt.Sprintf("cp -rf ./wrap/tmp/v%s/index.php ./index.php", version),
			// Experiment dependencies
			"mkdir -p ./exp",
			fmt.Sprintf("cp -rf ./wrap/tmp/v%s/timeline.js ./exp/timeline.js", version),
			fmt.Sprintf("cp -rf ./wrap/tmp/v%s/main.js ./exp/main.js", version),
			"cp -rf ./wrap/tmp/conf.js ./exp/conf.js",
			"cp -rf ./wrap/tmp/lang.js ./exp/lang.js",
			"cp -rf ./wrap/tmp/var.php ./exp/var.php",
			"echo \"// add local functions here \" >> ./exp/fn.js",
		}

		symbolicLinks := []string{
			"ln -s ./wrap/link/data.php ./data.php",
			"ln -s ./wrap/link/redirect.php ./redirect.php",
			"ln -s ./wrap/link/sync.sh ./sync.sh",
		}

		gitAdd := []string{
			"git add .gitignore",
			"git add .gitmodules",
			"git add data/.gitkeep",
			"git add --all",
		}

		gitCommit := []string{
			"git commit -m \"initialized experiment\"",
			"git branch -M main",
			fmt.Sprintf("git remote get-url origin || git remote add origin git@github.com:%s/%s", options.GitHubAccount.Output, options.ExperimentName.Output),
		}

		// Execute all commands
		for _, cmd := range gitFork {
			if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
				fmt.Printf("WARNING: Failed to execute git command '%s': %v\n", cmd, err)
			}
		}
		for _, cmd := range fileCommands {
			if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
				fmt.Printf("WARNING: Failed to execute git command '%s': %v\n", cmd, err)
			}
		}
		for _, cmd := range symbolicLinks {
			if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
				fmt.Printf("WARNING: Failed to execute git command '%s': %v\n", cmd, err)
			}
		}

		// Execute git setup commands
		for _, cmd := range gitAdd {
			if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
				fmt.Printf("WARNING: Failed to execute git command '%s': %v\n", cmd, err)
				// You might want to handle the error more gracefully depending on your application's requirements
			}
		}

		for _, cmd := range gitCommit {
			if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
				fmt.Printf("WARNING: Failed to execute git command '%s': %v\n", cmd, err)
				// You might want to handle the error more gracefully depending on your application's requirements
			}
		}

		// Push changes to the new remote in a loop until successful
		success := false
		for !success {
			err := exec.Command("bash", "-c", "git push -u origin main").Run()
			if err != nil {
				fmt.Println("WARNING: Failed to push changes to the remote GitHub repository.")
				fmt.Printf("Please make sure the repository 'belieflab/%s' has been created on GitHub and you have the correct access rights.\n", options.ExperimentName.Output)
				fmt.Println("Attempting to push again in 30 seconds...")
				time.Sleep(30 * time.Second) // Wait for 30 seconds before trying again
			} else {
				fmt.Println("Changes successfully pushed to GitHub.")
				fmt.Println("Please edit exp/conf.js to configure your experiment.")
				// Rename project directory
				err = os.Rename("../createExperiment", "../"+options.ExperimentName.Output)
				if err != nil {
					fmt.Printf("WARNING: Failed to rename experiment. Directory already exists with the name %s.\n", options.ExperimentName.Output)
				} else {
					fmt.Println("Experiment renamed successfully to " + options.ExperimentName.Output)
				}
				// Change working directory to the new project name
				err = os.Chdir("../" + options.ExperimentName.Output)
				if err != nil {
					// Handle the error, maybe the directory doesn't exist or there are permissions issues
					fmt.Printf("WARNING: Failed to change working directory to the new project: %v.\n", err)
				} else {
					fmt.Println("Current working directory changed to the new project.")
				}
				err = os.Remove("./jspsych") // Modify as necessary
				if err != nil {
					// Handle the error, maybe the file didn't exist or there were permissions issues
					fmt.Printf("WARNING: Failed to remove jsPsych binary: %v.\n", err)
				} else {
					fmt.Println("jsPsych binary removed successfully.")
				}
				success = true // Exit the loop since the push was successful
			}
		}
	},
}
