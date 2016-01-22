package main

import (
	log "github.com/Sirupsen/logrus"

	"flag"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	fetch := flag.String("fetch", "", "fetch some meetup, url required")
	printAll := flag.Bool("all", false, "print all quads")
	intersect := flag.String("intersect", "", "find intersecting users from meetups, separated by comma")

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

		// populating graph
		for _, v := range members {
			d.connectMemberMeetup(v, *fetch)
		}

		if err == nil {
			log.WithFields(log.Fields{
				"meetup": *fetch,
				"count":  len(members),
			}).Info("members added!")
		}

		return
	}

	// printing all quads
	if *printAll {
		d.printAllQuads()
		return
	}

	// looking for intersections
	if *intersect != "" {
		meetups := strings.Split(RemoveSpaces(*intersect), ",")
		members, err := d.findIntersectingMembers(meetups)

		if err != nil {
			log.WithFields(log.Fields{
				"meetups": meetups,
				"error":   err.Error(),
			}).Error("failed ")
			return
		}

		for _, v := range members {
			fmt.Printf("Member %s (ID %d) belongs to: %q\n", v.Name, v.ID, strings.Split(*intersect, ","))
		}

		log.WithFields(log.Fields{
			"meetups": *intersect,
			"count":   len(members),
		}).Info("intersection retrieved!")

		return

	}

	// nothing?
	log.WithFields(log.Fields{
		"nothing": "here",
	}).Warn("Nothing to do")
}
