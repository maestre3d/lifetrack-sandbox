package main

import (
	"context"
	"log"
	"time"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/google/uuid"
	"github.com/maestre3d/lifetrack-sanbox/pkg/application/occurrence"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence"
)

func main() {
	repo := persistence.NewOccurrenceInMemory()
	handler := occurrence.NewCreateCommandHandler(repo)

	id, err := handler.Invoke(occurrence.CreateCommand{
		Ctx:        context.Background(),
		ActivityID: uuid.New().String(),
		StartTime:  time.Now().Unix(),
		EndTime:    time.Now().Add(time.Minute * 32).Unix(),
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	log.Printf("category created: %s", id)

	query := occurrence.NewGetQuery(repo)
	oc, err := query.Query(context.Background(), id)
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	log.Printf("%+v", oc)

	time.Sleep(time.Second * 3)

	updateHandler := occurrence.NewUpdateCommandHandler(repo)
	err = updateHandler.Invoke(occurrence.UpdateCommand{
		Ctx:        context.Background(),
		ID:         id,
		StartTime:  time.Now().Add(time.Hour * 1).Unix(),
		EndTime:    time.Now().Add(time.Minute * 88).Unix(),
		ActivityID: uuid.New().String(),
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	oc2, err := query.Query(context.Background(), id)
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}

	log.Printf("%+v", oc2)

	/*
		oc, err := aggregate.NewOccurrence(uuid.New().String(), time.Now(), time.Now().Add(time.Minute*10))
		if err != nil {
			log.Fatal(exception.GetDescription(err))
		}

		ocJSON, err := oc.MarshalJSON()
		if err != nil {
			panic(err)
		}
		log.Print(string(ocJSON))

		time.Sleep(time.Second * 3)

		if err := oc.EditTimes(time.Now().Add(time.Hour*1), time.Now().Add(time.Minute*84)); err != nil {
			log.Fatal(exception.GetDescription(err))
		}

		for i, e := range oc.PullEvents() {
			log.Println(fmt.Sprintf("event %d", i))
			eJSON, err := e.MarshalBinary()
			if err != nil {
				log.Fatal(err)
			}

			log.Print(string(eJSON))
		}*/
}