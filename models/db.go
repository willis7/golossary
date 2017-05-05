package models

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

const bucketName = "glossary"

func InitDB(path string) {
	var err error
	// Open the path data file in your current directory.
	// It will be created if it doesn't exist.
	db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucketName: %s", err)
		}
		return nil
	})
}

func Update(word Word) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(word.Name), []byte(word.Description))
		return err
	})
}

func Get(word string) string {
	var v []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v = b.Get([]byte(word))
		return nil
	})
	return string(v)
}
