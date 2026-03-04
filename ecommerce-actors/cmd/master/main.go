package main

import (
	"log"
	"os"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"ecommerce-actors/actors"
)

func main() {
	masterHost := os.Getenv("MASTER_HOST")
	if masterHost == "" {
		masterHost = "127.0.0.1"
	}

	workerAddr := os.Getenv("WORKER_ADDR")
	if workerAddr == "" {
		workerAddr = "127.0.0.1:8090"
	}

	system := actor.NewActorSystem()
	r := remote.NewRemote(
		system,
		remote.Configure(masterHost, 8080),
	)
	r.Start()

	time.Sleep(500 * time.Millisecond)

	coordinatorPID, err := system.Root.SpawnNamed(
		actor.PropsFromProducer(func() actor.Actor {
			return actors.NewCoordinatorActor(r, workerAddr)
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

	log.Printf("MASTER NODE pokrenut na %s:8080, worker na %s\n", masterHost, workerAddr)
	select {}
}