package otp

import "database/sql"

type OtpRepository interface {
}

type otpRepository struct {
	db *sql.DB
}

func NewOtpRepository(db *sql.DB) OtpRepository {
	return &otpRepository{
		db: db,
	}
}
