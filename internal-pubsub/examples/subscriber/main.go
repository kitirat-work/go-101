package main

import (
	"context"
	"internal-pubsub/examples/subscriber/handler"
	"internal-pubsub/examples/subscriber/model"
	"internal-pubsub/pkg/pubsub"
	"internal-pubsub/pkg/subscriber"
	"time"
)

const (
	TopicOrderCreated pubsub.Topic = "order.created"
	TopicInvoiceGen   pubsub.Topic = "invoice.generated"
)

func main() {
	bus := pubsub.New(pubsub.DefaultOptions())

	invHandler := handler.NewInvoiceHandler(bus, TopicInvoiceGen)
	mailHandler := handler.NewMailHandler()

	// Invoice service: ฟัง order.created -> สร้าง invoice -> publish invoice.generated
	invSub := subscriber.New(bus, TopicOrderCreated, 8)
	go invSub.Run(context.Background(), invHandler.CreateInvoice)

	// Mailer: ฟัง invoice.generated (sequential)
	mailSub := subscriber.New(bus, TopicInvoiceGen, 4, subscriber.WithStopOnError(false))
	go mailSub.Run(context.Background(), mailHandler.SendMail)

	// Publisher
	_ = bus.Publish(context.Background(), TopicOrderCreated, model.Order{ID: "ORD-123", User: "alice"})

	time.Sleep(200 * time.Millisecond)

	// Shutdown
	_ = bus.Close(context.Background())
	invSub.Close()
	mailSub.Close()
}
