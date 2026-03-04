package main

import (
	"log"
	"time"

	"ecommerce-actors/messages"

	"github.com/asynkron/protoactor-go/actor"
)

func main() {
	system := actor.NewActorSystem()

	// 1. Remote PID (adresira aktora na serveru)
	pid := actor.NewPID(
		"127.0.0.1:8080",
		"user-actor",
	)

	// 2. Poruka
	msg := &messages.CreateUser{
		Username: "pera",
		Email:    "pera@test.com",
	}

	// 3. Slanje poruke
	system.Root.Send(pid, msg)

	log.Println("Poruka poslata serveru")

	time.Sleep(1 * time.Second)
}