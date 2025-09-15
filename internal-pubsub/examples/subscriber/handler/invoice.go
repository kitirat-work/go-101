package handler

import (
	"context"
	"fmt"
	"internal-pubsub/examples/subscriber/model"
	"internal-pubsub/pkg/pubsub"
)

type InvoiceHandler interface {
	CreateInvoice(ctx context.Context, ev pubsub.Event) error
}

type invoiceHandler struct {
	bus             pubsub.Bus
	topicInvoiceGen pubsub.Topic
}

func NewInvoiceHandler(bus pubsub.Bus, topicInvoiceGen pubsub.Topic) InvoiceHandler {
	return &invoiceHandler{bus: bus, topicInvoiceGen: topicInvoiceGen}
}

func (h *invoiceHandler) CreateInvoice(ctx context.Context, ev pubsub.Event) error {
	order := ev.Data.(model.Order)
	fmt.Println("[invoice] received order:", order.ID)
	inv := model.Invoice{ID: "INV-001", OrderID: order.ID, Amount: 1990}
	return h.bus.Publish(ctx, h.topicInvoiceGen, inv)
}
