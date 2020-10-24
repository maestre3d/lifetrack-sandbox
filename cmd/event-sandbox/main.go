package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/google/uuid"
	"github.com/maestre3d/lifetrack-sanbox/pkg/application/occurrence"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/configuration"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/eventbus"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/dboccurrence"
)

func main() {
	repo := dboccurrence.NewInMemory()
	eventBus := eventbus.NewInMemory(configuration.NewConfiguration())
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go subscribeToOccurrences(wg, eventBus)
	actRef := addOccurrences(occurrence.NewAddCommandHandler(repo, eventBus))
	wg.Wait()

	list := occurrence.NewListQuery(repo)
	ocs, _, err := list.Query(context.Background(), occurrence.Filter{
		ActivityID: actRef,
		Limit:      10,
		Token:      "",
	})
	if err != nil {
		panic(err)
	}

	for _, oc := range ocs {
		log.Print(oc)
	}
}

func subscribeToOccurrences(wg *sync.WaitGroup, eventBus *eventbus.InMemory) {
	ops := 2
	eventStream, err := eventBus.SubscribeTo(context.Background(), "lifetrack.tracker.occurrence.created")
	if err != nil {
		panic(err)
	}
	for ops > 0 {
		select {
		case ev := <-eventStream:
			evJSON, _ := ev.MarshalBinary()
			log.Print(string(evJSON))
			ops--
		}
	}
	wg.Done()
}

func addOccurrences(handler *occurrence.AddCommandHandler) string {
	act := uuid.New().String()
	id, err := handler.Invoke(occurrence.AddCommand{
		Ctx:        context.Background(),
		ActivityID: act,
		StartTime:  time.Now().Unix(),
		EndTime:    time.Now().Add(time.Minute * 32).Unix(),
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}
	log.Printf("category created: %s", id)

	id2, err := handler.Invoke(occurrence.AddCommand{
		Ctx:        context.Background(),
		ActivityID: act,
		StartTime:  time.Now().Add(time.Minute * 15).Unix(),
		EndTime:    time.Now().Add(time.Minute * 48).Unix(),
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}
	log.Printf("category created: %s", id2)

	return act
}
