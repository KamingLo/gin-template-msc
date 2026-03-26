package models

import (
	"template/utils"
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = utils.GenerateCustomID("bk", 4)
	return nil
}
