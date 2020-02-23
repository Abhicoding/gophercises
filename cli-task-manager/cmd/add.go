package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to todo list",
	Long:  `Adds a new task to your list.`,
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")

		db, err := bolt.Open("cli-task-manager.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal(err)
		}
		s := Store{db: db}
		defer s.db.Close()

		err = s.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("incomplete_tasks"))
			if err != nil {
				return fmt.Errorf("create bucket: %s\n", err)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		err = s.CreateTask(task)
		if err != nil {
			fmt.Printf("Failed to add task %s\n", task)
			fmt.Printf("Failed to add task %s\n", err.Error())
			return
		}
		fmt.Printf("Task added successfully: %s\n", task)
	},
}
