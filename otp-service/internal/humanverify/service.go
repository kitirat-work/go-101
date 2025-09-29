package humanverify

import (
	"context"
	"fmt"
)

type HumanVerifyService interface {
	Verify(ctx context.Context, req VerifyRequest) (bool, error)
}

type humanVerifyService struct {
	client TurnstileClient
}

func NewHumanVerifyService(client TurnstileClient) HumanVerifyService {
	return &humanVerifyService{
		client: client,
	}
}

// Verify implements HumanVerifyService.
func (h *humanVerifyService) Verify(ctx context.Context, req VerifyRequest) (bool, error) {
	resp, err := h.client.Verify(req)
	if err != nil {
		return false, err
	}

	if resp == nil {
		return false, fmt.Errorf("turnstile: empty response")
	}

	return resp.Success, nil
}
