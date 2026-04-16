package models

import (
	"template/utils"
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	ID           string    `gorm:"primaryKey" type:"varchar(36)" json:"id"`
	Email        string    `gorm:"index;not null" json:"email"`
	Code         string    `gorm:"type:varchar(6);not null" json:"code"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`                     // Wajib untuk melacak waktu update terakhir
	RequestCount int       `gorm:"default:1" json:"request_count"` // Penghitung request
}

func init() {
	RegisterModel(&OTP{})
}

func (o *OTP) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = utils.GenerateCustomID("otp", 6)
	return nil
}
