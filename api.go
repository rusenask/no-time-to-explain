package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/go-zoo/bone"
	"github.com/meatballhat/negroni-logrus"
)

type meetupDetailsResponse struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type allMeetupsResponse struct {
	Meetups []string `json:"meetups"`
}

type intersectResponse struct {
	Members     []Member                `json:"members"`
	Intersected int                     `json:"intersected"`
	Meetups     []meetupDetailsResponse `json:"meetups"`
}

func (h *Handler) startAdminInterface() {
	// starting admin interface
	mux := getBoneRouter(*h)
	n := negroni.Classic()

	loglevel := log.InfoLevel

	n.Use(negronilogrus.NewCustomMiddleware(loglevel, &log.JSONFormatter{}, "web"))
	n.UseHandler(mux)

	// admin interface starting message
	log.WithFields(log.Fields{
		"port": h.cfg.port,
	}).Info("web interface is starting...")

	n.Run(fmt.Sprintf(":%s", h.cfg.port))
}

// getBoneRouter returns mux for admin interface
func getBoneRouter(d Handler) *bone.Mux {
	mux := bone.New()

	mux.Get("/api/intersect", http.HandlerFunc(d.IntersectionHandler))
	mux.Get("/api/fetch", http.HandlerFunc(d.FetchMeetupHandler))
	mux.Get("/api/meetups", http.HandlerFunc(d.GetAllMeetupsHandler))

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/")
	})

	return mux
}

// FetchMeetupHandler - fetches supplied meetup:
// http://localhost:8080/api/fetch?meetup=raspberry-pint-london
func (h *Handler) FetchMeetupHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()

	meetup := q["meetup"][0]
	log.WithFields(log.Fields{
		"meetup": meetup,
	}).Info("got query to fetch meetup")

	w.Header().Set("Content-Type", "application/json")

	members, err := h.fetchMeetupData(meetup)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var mr meetupDetailsResponse

	mr.Size = len(members)
	mr.Name = meetup

	b, err := json.Marshal(mr)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(b)
		return
	}

	w.WriteHeader(200)

}

func (h *Handler) GetAllMeetupsHandler(w http.ResponseWriter, req *http.Request) {
	meetups := h.getAllMeetups()

	var mr allMeetupsResponse

	mr.Meetups = meetups

	b, err := json.Marshal(mr)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(b)
		return
	}

	w.WriteHeader(200)
}

// IntersectionHandler returns intersected members for given meetups
// http://localhost:8080/api/intersect?q=kubernetes-london&q=docker-london
func (h *Handler) IntersectionHandler(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()

	log.WithFields(log.Fields{
		"query": q["q"],
	}).Info("got query")

	w.Header().Set("Content-Type", "application/json")

	var mDetails []meetupDetailsResponse
	for _, v := range q["q"] {
		size := h.getTotalFollowersCount(v)
		mDetails = append(mDetails, meetupDetailsResponse{Size: size, Name: v})
	}

	intersectingMembers, err := h.findIntersectingMembers(q["q"])

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response intersectResponse

	response.Meetups = mDetails
	response.Members = intersectingMembers
	response.Intersected = len(intersectingMembers)

	b, err := json.Marshal(response)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(b)
		return
	}

	w.WriteHeader(200)

}
