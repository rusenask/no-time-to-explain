package main

import (
	log "github.com/Sirupsen/logrus"

	"flag"
	"net/http"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})

	fetch := flag.String("fetch", "", "fetch some meetup, url required")

	flag.Parse()

	// getting settings
	cfg := InitSettings()

	d := Handler{
		http: &http.Client{},
		cfg:  cfg,
	}

	// deciding what to do
	if *fetch != "" {
		log.WithFields(log.Fields{
			"meetup": *fetch,
		}).Info("Fetching meetup data!")

		members, err := d.getMembers(*fetch, 200)

		if err == nil {
			for _, v := range members {
				log.Printf("Member %s, ID %d ", v.Name, v.ID)
			}
			return
		}

		return
	}

	// nothing?
	log.WithFields(log.Fields{
		"nothing": "here",
	}).Warn("Nothing to do")
}
