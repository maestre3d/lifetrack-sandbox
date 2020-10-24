package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/application/activity"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/inmemactivity"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/google/uuid"
	"github.com/maestre3d/lifetrack-sanbox/pkg/application/occurrence"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/configuration"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/eventbus"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/inmemoccurrence"
)

func main() {
	//ocRepo := inmemoccurrence.NewInMemory()
	actRepo := inmemactivity.NewInMemory()
	eventBus := eventbus.NewInMemory(configuration.NewConfiguration())

	physicsRef := uuid.New().String()
	act1 := addActivity1(physicsRef, activity.NewAddCommandHandler(actRepo, eventBus))
	log.Printf("activity created: %s", act1)
	act2 := addActivity2(physicsRef, activity.NewAddCommandHandler(actRepo, eventBus))
	log.Printf("activity created: %s", act2)

	query := activity.NewGetQuery(actRepo)
	_, err := query.Query(context.Background(), act1)
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	// log.Printf("%+v", actM)

	list := activity.NewListQuery(actRepo)
	acts, nextToken, err := list.Query(context.Background(), activity.Filter{
		CategoryID: physicsRef,
		Title:      "",
		Limit:      1,
		Token:      "",
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}
	log.Print(nextToken)

	for _, a := range acts {
		log.Printf("%+v", a)
	}

	// initOccurrenceWorkflow(act1, ocRepo, eventBus)
}

func initOccurrenceWorkflow(activityID string, r *inmemoccurrence.InMemory, b *eventbus.InMemory) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go subscribeToOccurrences(wg, b)
	oc1 := addOccurrence1(activityID, occurrence.NewAddCommandHandler(r, b))
	log.Printf("category created: %s", oc1)

	oc2 := addOccurrence2(activityID, occurrence.NewAddCommandHandler(r, b))
	log.Printf("category created: %s", oc2)
	wg.Wait()

	list := occurrence.NewListQuery(r)
	ocs, _, err := list.Query(context.Background(), occurrence.Filter{
		ActivityID: activityID,
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

func addOccurrence1(activityID string, handler *occurrence.AddCommandHandler) string {
	id, err := handler.Invoke(occurrence.AddCommand{
		Ctx:        context.Background(),
		ActivityID: activityID,
		StartTime:  time.Now().Unix(),
		EndTime:    time.Now().Add(time.Minute * 32).Unix(),
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	return id
}

func addOccurrence2(activityID string, handler *occurrence.AddCommandHandler) string {
	id, err := handler.Invoke(occurrence.AddCommand{
		Ctx:        context.Background(),
		ActivityID: activityID,
		StartTime:  time.Now().Add(time.Minute * 15).Unix(),
		EndTime:    time.Now().Add(time.Minute * 48).Unix(),
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	return id
}

func addActivity1(categoryID string, handler *activity.AddCommandHandler) string {
	id, err := handler.Invoke(activity.AddCommand{
		Ctx:        context.Background(),
		CategoryID: categoryID,
		Title:      "Quantum mechanics",
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	return id
}

func addActivity2(categoryID string, handler *activity.AddCommandHandler) string {
	id, err := handler.Invoke(activity.AddCommand{
		Ctx:        context.Background(),
		CategoryID: categoryID,
		Title:      "Relativity theory",
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	return id
}
