package cmd

import (
	"fmt"

	"github.com/gophercises/cli-task-manager/models"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of your incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		s := models.S
		defer s.DB.Close()

		tasks, err := s.GetTask()
		if err != nil {
			fmt.Errorf("Failed to get tasks: %s\n", err)
			return
		}
		if len(tasks) == 0 {
			fmt.Printf("No pending tasks for today!\n")
			return
		}

		printTasks(tasks)
	},
}
