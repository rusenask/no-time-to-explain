package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/google/cayley"
)

// DBClient provides access to cache, http client and configuration
type Handler struct {
	http *http.Client
	cfg  *Configuration
	g    *cayley.Handle
	qs   *cayley.QuadStore
}

// Member struct holds information about each member
type Member struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Link    string  `json:"link"`
	Country string  `json:"country"`
	City    string  `json:"city"`
	Joined  int     `json:"joined"`
	Lon     float64 `json:"lon"`
	Lat     float64 `json:"lat"`
	Visited int     `json:"visited"`
	Status  string  `json:"status"`
}

// ResponseMeta holds vital information for navigating through members API
type ResponseMeta struct {
	Next       string `json:"next"` // next URL to get more members
	TotalCount int    `json:"total_count"`
	Count      int    `json:"count"`
}

// MembersResponse - this is top level structure returned by API
type MembersResponse struct {
	Meta    ResponseMeta `json:"meta"`
	Results []Member     `json:"results"`
}

// getMembers - gets all members for given meetup
func (h *Handler) getMembers(groupName string, pageSize int) ([]Member, error) {
	// creating initial url
	url := fmt.Sprintf("%smembers?group_urlname=%s&page=%d&key=%s", h.cfg.meetupEndpoint, groupName, pageSize, h.cfg.appKey)

	return h._getMembers(url, pageSize)

}

// _getMembers - recursively dives into meetup, fetching all pages till the end
func (h *Handler) _getMembers(url string, pageSize int) ([]Member, error) {
	//	https://api.meetup.com/2/members?group_urlname=frontend&page=200&key=343f7567781b654151e2c635c5445a&order=name
	//	url := fmt.Sprintf("%smembers?group_urlname=%s&page=%d&key=%s", h.cfg.meetupEndpoint, groupName, pageSize, h.cfg.appKey)

	var members []Member

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"url":   url,
			"key":   h.cfg.appKey,
		}).Error("failed to create request")
		return members, err
	}

	resp, err := h.http.Do(request)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"url":   url,
			"key":   h.cfg.appKey,
		}).Error("failed to query API")
		return members, err
	}

	mr := MembersResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"url":   url,
			"key":   h.cfg.appKey,
		}).Error("failed to read body")
		return members, err
	}

	err = json.Unmarshal(body, &mr)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"url":   url,
			"key":   h.cfg.appKey,
		}).Error("failed to unmarshal response from API")
		return members, err
	}

	if mr.Meta.Next != "" {
		moreMembers, err := h._getMembers(mr.Meta.Next, pageSize)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
				"url":   mr.Meta.Next,
				"key":   h.cfg.appKey,
			}).Error("failed to follow trail of members")
		} else {
			mr.Results = append(mr.Results, moreMembers...)
		}

	}

	return mr.Results, nil
}

// connectMemberMeetup - connects members with meetups
func (h *Handler) connectMemberMeetup(member Member, meetup string) (err error) {

	// [member] ----follows----> [meetup]
	err = h.addQuad(strconv.Itoa(member.ID), "follows", meetup)

	// adding kind and details
	err = h.addQuad(strconv.Itoa(member.ID), "kind", "user")
	err = h.addQuad(strconv.Itoa(member.ID), "named", member.Name)
	err = h.addQuad(strconv.Itoa(member.ID), "lives", member.City)
	err = h.addQuad(strconv.Itoa(member.ID), "marked", member.Status)
	return
}
