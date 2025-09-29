package otp

type CreateSessionRequest struct {
	CitizenID      string `json:"citizen_id"`
	Phone          string `json:"phone"`
	TurnstileToken string `json:"turnstile_token"`
}
