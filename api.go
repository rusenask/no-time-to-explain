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

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/")
	})

	return mux
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

	var response intersectResponse

	response.Meetups = mDetails
	b, err := json.Marshal(response)

	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(b)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)

}
