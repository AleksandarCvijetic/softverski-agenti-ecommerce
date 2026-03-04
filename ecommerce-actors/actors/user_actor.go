package actors

import (
	"fmt"
	"time"

	"ecommerce-actors/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type UserActor struct {
	coordinatorPID *actor.PID
}

func NewUserActor(coordinator *actor.PID) actor.Actor {
	return &UserActor{coordinatorPID: coordinator}
}

func (u *UserActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		fmt.Println("[UserActor] started")

		// Čekamo da Coordinator spawnuje remote aktore
		time.Sleep(2 * time.Second)

		ctx.Send(u.coordinatorPID, &messages.UserAddItem{
			ProductId: "chocolate",
			Quantity:  2,
		})

		ctx.Send(u.coordinatorPID, &messages.UserAddItem{
			ProductId: "cookies",
			Quantity:  1,
		})

		ctx.Send(u.coordinatorPID, &messages.UserCheckout{
			ReplyTo: ctx.Self(),
		})

	case *messages.OrderConfirmed:
		fmt.Println("[UserActor] order confirmed! OrderId:", msg.OrderId)
	}
}