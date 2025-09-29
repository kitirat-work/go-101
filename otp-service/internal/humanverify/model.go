package humanverify

import "net"

type VerifyRequest struct {
	Token          string
	RemoteIP       net.IP
	IdempotencyKey string
	ExpectAction   string
	ExpectHostname string
}

type VerifyResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Action      string   `json:"action"`
	CData       string   `json:"cdata"`
	ErrorCodes  []string `json:"error-codes"`
}

type Error struct {
	Message string
	Codes   []string
}

func (e *Error) Error() string { return e.Message }
