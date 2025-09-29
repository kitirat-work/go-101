package otp

import (
	"context"
	"fmt"
	"net"
	"otp/internal/humanverify"
)

type OtpService interface {
	VerifyHumanity(ctx context.Context, token string, uuid string, remoteIP string) (bool, error)
	CreateSession(ctx context.Context, req CreateSessionRequest) error
}

type otpService struct {
	repo OtpRepository

	humanVerifyService humanverify.HumanVerifyService
}

func NewOtpService(repo OtpRepository, humanVerifyService humanverify.HumanVerifyService) OtpService {
	return &otpService{
		repo:               repo,
		humanVerifyService: humanVerifyService,
	}
}

// VerifyHumanity implements OtpService.
func (o *otpService) VerifyHumanity(ctx context.Context, token string, uuid string, remoteIP string) (bool, error) {
	req := humanverify.VerifyRequest{
		Token: token,
	}

	if remoteIP != "" {
		ip := net.ParseIP(remoteIP)
		if ip == nil {
			return false, fmt.Errorf("invalid remote IP: %s", remoteIP)
		}
		req.RemoteIP = ip
	}

	if uuid != "" {
		req.IdempotencyKey = uuid
	}

	return o.humanVerifyService.Verify(ctx, req)
}

// CreateSession implements OtpService.
func (o *otpService) CreateSession(ctx context.Context, req CreateSessionRequest) error {
	return nil
}
