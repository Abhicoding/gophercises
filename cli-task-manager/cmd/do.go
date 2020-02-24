package cmd

import (
	"fmt"
	"strconv"

	"github.com/gophercises/cli-task-manager/models"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete",
	Long:  `Marks a task as complete`,
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
			err := s.CompleteTask(task.Key)
			if err != nil {
				fmt.Errorf("%s", err)
				continue
			}
			fmt.Printf("You have completed task: \"%s\" \n", task.Value)
		}
		return
	},
}
