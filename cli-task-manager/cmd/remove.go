package cmd

import (
	"fmt"
	"strings"

	"github.com/gophercises/cli-task-manager/models"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes incomplete tasks",
	Long:  `Removes task from the list`,
	Run: func(cmd *cobra.Command, args []string) {
		id := strings.Join(args, "")

		s := models.S
		defer s.DB.Close()

		task, err := s.RemoveIncompleteTask(id)
		if err != nil {
			fmt.Errorf("Something went wrong: %s", err)
			return
		}
		if task == nil {
			fmt.Printf("No task found for id %s\n", id)
			return
		}
		fmt.Printf("Successfully removed task: %s\n", *task)
	},
}
