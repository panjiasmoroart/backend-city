package models

import (
	"time"
)

type Photo struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	Image       string    `json:"image"`
	Caption     string    `json:"caption"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
