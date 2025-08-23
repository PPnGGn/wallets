package walletsService

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type WalletsRepository interface {
	CreateWallet(wallet *Wallet) error
	GetLastTransactions(n int) ([]Transaction, error)
	GetBalanceByAddress(address string) (string, error)
	TransferFunds(from string, to string, amount string) error
	CreateTransaction(transaction *Transaction) error
}

func NewWalletRepository(db *gorm.DB) WalletsRepository {
	return &walletsRepository{db: db}
}

type walletsRepository struct {
	db *gorm.DB
}

func (r *walletsRepository) GetBalanceByAddress(address string) (string, error) {
	log.Printf("Searching for wallet with address: '%s'", address)
	var wallet Wallet
	println(address)
	if err := r.db.Where("address = ?", address).First(&wallet).Error; err != nil {
		log.Printf("Error finding wallet with address '%s': %v", address, err)
		return "", err
	}
	log.Printf("Found wallet: %+v", wallet)
	return wallet.Balance, nil
}

func (r *walletsRepository) CreateWallet(wallet *Wallet) error {
	log.Printf("Creating wallet: %+v", wallet)
	err := r.db.Create(wallet).Error
	if err != nil {
		log.Printf("Error creating wallet: %v", err)
	} else {
		log.Printf("Wallet created successfully")
	}
	return err
}

func (r *walletsRepository) GetLastTransactions(n int) ([]Transaction, error) {
	var transactions []Transaction
	if err := r.db.Order("created_at desc").Limit(n).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *walletsRepository) TransferFunds(from string, to string, amount string) error {
	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("ошибка при начале транзакции: %v", tx.Error)
	}

	// Получаем кошелек отправителя
	var fromWallet Wallet
	if err := tx.Where("address = ?", from).First(&fromWallet).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("кошелек отправителя не найден")
		}
		return fmt.Errorf("ошибка при получении кошелька отправителя: %v", err)
	}

	// Получаем кошелек получателя
	var toWallet Wallet
	if err := tx.Where("address = ?", to).First(&toWallet).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("кошелек получателя не найден")
		}
		return fmt.Errorf("ошибка при получении кошелька получателя: %v", err)
	}

	// Парсим сумму перевода
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("некорректный формат суммы перевода")
	}

	// Парсим текущие балансы
	fromBalance, err := strconv.ParseFloat(fromWallet.Balance, 64)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при обработке баланса отправителя")
	}

	toBalance, err := strconv.ParseFloat(toWallet.Balance, 64)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при обработке баланса получателя")
	}

	// Проверяем достаточно ли средств у отправителя
	if fromBalance < amountFloat {
		tx.Rollback()
		return fmt.Errorf("недостаточно средств на счете. Доступно: %.2f, требуется: %.2f", fromBalance, amountFloat)
	}

	// Обновляем балансы
	fromWallet.Balance = strconv.FormatFloat(fromBalance-amountFloat, 'f', 2, 64)
	toWallet.Balance = strconv.FormatFloat(toBalance+amountFloat, 'f', 2, 64)

	// Сохраняем изменения
	if err := tx.Save(&fromWallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при обновлении баланса отправителя: %v", err)
	}

	if err := tx.Save(&toWallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при обновлении баланса получателя: %v", err)
	}

	// Создаем запись о транзакции
	transaction := &Transaction{
		From:   from,
		To:     to,
		Amount: amount,
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка при создании записи о транзакции: %v", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("ошибка при подтверждении транзакции: %v", err)
	}

	return nil
}

func (r *walletsRepository) CreateTransaction(transaction *Transaction) error {
	return r.db.Create(transaction).Error
}
