package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"time"
)

func main() {
	// Create an occurrence with total duration of 24 hours
	// 15 min >= Occurrence <= 1440 min (24 h)
	oc := &Occurrence{
		ID:            uuid.New().String(),
		StartTime:     time.Now().Add(time.Hour * 24).UTC(),
		EndTime:       time.Now().Add(time.Hour*48).UTC(),
		TotalDuration: 0,
	}
	oc.TotalDuration = int(oc.EndTime.Sub(oc.StartTime).Minutes())

	oc2 := &Occurrence{
		ID:            uuid.New().String(),
		StartTime:     time.Now().Add(time.Hour * 24).UTC(),
		EndTime:       time.Now().Add(time.Hour*40).UTC(),
		TotalDuration: 0,
	}
	oc2.TotalDuration = int(oc2.EndTime.Sub(oc2.StartTime).Minutes())


	// Create an activity with associated occurrence
	// 1440 min (24 h) >= Activity <= 518,400 min (1 y)
	act := &Activity{
		ID:            uuid.New().String(),
		Name:          "Heisenberg Uncertainty Principle",
		AppointedTime: 518400,
		Occurrences: []*Occurrence{oc},
	}

	// Length = 7 days (1 week)
	act2 := &Activity{
		ID:            uuid.New().String(),
		Name:          "Pauli Exclusion Principle",
		AppointedTime: 10080,
		Occurrences: []*Occurrence{oc2},
	}
	
	cat := &Category{
		ID:          uuid.New().String(),
		Title:       "Quantum Mechanics",
		Description: "Quantum mechanics is a fundamental theory in physics that provides a description of the physical " +
			"properties of nature at the scale of atoms and subatomic particles",
		Activities: []*Activity{act, act2},
	}

	catJSON, _ := json.Marshal(cat)
	log.Print(string(catJSON))
}

type Category struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Activities []*Activity `json:"activities"`
}

type Activity struct {
	ID string `json:"id"`
	Name string `json:"title"`
	// In min.
	AppointedTime int `json:"appointed_time"`
	Occurrences []*Occurrence `json:"occurrences"`
}

type Occurrence struct {
	ID string `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime time.Time `json:"end_time"`
	TotalDuration int `json:"total_duration"`
}
