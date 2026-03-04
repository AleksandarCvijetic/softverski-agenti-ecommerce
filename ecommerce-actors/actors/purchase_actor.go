package actors

import (
	"fmt"

	"ecommerce-actors/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type PurchaseActor struct{}

func NewPurchaseActor() actor.Actor {
	return &PurchaseActor{}
}

func (p *PurchaseActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		fmt.Println("[PurchaseActor] started")

	case *actor.Stopping:
		fmt.Println("[PurchaseActor] stopping")

	case *actor.Stopped:
		fmt.Println("[PurchaseActor] stopped")

	case *messages.ProcessOrder:
		fmt.Println("[PurchaseActor] processing order, items:", msg.Items)

		orderId := fmt.Sprintf("ORDER-%d", len(msg.Items))

		if msg.ReplyTo != nil {
			ctx.Send(msg.ReplyTo, &messages.OrderConfirmed{
				OrderId: orderId,
			})
			fmt.Println("[PurchaseActor] OrderConfirmed sent to:", msg.ReplyTo)
		} else {
			fmt.Println("[PurchaseActor] ERROR: replyTo is nil!")
		}
	}
}