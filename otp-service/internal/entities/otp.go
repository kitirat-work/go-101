package entities

type OtpChannel string

const (
	OtpChannelSMS   OtpChannel = "sms"
	OtpChannelEmail OtpChannel = "email"
)

type OtpCode struct {
	ID          int64      `db:"id" gorm:"primaryKey;autoIncrement"`
	SessionID   string     `db:"session_id" gorm:"column:session_id;size:64;not null"`
	Channel     OtpChannel `db:"channel" gorm:"column:channel;type:enum('sms','email');not null"`
	OtpRef      string     `db:"otp_ref" gorm:"column:otp_ref;size:64;not null"`
	OtpCodeHash []byte     `db:"otp_code_hash" gorm:"column:otp_code_hash;size:32;not null"`
	MaxAttempts int        `db:"max_attempts" gorm:"column:max_attempts;not null;default:6"`
	Attempts    int        `db:"attempts" gorm:"column:attempts;not null;default:0"`
	ExpiresAt   string     `db:"expires_at" gorm:"column:expires_at;type:datetime(6);not null"`
	CreatedAt   string     `db:"created_at" gorm:"column:created_at;type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
	UpdatedAt   string     `db:"updated_at" gorm:"column:updated_at;type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (OtpCode) TableName() string {
	return "otp_code"
}
