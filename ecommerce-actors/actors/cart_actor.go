package actors

import (
	"fmt"

	"ecommerce-actors/messages"

	"github.com/asynkron/protoactor-go/actor"
)

type CartActor struct {
	behavior    actor.Behavior
	items       map[string]int32
	totalPrice  float64
	purchasePID *actor.PID
}

func NewCartActor() actor.Actor {
	cart := &CartActor{
		behavior:   actor.NewBehavior(),
		items:      make(map[string]int32),
		totalPrice: 0,
	}
	cart.behavior.Become(cart.open)
	return cart
}

func (c *CartActor) Receive(ctx actor.Context) {
	c.behavior.Receive(ctx)
}

func (c *CartActor) open(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		fmt.Println("[CartActor] started")

	case *actor.Stopping:
		fmt.Println("[CartActor] stopping")

	case *actor.Stopped:
		fmt.Println("[CartActor] stopped")

	case *messages.SetPurchasePID:
		c.purchasePID = msg.Pid
		fmt.Println("[CartActor] received PurchaseActor PID")

	case *messages.AddItem:
		c.items[msg.ProductId] += msg.Quantity
		c.totalPrice += float64(msg.Quantity) * 10
		fmt.Println("[CartActor] item added:", msg.ProductId, "| cart:", c.items)

	case *messages.Checkout:
		fmt.Println("[CartActor] checkout started")

		if c.purchasePID == nil {
			fmt.Println("[CartActor] ERROR: purchasePID is nil!")
			return
		}

		ctx.Send(c.purchasePID, &messages.ProcessOrder{
			ReplyTo:    msg.ReplyTo,
			Items:      c.items,
			TotalPrice: c.totalPrice,
		})

		c.behavior.Become(c.checkedOut)
	}
}

func (c *CartActor) checkedOut(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *messages.AddItem:
		fmt.Println("[CartActor] cannot add items, cart already checked out")
	case *messages.Checkout:
		fmt.Println("[CartActor] checkout already completed")
	}
}