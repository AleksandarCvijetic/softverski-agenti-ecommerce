package main

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"ecommerce-actors/actors"
)

func main() {
	system := actor.NewActorSystem()

	r := remote.NewRemote(
		system,
		remote.Configure("127.0.0.1", 8080),
	)
	r.Start()

	// Malo cekamo da worker node bude spreman
	time.Sleep(500 * time.Millisecond)

	coordinatorPID, err := system.Root.SpawnNamed(
		actor.PropsFromProducer(func() actor.Actor {
			return actors.NewCoordinatorActor(r)
		}),
		"coordinator",
	)
	if err != nil {
		log.Fatal("Greška pri spawnovanju Coordinator-a:", err)
	}

	system.Root.Spawn(
		actor.PropsFromProducer(func() actor.Actor {
			return actors.NewUserActor(coordinatorPID)
		}),
	)

	log.Println("MASTER NODE pokrenut na 8080")
	select {}
}