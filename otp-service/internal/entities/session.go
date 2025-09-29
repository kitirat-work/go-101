package entities

type SessionStatus string

const (
	SessionStatusHumanVerified   SessionStatus = "human_verified"
	SessionStatusUserVerified    SessionStatus = "user_verified"
	SessionStatusSmsOtpVerified  SessionStatus = "sms_otp_verified"
	SessionStatusDocumentRequest SessionStatus = "document_requested"
	SessionStatusCompleted       SessionStatus = "completed"
)

type Session struct {
	ID            string        `db:"id" gorm:"column:id;primaryKey;size:64"`
	HumanVerified bool          `db:"human_verified" gorm:"column:human_verified;not null;default:0"`
	IdcardHash    []byte        `db:"idcard_hash" gorm:"column:idcard_hash;size:32"`
	PhoneHash     []byte        `db:"phone_hash" gorm:"column:phone_hash;size:32"`
	PhoneVerified bool          `db:"phone_verified" gorm:"column:phone_verified;not null;default:0"`
	EmailHash     []byte        `db:"email_hash" gorm:"column:email_hash;size:32"`
	EmailVerified bool          `db:"email_verified" gorm:"column:email_verified;not null;default:0"`
	Status        SessionStatus `db:"status" gorm:"column:status;type:enum('human_verified','user_verified','sms_otp_verified','document_requested','completed');not null;default:'human_verified'"`
	CreatedAt     string        `db:"created_at" gorm:"column:created_at;not null"`
	ExpiresAt     string        `db:"expires_at" gorm:"column:expires_at;not null"`
	UpdatedAt     string        `db:"updated_at" gorm:"column:updated_at;not null"`
}

func (Session) TableName() string {
	return "session"
}
