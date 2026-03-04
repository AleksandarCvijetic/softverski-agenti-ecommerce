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
	workerAddr  string
	cartPID     *actor.PID
	purchasePID *actor.PID
}

func NewCoordinatorActor(remoting *remote.Remote, workerAddr string) actor.Actor {
	return &CoordinatorActor{
		remoting:   remoting,
		workerAddr: workerAddr,
	}
}

func (c *CoordinatorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {

	case *actor.Started:
		fmt.Println("[Coordinator] started, spawning remote actors...")

		cartResp, err := c.remoting.SpawnNamed(c.workerAddr, "cart-actor", "cart", time.Second*5)
		if err != nil {
			fmt.Println("[Coordinator] ERROR spawning CartActor:", err)
			return
		}
		c.cartPID = cartResp.Pid
		fmt.Println("[Coordinator] CartActor spawned:", c.cartPID)

		purchaseResp, err := c.remoting.SpawnNamed(c.workerAddr, "purchase-actor", "purchase", time.Second*5)
		if err != nil {
			fmt.Println("[Coordinator] ERROR spawning PurchaseActor:", err)
			return
		}
		c.purchasePID = purchaseResp.Pid
		fmt.Println("[Coordinator] PurchaseActor spawned:", c.purchasePID)

		ctx.Send(c.cartPID, &messages.SetPurchasePID{Pid: c.purchasePID})

	case *actor.Stopping:
		fmt.Println("[Coordinator] stopping")

	case *messages.UserAddItem:
		if c.cartPID == nil {
			fmt.Println("[Coordinator] cartPID not ready yet")
			return
		}
		fmt.Printf("[Coordinator] routing AddItem(%s, %d) → CartActor na workeru\n", msg.ProductId, msg.Quantity)
		ctx.Send(c.cartPID, &messages.AddItem{
			ProductId: msg.ProductId,
			Quantity:  msg.Quantity,
		})

	case *messages.UserCheckout:
		if c.cartPID == nil {
			fmt.Println("[Coordinator] cartPID not ready yet")
			return
		}
		fmt.Println("[Coordinator] routing Checkout → CartActor na workeru")
		ctx.Send(c.cartPID, &messages.Checkout{
			ReplyTo: msg.ReplyTo,
		})

	case *messages.OrderConfirmed:
		fmt.Println("[Coordinator] forwarding OrderConfirmed to parent")
		ctx.Send(ctx.Parent(), msg)
	}
}