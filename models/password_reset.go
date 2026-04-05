package models

import (
	"template/utils"
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UpdatedAt time.Time `json:"updated_at"`

	Email        string    `gorm:"primaryKey;type:varchar(255)" json:"email"`
	Token        string    `gorm:"type:varchar(255)" json:"token"`
	ExpiredAt    time.Time `json:"expired_at"`
	RequestCount int       `gorm:"default:0" json:"request_count"`
}

func init() {
	RegisterModel(&PasswordReset{})
}

func (p *PasswordReset) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = utils.GenerateCustomID("pr", 6)
	return nil
}
