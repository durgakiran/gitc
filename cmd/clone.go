/*
Copyright © 2022 NAME HERE durgakiran
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"gitc/sshclient"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type promptContent struct {
	errorMsg string
	label    string
}

/**
* Alternative to golang, but with the power of choosing which git config to config to clone
* 1. choose which config to clone
*
 */
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedConfig := chooseConfig()
		activateConfig(selectedConfig)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// activates selected config
func activateConfig(selectedConfig string) {
	fmt.Println("selected config", selectedConfig)
	sshclient.SshToGit("/.ssh/id_personal")
	// current key to ssh
	// make sure ssh -T works
}

// asks user to select config
func chooseConfig() string {
	hosts := getHost()
	categoryPromptContent := promptContent{
		"Please select a Host.",
		"Which git config to use?",
	}
	selectedOption := promptGetSelect(categoryPromptContent, hosts)
	return selectedOption
}

// get's all the hosts in the present config
func getHost() []string {
	dirname, err := os.UserHomeDir()
	file, err := os.Open(dirname + "/.ssh/config")
	if err != nil {
		fmt.Println("error reading config file in .ssh directory")
		os.Exit(1)
	}

	defer func() {
		if err = file.Close(); err != nil {
			os.Exit(1)
		}
	}()

	// read line by line
	scanner := bufio.NewScanner(file)
	hosts := make([]string, 0, 10)
	for scanner.Scan() {
		text := scanner.Text()
		textArray := strings.Split(text, " ")

		if strings.ToLower(textArray[0]) == "host" {
			hosts = append(hosts, text)
		}
	}

	return hosts
}

func promptGetSelect(pc promptContent, hosts []string) string {
	items := hosts
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    pc.label,
			Items:    items,
			AddLabel: "Add new config",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}
