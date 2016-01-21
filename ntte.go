package main

import (
	log "github.com/Sirupsen/logrus"

	"flag"
	"net/http"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	fetch := flag.String("fetch", "", "fetch some meetup, url required")
	printAll := flag.Bool("all", false, "print all quads")

	flag.Parse()

	// getting settings
	cfg := InitSettings()
	graph := InitDB()

	d := Handler{
		http: &http.Client{},
		cfg:  cfg,
		g:    graph,
	}

	// deciding what to do
	if *fetch != "" {
		log.WithFields(log.Fields{
			"meetup": *fetch,
		}).Info("Fetching meetup data!")

		members, err := d.getMembers(*fetch, 200)

		if err == nil {
			log.WithFields(log.Fields{
				"meetup": *fetch,
				"count":  len(members),
			}).Info("members added!")
		}

		return
	}

	if *printAll {
		d.printAllQuads()
		return
	}

	// nothing?
	log.WithFields(log.Fields{
		"nothing": "here",
	}).Warn("Nothing to do")
}
