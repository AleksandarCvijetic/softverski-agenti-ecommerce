package main

import (
	"fmt"
	"time"

	"ecommerce-actors/actors"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

func main() {
	fmt.Println("=== SERVER START ===")

	system := actor.NewActorSystem()

	// REMOTE
	remoteConfig := remote.Configure("127.0.0.1", 8080)
	remote.NewRemote(system, remoteConfig)
	remote.Start(system)

	root := system.Root

	// 1. Spawn Coordinator (CENTRALA)
	coordinatorPID := root.SpawnNamed(
		actor.PropsFromProducer(actors.NewCoordinatorActor),
		"coordinator",
	)

	// 2. Spawn UserActor i prosledi mu coordinator PID
	root.SpawnNamed(
		actor.PropsFromProducer(func() actor.Actor {
			return actors.NewUserActor(coordinatorPID)
		}),
		"user-actor",
	)

	fmt.Println("SERVER pokrenut na 127.0.0.1:8080")

	for {
		time.Sleep(time.Second)
	}
}