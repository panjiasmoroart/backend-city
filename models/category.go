package models

import (
	"time"
)

type Category struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug" gorm:"unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
