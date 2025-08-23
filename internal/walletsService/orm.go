package walletsService

import (
	"time"
)

// Wallet представляет кошелек пользователя
type Wallet struct {
	Address   string    `json:"address" gorm:"primaryKey;size:64"`
	Balance   string    `json:"balance" gorm:"type:varchar(100);not null;default:'100.00'"`
}

// Transaction представляет транзакцию между кошельками
type Transaction struct {
	ID        int      `json:"id" gorm:"primaryKey;autoIncrement"`
	From      string    `json:"from" gorm:"index;size:64;not null"`
	To        string    `json:"to" gorm:"index;size:64;not null"`
	Amount    string    `json:"amount" gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
