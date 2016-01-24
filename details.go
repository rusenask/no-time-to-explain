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

// Get - searches for given key in the cache and returns value if found
func (c *DetailsDB) Get(key []byte) (value []byte, err error) {

	err = c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(c.bucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", c.bucket)
		}
		// "Byte slices returned from Bolt are only valid during a transaction."
		var buffer bytes.Buffer
		val := bucket.Get(key)

		// If it doesn't exist then it will return nil
		if val == nil {
			return fmt.Errorf("not found")
		}

		buffer.Write(val)
		value = buffer.Bytes()
		return nil
	})

	return
}

// GetAllMembers - returns all captured requests/responses
func (c *DetailsDB) GetAllMembers() (payloads []Member, err error) {
	err = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucket)
		if b == nil {
			// bucket doesn't exist
			return nil
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			pl, err := decodeMember(v)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
					"json":  v,
				}).Warning("Failed to deserialize bytes to payload.")
			} else {
				payloads = append(payloads, pl)
			}
		}
		return nil
	})
	return
}
