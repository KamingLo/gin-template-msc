package models

import (
	"template/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Dont forget to binding the atributte that needed to fill from body
	Username string `json:"username" binding:"required"`
	Email    string `gorm:"unique" json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GenerateCustomID("us", 4)
	return nil
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
