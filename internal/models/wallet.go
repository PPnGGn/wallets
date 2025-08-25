package models

type Wallet struct {
	Address string `json:"address" gorm:"primaryKey;size:64"`
	Balance string `json:"balance" gorm:"type:varchar(100);not null;default:'100.00'"`
}
