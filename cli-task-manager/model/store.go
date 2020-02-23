package models

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var S Store

type Store struct {
	db *bolt.DB
}

func init() {
	db, err := bolt.Open("cli-task-manager.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	s = Store{db: db}
}
