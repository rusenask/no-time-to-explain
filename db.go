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
