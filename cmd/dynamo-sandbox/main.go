package main

import (
	"context"
	"log"
	"time"

	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/dbpool"

	"github.com/alexandria-oss/common-go/exception"
	"github.com/maestre3d/lifetrack-sanbox/pkg/application/category"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/configuration"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/eventbus"
	"github.com/maestre3d/lifetrack-sanbox/pkg/infrastructure/persistence/dynamocategory"
)

func main() {
	cfg := configuration.NewConfiguration()
	catRepo := dynamocategory.NewDynamo(cfg, dbpool.NewDynamoDBPool(cfg))
	eventBus := eventbus.NewInMemory(cfg)

	addH := category.NewAddCommandHandler(catRepo, eventBus)
	id, err := addH.Invoke(category.AddCommand{
		Ctx:    context.Background(),
		UserID: "aruiz",
		Name:   "Music",
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}
	log.Printf("category %s created", id)

	time.Sleep(time.Second * 15)

	h := category.NewUpdateCommandHandler(catRepo, eventBus)
	err = h.Invoke(category.UpdateCommand{
		Ctx:         context.Background(),
		ID:          id,
		UserID:      "aruiz",
		Name:        "Music II",
		Description: "A brief sample explanation",
		TargetTime:  68,
		Picture:     "https://cdn.damascus-engineering.com/alexandria/user/1b4cc750-c551-4767-a232-e91b52e68fa0.jpeg",
		State:       "true",
	})
	if err != nil {
		log.Fatal(exception.GetDescription(err))
	}
}
