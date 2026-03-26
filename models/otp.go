package models

import (
	"template/utils"
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	ID        string    `gorm:"primaryKey" type:"varchar(36)" json:"id"`
	Email     string    `gorm:"index;not null" json:"email"`
	Code      string    `gorm:"type:varchar(6);not null" json:"code"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (o *OTP) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = utils.GenerateCustomID("otp", 4)
	return nil
}
