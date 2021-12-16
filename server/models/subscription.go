package models

import "github.com/jinzhu/gorm"

type Subscription struct {
	gorm.Model
	Following   User `gorm:"foreignKey:FollowingId"`
	FollowingId uint
	Follower    User `gorm:"foreignKey:FollowerId"`
	FollowerId  uint
}

func (Subscription) TableName() string {
	return "subscriptions"
}
