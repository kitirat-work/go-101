package humanverify

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

const siteVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type TurnstileClient interface {
	Verify(req VerifyRequest) (*VerifyResponse, error)
}

type Client struct {
	Secret    string
	UserAgent string
}

func NewClient(secret string) TurnstileClient {
	return &Client{
		Secret:    secret,
		UserAgent: "turnstile-fiber-go/1.1",
	}
}

func (c *Client) Verify(req VerifyRequest) (*VerifyResponse, error) {
	if c.Secret == "" {
		return nil, errors.New("turnstile: missing secret")
	}
	if req.Token == "" {
		return nil, &Error{Message: "turnstile: empty token"}
	}

	// payload JSON (Cloudflare รับ form หรือ JSON ก็ได้)
	payload := fiber.Map{
		"secret":   c.Secret,
		"response": req.Token,
	}
	if ip := req.RemoteIP; ip != nil && !ip.IsUnspecified() {
		payload["remoteip"] = ip.String()
	}
	if req.IdempotencyKey != "" {
		payload["idempotency_key"] = req.IdempotencyKey
	}

	// เรียกด้วย Fiber Agent:
	// - ใช้ Post(url) → JSON(v) → UserAgent(...) → Set(...headers) → Bytes()
	//   แล้ว Unmarshal body เอง (pattern ตามเอกสาร) :contentReference[oaicite:2]{index=2} :contentReference[oaicite:3]{index=3}
	agent := fiber.Post(siteVerifyURL).
		JSON(payload).         // ตั้ง Content-Type: application/json อัตโนมัติ :contentReference[oaicite:4]{index=4}
		UserAgent(c.UserAgent) // ตั้ง UA ให้ชัดเจน :contentReference[oaicite:5]{index=5}
	if req.IdempotencyKey != "" {
		agent.Set("Idempotency-Key", req.IdempotencyKey) // header เสริมตอน retry
	}

	status, body, errs := agent.Bytes() // ได้ status, body, errs slice ตามสไตล์ Agent :contentReference[oaicite:6]{index=6}
	if len(errs) > 0 {
		return nil, &Error{Message: "turnstile: request failed: " + errs[0].Error()}
	}
	if status < 200 || status >= 300 {
		return nil, &Error{Message: "turnstile: non-2xx from siteverify"}
	}

	var vr VerifyResponse
	if err := json.Unmarshal(body, &vr); err != nil {
		return nil, &Error{Message: "turnstile: decode response failed: " + err.Error()}
	}

	if !vr.Success {
		return nil, &Error{
			Message: "turnstile: verification failed",
			Codes:   vr.ErrorCodes,
		}
	}

	// เสริมความเข้มงวดตาม best-practices Cloudflare (hostname/action)
	if req.ExpectAction != "" && vr.Action != "" && vr.Action != req.ExpectAction {
		return nil, &Error{Message: "turnstile: action mismatch"}
	}
	if req.ExpectHostname != "" && vr.Hostname != "" && vr.Hostname != req.ExpectHostname {
		return nil, &Error{Message: "turnstile: hostname mismatch"}
	}

	return &vr, nil
}
