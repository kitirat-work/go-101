package otp

import "github.com/gofiber/fiber/v2"

type OtpHandler interface {
	CreateSession(c *fiber.Ctx) error
	ClearSession(c *fiber.Ctx) error

	GetOtp(c *fiber.Ctx) error
	VerifyOtp(c *fiber.Ctx) error

	GetDocumentTypes(c *fiber.Ctx) error
	RequestDocument(c *fiber.Ctx) error
}

type otpHandler struct {
	service OtpService
}

func NewOtpHandler(service OtpService) OtpHandler {
	return &otpHandler{
		service: service,
	}
}

// CreateSession implements OtpHandler.
func (o *otpHandler) CreateSession(c *fiber.Ctx) error {
	ctx := c.Context()
	var req CreateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// get transactionId from header
	transactionId := c.Get("transactionId")
	if transactionId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing transactionId header")
	}

	if pass, err := o.service.VerifyHumanity(ctx, req.TurnstileToken, transactionId, c.IP()); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to verify humanity")
	} else if !pass {
		return fiber.NewError(fiber.StatusForbidden, "failed to verify humanity")
	}

	if err := o.service.CreateSession(ctx, req); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create session")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ClearSession implements OtpHandler.
func (o *otpHandler) ClearSession(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetOtp implements OtpHandler.
func (o *otpHandler) GetOtp(c *fiber.Ctx) error {
	panic("unimplemented")
}

// VerifyOtp implements OtpHandler.
func (o *otpHandler) VerifyOtp(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetDocumentTypes implements OtpHandler.
func (o *otpHandler) GetDocumentTypes(c *fiber.Ctx) error {
	panic("unimplemented")
}

// RequestDocument implements OtpHandler.
func (o *otpHandler) RequestDocument(c *fiber.Ctx) error {
	panic("unimplemented")
}
