package main

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"ecommerce-actors/actors"
)

func main() {
	system := actor.NewActorSystem()

	r := remote.NewRemote(
		system,
		remote.Configure("127.0.0.1", 8090),
	)
	r.Start()

	// Registrujemo aktore po imenu - master ce ih remotely spawnovati
	r.Register("cart", actor.PropsFromProducer(actors.NewCartActor))
	r.Register("purchase", actor.PropsFromProducer(actors.NewPurchaseActor))

	log.Println("WORKER NODE pokrenut na 8090, ceka remote spawn...")
	select {}
}