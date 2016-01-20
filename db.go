package main

import (
	"github.com/google/cayley"
	"github.com/google/cayley/graph"
	_ "github.com/google/cayley/graph/bolt"

	log "github.com/Sirupsen/logrus"
)

func InitDB() *cayley.Handle {
	// Initialize the database
	graph.InitQuadStore("bolt", "cayley.db", nil)

	// Open and use the database
	db, err := cayley.NewGraph("bolt", "cayley.db", nil)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to init DB")
	}

	return db
}
