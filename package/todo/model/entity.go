package model

import "time"

type Todo struct {
	ID        int64  `gorm:"primaryKey;autoIncrement;not null"`
	Title     string `gorm:"type:varchar(255);not null"`
	UserID    int64  `gorm:"not null"`
	CreatedAt time.Time
}
