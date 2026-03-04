package actors

import (
	"fmt"
	"time"

	"ecommerce-actors/messages"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
)

type CoordinatorActor struct {
	remoting    *remote.Remote
	cartPID     *actor.PID
	purchasePID *actor.PID
}

func NewCoordinatorActor(remoting *remote.Remote) actor.Actor {
	return &CoordinatorActor{remoting: remoting}
}

func (c *CoordinatorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		fmt.Println("[Coordinator] started, spawning remote actors...")

		// Remotely spawnuj CartActor na workeru
		cartResp, err := c.remoting.SpawnNamed("127.0.0.1:8090", "cart-actor", "cart", time.Second*5)
		if err != nil {
			fmt.Println("[Coordinator] ERROR spawning CartActor:", err)
			return
		}
		c.cartPID = cartResp.Pid
		fmt.Println("[Coordinator] CartActor spawned:", c.cartPID)

		// Remotely spawnuj PurchaseActor na workeru
		purchaseResp, err := c.remoting.SpawnNamed("127.0.0.1:8090", "purchase-actor", "purchase", time.Second*5)
		if err != nil {
			fmt.Println("[Coordinator] ERROR spawning PurchaseActor:", err)
			return
		}
		c.purchasePID = purchaseResp.Pid
		fmt.Println("[Coordinator] PurchaseActor spawned:", c.purchasePID)

		// Prosledi PurchaseActor PID CartActor-u
		ctx.Send(c.cartPID, &messages.SetPurchasePID{Pid: c.purchasePID})

	case *actor.Stopping:
		fmt.Println("[Coordinator] stopping")

	case *messages.UserAddItem:
		if c.cartPID == nil {
			fmt.Println("[Coordinator] cartPID not ready yet")
			return
		}
		ctx.Send(c.cartPID, &messages.AddItem{
			ProductId: msg.ProductId,
			Quantity:  msg.Quantity,
		})

	case *messages.UserCheckout:
		if c.cartPID == nil {
			fmt.Println("[Coordinator] cartPID not ready yet")
			return
		}
		// Uzimamo ReplyTo iz UserCheckout poruke, ne iz ctx.Sender()
		ctx.Send(c.cartPID, &messages.Checkout{
			ReplyTo: msg.ReplyTo,
		})

	case *messages.OrderConfirmed:
		fmt.Println("[Coordinator] forwarding OrderConfirmed to parent")
		ctx.Send(ctx.Parent(), msg)
	}
}