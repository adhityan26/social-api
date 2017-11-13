package models

import "time"

type Subscribe struct {
	Id                  int32      `gorm:"primary_key autoincrement"; json:"id"`
	RequestorId         int32      `json:requestor_id`
	TargetId         	int32      `json:target_id`
	CreatedAt			time.Time  `json:"created_at"`
	UpdatedAt			time.Time  `json:"updated_at"`
	UserDetail			User	   `gorm:"ForeignKey:Id;AssociationForeignKey:RequestorId"`
	SubscribeDetail		User	   `gorm:"ForeignKey:Id;AssociationForeignKey:TargetId"`
}

