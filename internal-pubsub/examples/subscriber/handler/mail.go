package handler

import (
	"context"
	"fmt"
	"internal-pubsub/examples/subscriber/model"
	"internal-pubsub/pkg/pubsub"
)

type MailHandler interface {
	SendMail(ctx context.Context, ev pubsub.Event) error
}

type mailHandler struct {
	// สมมติว่ามี service อื่นๆ ที่ใช้ส่งเมล
}

func NewMailHandler() MailHandler {
	return &mailHandler{}
}

func (h *mailHandler) SendMail(ctx context.Context, ev pubsub.Event) error {
	inv := ev.Data.(model.Invoice)
	fmt.Println("[mailer] send email for invoice:", inv.ID, "order:", inv.OrderID)
	return nil
}
