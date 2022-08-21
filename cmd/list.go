/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists existing config for different git accounts",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := cmd.Flags().GetString("host")
		dirname, err := os.UserHomeDir()

		if err != nil {
			fmt.Println("unable to read home directory")
			os.Exit(1)
		}

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

		if len(host) > 0 {
			scanner := bufio.NewScanner(file)
			start := false
			for scanner.Scan() {
				text := scanner.Text()
				textArray := strings.Split(text, " ")

				if start && len(textArray[0]) == 0 {
					fmt.Println(text)
				} else {
					start = false
				}

				if strings.ToLower(textArray[0]) == "host" && textArray[1] == host {
					start = true
					fmt.Println(text)
				}
			}
		} else {
			b, _ := ioutil.ReadAll(file)
			fmt.Print(string(b))
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.PersistentFlags().String("host", "", "see config for a host")
}
