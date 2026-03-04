package main

import (
	"log"
	"os"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"ecommerce-actors/actors"
)

func main() {
	host := os.Getenv("WORKER_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("WORKER_PORT")
	if port == "" {
		port = "8090"
	}

	portInt := 8090
	if port == "8090" {
		portInt = 8090
	}

	system := actor.NewActorSystem()
	r := remote.NewRemote(
		system,
		remote.Configure(host, portInt),
	)
	r.Start()

	r.Register("cart", actor.PropsFromProducer(actors.NewCartActor))
	r.Register("purchase", actor.PropsFromProducer(actors.NewPurchaseActor))

	log.Printf("WORKER NODE pokrenut na %s:%d\n", host, portInt)
	select {}
}