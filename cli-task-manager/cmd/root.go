package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli-task-manager",
	Short: "A CLI based todo app",
	Long:  `This is a simple todo app. You can add, do and list commands to manage your tasks. `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(completedCmd)
	rootCmd.AddCommand(removeCmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
