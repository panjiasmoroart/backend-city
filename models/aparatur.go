package models

import (
	"time"
)

type Aparatur struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	Image       string    `json:"image"`
	Name        string    `json:"name"`
	Position    string    `json:"position"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
