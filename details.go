package main

import (
	"bytes"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
)

// MembersBucketName - default name for members bucket name
const MembersBucketName = "memberbucket"

// DetailsDB - provides access to BoltDB and holds current bucket name
type DetailsDB struct {
	db     *bolt.DB
	bucket []byte
}

// getDB - return bolt db connection
func getDB(name string) *bolt.DB {
	log.WithFields(log.Fields{
		"databaseName": name,
	}).Info("Initiating database")
	db, err := bolt.Open(name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Set - saves given key and value pair to cache
func (c *DetailsDB) Set(key, value []byte) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(c.bucket)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
