package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/boltdb/bolt"
)

var (
	S                    Store
	incompleteTaskBucket = []byte("incomplete_tasks")
	completedTaskBucket  = []byte("completed_tasks")
)

//Store : DB singleton
type Store struct {
	DB *bolt.DB
}

//Task : defines a todo struct with value as the task
type Task struct {
	Key   int
	Value string
}

func Init(homePath string) {
	dbFilePath := path.Join(homePath, "cli-task-manager.db")
	db, err := bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	S = Store{DB: db}
	S.createBucketIfNotExists(incompleteTaskBucket)
	S.createBucketIfNotExists(completedTaskBucket)
}

//CreateTask: creates a task
func (s *Store) CreateTask(task string) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(incompleteTaskBucket)
		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		id, _ := b.NextSequence()

		// Persist bytes to users bucket.
		return b.Put(itob(int(id)), []byte(task))
	})
}

//GetTask : returns list of incomplete tasks
func (s *Store) GetTask() ([]Task, error) {
	result := []Task{}
	err := s.DB.View(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket(incompleteTaskBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			result = append(result, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	return result, err
}

//GetCompletedTask : returns list of completed tasks
func (s *Store) GetCompletedTasks() ([]string, error) {
	result := []string{}
	err := s.DB.View(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		c := tx.Bucket(completedTaskBucket).Cursor()
		max := []byte(time.Now().Format(time.RFC3339))
		min := []byte(time.Now().AddDate(0, 0, -1).Format(time.RFC3339))

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			result = append(result, string(v))
		}
		return nil
	})
	return result, err
}

//CompleteTask : completes the task with given id
func (s *Store) CompleteTask(id int) error {
	// return s.DB.Update(func(tx *bolt.Tx) error {
	// 	taskID := itob(id)
	// 	bOne := tx.Bucket(incompleteTaskBucket)
	// 	task := bOne.Get(taskID)
	// 	if err := bOne.Delete(taskID); err != nil {
	// 		return err
	// 	}
	// 	bTwo := tx.Bucket(completedTaskBucket)
	// 	t := time.Now()

	// 	if err := bTwo.Put([]byte(t.Format(time.RFC3339)), task); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
	tx, err := s.DB.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	taskID := itob(id)
	bOne := tx.Bucket(incompleteTaskBucket)
	task := bOne.Get(taskID)
	if err = bOne.Delete(taskID); err != nil {
		return err
	}

	bTwo := tx.Bucket(completedTaskBucket)
	t := time.Now()

	if err = bTwo.Put([]byte(t.Format(time.RFC3339)), task); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

//RemoveIncompleteTask: Removes incomplete task from the list
func (s *Store) RemoveIncompleteTask(id int) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(incompleteTaskBucket)
		return b.Delete(itob(id))
	})
}

//CreateBucketIfNotExists: creates a bucket with given name
func (s *Store) createBucketIfNotExists(bucketName []byte) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s\n", err)
		}
		return nil
	})
}

//itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func stringFromByteArray(b []byte) string {
	return string(b)
}
