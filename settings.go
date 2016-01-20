package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

// Configuration - it is what it is
type Configuration struct {
	port           string
	appKey         string
	databaseName   string
	meetupEndpoint string
}

// DefaultPort - as name suggest
const DefaultPort = "8080"

// DefaultDatabaseName - default database name that will be created
const DefaultDatabaseName = "requests.db"

const DefaultMeetupEndpoint = "https://api.meetup.com/2/"

func InitSettings() *Configuration {
	var appConfig Configuration

	appConfig.port = DefaultPort
	appConfig.databaseName = DefaultDatabaseName
	appConfig.meetupEndpoint = DefaultMeetupEndpoint

	appConfig.appKey = os.Getenv("MeetupKey")

	if appConfig.appKey == "" {
		log.Fatal("Failed to get Meetup API key")
	}

	return &appConfig
}
