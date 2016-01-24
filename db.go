package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/google/cayley"
	"github.com/google/cayley/graph"
	_ "github.com/google/cayley/graph/leveldb"
)

// InitDB - initialises graph database
func InitDB(databaseName string) *cayley.Handle {

	databseType := "leveldb"

	// Initialize the database
	graph.InitQuadStore(databseType, databaseName, nil)

	// Open and use the database
	db, err := cayley.NewGraph(databseType, databaseName, nil)

	if err != nil {
		log.WithFields(log.Fields{
			"error":       err.Error(),
			"databseType": databseType,
			"databseName": databaseName,
		}).Fatal("failed to init DB")
	}

	return db
}

// addQuad - adds quad to graph. In our case - it abstract graph database
// so if something bad happens - we can swap it to different DB.
// [subject] ---predicate---> [object]
func (h *Handler) addQuad(subject, predicate, object string) (err error) {

	quad := cayley.Quad(subject, predicate, object, "")
	err = h.g.AddQuad(quad)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"subject":   subject,
			"predicate": predicate,
			"object":    object,
		}).Error("failed to add quad")
		return
	}
	log.WithFields(log.Fields{
		"subject":   subject,
		"predicate": predicate,
		"object":    object,
	}).Info("quad added")

	return
}

func (h *Handler) printAllQuads() {
	iter := h.g.QuadsAllIterator()
	n, _ := iter.Size()

	for i := int64(0); i < n; i++ {
		graph.Next(iter)
		r := iter.Result()
		q := h.g.Quad(r)
		fmt.Printf("%s - %s -> %s [%s]\n", q.Subject, q.Predicate, q.Object, q.Label)
	}
}
