package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "completed",
	Short: "Lists completed tasks",
	Long:  `Gives a list of tasks completed in last 24 hrs.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("cli-task-manager.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal(err)
		}
		s := Store{db: db}
		defer s.db.Close()

		t, err := s.GetCompletedTasks()
		if err != nil {
			fmt.Errorf("%s", err)
			return
		}
		if t == nil {
			fmt.Printf("No tasks completed in last 24hrs :(\n")
			return
		}
		printCompletedTasks(t)
		return
	},
}
