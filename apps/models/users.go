package models

import (
	"time"
)

type User struct {
	Id					int32   `gorm:"primary_key autoincrement"; json:"id"`
	Name				string	`json:"name" binding:"required"`
	Email				string	`json:"email" binding:"required"`
	Status				bool	`json:"status"`
	CreatedAt			time.Time 	`json:"created_at"`
	UpdatedAt			time.Time  `json:"updated_at"`
}

type UserOutput struct {
	Id					int32 	`json:"id"`
	Name				string	`json:"type"`
	Email				string	`json:"path"`
	Status				bool	`json:"status"`
	CreatedAt			time.Time 	`json:"created_at"`
	UpdatedAt			time.Time  `json:"updated_at"`
}
