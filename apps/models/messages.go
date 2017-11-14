package models

import "time"

type Message struct {
	Id           int32     `gorm:"primary_key autoincrement"; json:"id"`
	Text         string    `json:text`
	SenderId     int32     `json:sender_id`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	SenderDetail User      `gorm:"ForeignKey:Id;AssociationForeignKey:SenderId"`
}
