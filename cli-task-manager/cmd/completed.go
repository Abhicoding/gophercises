package cmd

import (
	"fmt"

	"github.com/gophercises/cli-task-manager/models"
	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "Lists completed tasks",
	Long:  `Gives a list of tasks completed in last 24 hrs.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := models.S
		defer s.DB.Close()

		t, err := s.GetCompletedTasks()
		if err != nil {
			fmt.Errorf("Something went wrong: %s", err)
			return
		}

		if len(t) == 0 {
			fmt.Printf("No tasks completed in last 24hrs :(\n")
			return
		}

		printCompletedTasks(t)
		return
	},
}
