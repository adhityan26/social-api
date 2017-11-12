package models

import "time"

type Connection struct {
	Id                  int32      `gorm:"primary_key autoincrement"; json:"id"`
	UserId              int32      `json:user_id`
	FriendId            int32      `json:friend_id`
	CreatedAt			time.Time  `json:"created_at"`
	UpdatedAt			time.Time  `json:"updated_at"`
	FriendDetail		User	   `gorm:"ForeignKey:Id;AssociationForeignKey:UserId"`
}
