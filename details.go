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

