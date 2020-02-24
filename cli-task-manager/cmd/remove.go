package cmd

import (
	"fmt"
	"strconv"

	"github.com/gophercises/cli-task-manager/models"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes incomplete tasks",
	Long:  `Removes task from the list`,
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int

		s := models.S
		defer s.DB.Close()

		allTasks, _ := s.GetTask()

		for _, k := range args {
			id, err := strconv.Atoi(k)
			if err != nil || id <= 0 || id > len(allTasks) {
				fmt.Printf("Invalid task number: %s\n", k)
				continue
			}
			ids = append(ids, id)
		}

		for _, id := range ids {
			task := allTasks[id-1]
			err := s.RemoveIncompleteTask(task.Key)
			if err != nil {
				fmt.Errorf("Something went wrong: %s", err)
				continue
			}
			fmt.Printf("Successfully removed task: %s\n", task.Value)
		}
		return
	},
}
