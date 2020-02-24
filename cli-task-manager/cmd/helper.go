package cmd

import (
	"fmt"

	"github.com/gophercises/cli-task-manager/models"
)

func printCompletedTasks(tasks []string) {
	for _, task := range tasks {
		fmt.Printf("\u2713 %s\n", task)
	}
}

func printTasks(tasks []models.Task) {
	for i, task := range tasks {
		fmt.Printf("%d. %s\n", i+1, string(task.Value))
	}
}
