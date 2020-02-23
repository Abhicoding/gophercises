package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as complete",
	Long:  `Marks a task as complete`,
	Run: func(cmd *cobra.Command, args []string) {
		id := strings.Join(args, " ")

		db, err := bolt.Open("cli-task-manager.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal(err)
		}
		s := Store{db: db}
		defer s.db.Close()

		err = s.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("completed_tasks"))
			if err != nil {
				return fmt.Errorf("create bucket: %s\n", err)
			}
			return nil
		})

		t, err := s.CompleteTask(id)
		if err != nil {
			fmt.Errorf("%s", err)
			return
		}
		fmt.Printf("You have completed the \"%s\" task.\n", *t)
		return
	},
}
