package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	log "github.com/Sirupsen/logrus"
	"github.com/google/cayley"
	"github.com/google/cayley/graph/path"
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
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Link     string  `json:"link"`
	Country  string  `json:"country"`
	City     string  `json:"city"`
	Hometown string  `json:"hometown"`
	Joined   int     `json:"joined"`
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
	Visited  int     `json:"visited"`
	Status   string  `json:"status"`
}

// encode method encodes all exported Member fields to bytes
func (m *Member) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Member) binaryID() []byte {
	return []byte(strconv.Itoa(m.ID))
}

// decodeMember decodes supplied bytes into Member structure
func decodeMember(data []byte) (Member, error) {
	var p Member
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&p)
	if err != nil {
		return p, err
	}
	return p, nil
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

func (h *Handler) findMember(id string) (member Member) {
	member.ID, _ = strconv.Atoi(id)
	// getting name
	p := cayley.StartPath(h.g, id).Out("named")
	it := p.BuildIterator()
	for cayley.RawNext(it) {
		member.Name = h.g.NameOf(it.Result())
	}
	// getting city
	p_lives := cayley.StartPath(h.g, id).Out("lives")
	it_lives := p_lives.BuildIterator()
	for cayley.RawNext(it_lives) {
		member.City = h.g.NameOf(it_lives.Result())
	}

	return
}

// RemoveSpaces - surprisingly removes spaces
func RemoveSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func (h *Handler) _getMasterPath(nodes []string) *path.Path {
	//	 getting initial path
	p := cayley.StartPath(h.g, nodes[0]).In("follows")

	// add what should users follow as well
	for _, node := range nodes[1:] {
		p.Has("follows", node)
	}

	for _, n := range nodes[1:] {
		thisPath := h._getLesserPath(n, nodes)
		p.And(thisPath)
	}

	return p
}

func (h *Handler) _getLesserPath(current string, nodes []string) *path.Path {
	// getting initial path
	p := cayley.StartPath(h.g, current).In("follows")

	// add what should users follow as well
	for _, node := range nodes {
		if node != current {
			p.Has("follows", node)
		}
	}
	return p
}

func (h *Handler) findIntersectingMembers(meetups []string) (members []Member, err error) {
	log.WithFields(log.Fields{
		"meetups0": meetups[0],
		"meetups1": meetups[1],
	}).Info("starting intersect!")

	p := h._getMasterPath(meetups)

	it := p.BuildIterator()
	for cayley.RawNext(it) {
		//		log.Println(h.g.NameOf(it.Result()))
		members = append(members, h.findMember(h.g.NameOf(it.Result())))
	}
	return members, nil
}
