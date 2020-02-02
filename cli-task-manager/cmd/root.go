/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/boltdb/bolt"
)

var (
	cfgFile string
)

type Store struct {
	db *bolt.DB
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A CLI based todo app",
	Long:  `This is a simple todo app. You can add, do and list commands to manage your tasks. `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("adding task")
	// },
}

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
			fmt.Println("Failed to add task %s\n", task)
			fmt.Errorf("Failed to add task %s\n", err)
			return
		}
		fmt.Printf("Task added successfully: %s\n", task)
	},
}

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(completedCmd)
	rootCmd.AddCommand(removeCmd)

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli-task-manager.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cli-task-manager" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cli-task-manager")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func (s *Store) CreateTask(t string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("incomplete_tasks"))
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		id, _ := b.NextSequence()

		// Persist bytes to users bucket.
		return b.Put(itob(int(id)), []byte(t))
	})
}

func (s *Store) GetTask() (map[string]string, error) {
	result := map[string]string{}
	err := s.db.View(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("incomplete_tasks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			result[sFromb(k)] = string(v)
		}

		return nil
	})
	return result, err
}

func (s *Store) CompleteTask(id string) (*string, error) {
	tx, err := s.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bOne := tx.Bucket([]byte("incomplete_tasks"))
	task := bOne.Get([]byte(id))
	if err = bOne.Delete([]byte(id)); err != nil {
		return nil, err
	}

	bTwo := tx.Bucket([]byte("completed_tasks"))
	t := time.Now()

	if err = bTwo.Put([]byte(t.Format(time.RFC3339)), task); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	T := string(task)
	return &T, nil
}

func (s *Store) GetCompletedTasks() ([]string, error) {
	result := []string{}
	err := s.db.View(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		c := tx.Bucket([]byte("completed_tasks")).Cursor()
		max := []byte(time.Now().Format(time.RFC3339))
		min := []byte(time.Now().AddDate(0, 0, -1).Format(time.RFC3339))

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			result = append(result, string(v))
		}
		return nil
	})
	return result, err
}

func (s *Store) RemoveInCompleteTask(id string) (*string, error) {
	var task string
	return &task, s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("incomplete_tasks"))
		task = string(b.Get([]byte(id)))
		return b.Delete([]byte(id))
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := []byte(strconv.Itoa(v))
	return b
}

func printTasks(tasks map[string]string) {
	for k, task := range tasks {
		fmt.Printf("%s) %s\n", k, task)
	}
}

func printCompletedTasks(tasks []string) {
	for _, task := range tasks {
		fmt.Printf("%s\n", task)
	}
}

func sFromb(b []byte) string {
	return string(b)
}
