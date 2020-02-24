package cmd

import (
	"fmt"
	"strings"

	"github.com/gophercises/cli-task-manager/models"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to todo list",
	Long:  `Adds a new task to your list.`,
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")

		s := models.S
		defer s.DB.Close()

		err := s.CreateTask(task)
		if err != nil {
			fmt.Errorf("Failed to add task: %s", err)
			return
		}

		fmt.Printf("Task added successfully: %s\n", task)
	},
}
