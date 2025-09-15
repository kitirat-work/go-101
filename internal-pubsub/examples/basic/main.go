package main

import (
	"context"
	"fmt"
	"time"

	"internal-pubsub/pkg/pubsub"
)

const (
	TopicOrderCreated pubsub.Topic = "order.created"
	TopicInvoiceGen   pubsub.Topic = "invoice.generated"
)

type Order struct {
	ID   string
	User string
}

type Invoice struct {
	ID      string
	OrderID string
	Amount  int
}

func main() {
	bus := pubsub.New(pubsub.DefaultOptions())

	// --- Module A: Order Service -> Publish order.created ---
	// --- Module B: Invoice Service -> Subscribe order.created, then Publish invoice.generated ---
	// --- Module C: Mailer -> Subscribe invoice.generated ---

	// B: Invoice Service
	subOrder := bus.Subscribe(TopicOrderCreated, 8)
	go func() {
		for ev := range subOrder.C() {
			order := ev.Data.(Order)
			fmt.Println("[invoice] received order:", order.ID)

			// ทำงานจริง...
			inv := Invoice{ID: "INV-001", OrderID: order.ID, Amount: 1990}

			_ = bus.Publish(context.Background(), TopicInvoiceGen, inv)
		}
	}()

	// C: Mailer
	subInvoice := bus.Subscribe(TopicInvoiceGen, 4)
	go func() {
		for ev := range subInvoice.C() {
			inv := ev.Data.(Invoice)
			fmt.Println("[mailer] send email for invoice:", inv.ID, "order:", inv.OrderID)
		}
	}()

	// A: Order Service -> trigger
	_ = bus.Publish(context.Background(), TopicOrderCreated, Order{ID: "ORD-123", User: "alice"})

	time.Sleep(200 * time.Millisecond)

	// ปิดระบบ
	_ = bus.Close(context.Background())

	// Unsubscribe ไม่จำเป็นหลัง Close() แต่เรียกได้
	subOrder.Unsubscribe()
	subInvoice.Unsubscribe()
}
