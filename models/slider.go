package models

import (
	"time"
)

type Slider struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
