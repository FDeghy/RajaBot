package database

import "time"

func NewSubscription(userId int64) *Subscription {
	return &Subscription{
		UserID:         userId,
		IsTrial:        false,
		ExpirationDate: 0,
		IsEnabled:      false,
		RegisteryDate:  time.Now().Unix(),
	}
}
