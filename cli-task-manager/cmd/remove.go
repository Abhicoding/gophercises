package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes incomplete tasks",
	Long:  `Removes task from the list`,
	Run: func(cmd *cobra.Command, args []string) {
		id := strings.Join(args, "")
		db, err := bolt.Open("cli-task-manager.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal(err)
		}
		s := Store{db: db}
		defer s.db.Close()
		task, err := s.RemoveInCompleteTask(id)
		if err != nil {
			fmt.Errorf("%s", err)
			return
		}
		if task == nil {
			fmt.Printf("No task found for id %s\n", id)
			return
		}
		fmt.Printf("Successfully removed task: %s\n", *task)
	},
}
