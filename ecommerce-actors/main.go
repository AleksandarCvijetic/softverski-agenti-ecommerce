package main

import (
	"fmt"

	"ecommerce-actors/actors"

	"github.com/asynkron/protoactor-go/actor"
)

func main() {

	fmt.Println("=== E-COMMERCE ACTOR SYSTEM ===")

	system := actor.NewActorSystem()
	root := system.Root

	coordinatorPID := root.Spawn(
		actor.PropsFromProducer(actors.NewCoordinatorActor),
	)

	root.Spawn(
		actor.PropsFromProducer(func() actor.Actor {
			return actors.NewUserActor(coordinatorPID)
		}),
	)

	select {}
}