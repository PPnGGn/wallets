package repository

import (
	"wallets/internal/models"

	"gorm.io/gorm"
)

type WalletsRepository interface {
	CreateWallet(wallet *models.Wallet) error
	GetWallet(address string) (*models.Wallet, error)
	UpdateWallet(wallet *models.Wallet) error
	GetLast(n int) ([]models.Transaction, error)
	CreateTransaction(tx *models.Transaction) error
	CountWallets() (int64, error)
}

func NewWalletRepository(db *gorm.DB) WalletsRepository {
	return &walletsRepository{db: db}
}

type walletsRepository struct {
	db *gorm.DB
}

// Создание нового кошелька в бд
func (r *walletsRepository) CreateWallet(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

// Поиск кошелька по адресу
func (r *walletsRepository) GetWallet(address string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("address = ?", address).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

// Обновление данных кошелька в бд
func (r *walletsRepository) UpdateWallet(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

// Получение последних n транзакций из бд
// Если n <= 0, возвращает все доступные транзакции
func (r *walletsRepository) GetLast(n int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.Order("created_at desc").Limit(n).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// Сохранение новой транзакции в бд
func (r *walletsRepository) CreateTransaction(txn *models.Transaction) error {
	return r.db.Create(txn).Error
}

// Получение общего количества кошельков в бд
func (r *walletsRepository) CountWallets() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Wallet{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
