package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all of your incomplete tasks",
	Long:  `Lists all of your incomplete tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("cli-task-manager.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal(err)
		}
		s := Store{db: db}
		defer s.db.Close()
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
