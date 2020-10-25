package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/google/uuid"
	"github.com/maestre3d/lifetrack-sanbox/pkg/application/activity"
	"github.com/maestre3d/lifetrack-sanbox/pkg/application/category"
	"github.com/maestre3d/lifetrack-sanbox/pkg/application/occurrence"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/configuration"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/eventbus"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/inmemactivity"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/inmemcategory"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/inmemoccurrence"
)

func main() {
	catRepo := inmemcategory.NewInMemory()
	eventBus := eventbus.NewInMemory(configuration.NewConfiguration())

	// Category domain
	userRef := uuid.New().String()
	physicsRef := addCategory1(userRef, category.NewAddCommandHandler(catRepo, eventBus))
	log.Printf("category created: %s", physicsRef)

	q := category.NewGetQuery(catRepo)
	physics, err := q.Query(context.Background(), physicsRef)
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	log.Printf("%+v", physics)

	// Activity domain
	initActivityWorkflow(physicsRef, inmemactivity.NewInMemory(), eventBus, inmemoccurrence.NewInMemory())
}

func addCategory1(userID string, h *category.AddCommandHandler) string {
	id, err := h.Invoke(category.AddCommand{
		Ctx:    context.Background(),
		UserID: userID,
		Name:   "Physics",
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	return id
}

func addCategory2(userID string, h *category.AddCommandHandler) string {
	id, err := h.Invoke(category.AddCommand{
		Ctx:    context.Background(),
		UserID: userID,
		Name:   "Sports",
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

func initActivityWorkflow(physicsRef string, r *inmemactivity.InMemory, b *eventbus.InMemory,
	rOc *inmemoccurrence.InMemory) {
	quantumActRef := addActivity1(physicsRef, activity.NewAddCommandHandler(r, b))
	log.Printf("activity created: %s", quantumActRef)
	relativityActRef := addActivity2(physicsRef, activity.NewAddCommandHandler(r, b))
	log.Printf("activity created: %s", relativityActRef)

	query := activity.NewGetQuery(r)
	_, err := query.Query(context.Background(), quantumActRef)
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	// log.Printf("%+v", actM)

	list := activity.NewListQuery(r)
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

	// Occurrence domain
	initOccurrenceWorkflow(quantumActRef, rOc, b)
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
		StartTime:  time.Now().Add(time.Minute * 32).Unix(),
		EndTime:    time.Now().Add(time.Minute * 52).Unix(),
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	return id
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

func initOccurrenceWorkflow(activityID string, r *inmemoccurrence.InMemory, b *eventbus.InMemory) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go subscribeToOccurrences(wg, b)
	oc1 := addOccurrence1(activityID, occurrence.NewAddCommandHandler(r, b))
	log.Printf("occurrence created: %s", oc1)

	oc2 := addOccurrence2(activityID, occurrence.NewAddCommandHandler(r, b))
	log.Printf("occurrence created: %s", oc2)
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
		log.Printf("%+v", oc)
	}
}
