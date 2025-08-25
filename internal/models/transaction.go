package models

import "time"

type Transaction struct {
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	From   string `json:"from" gorm:"index;size:64;not null"`
	To     string `json:"to" gorm:"index;size:64;not null"`
	Amount string `json:"amount" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime;not null"`
}
