package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"wallets/internal/models"
	"wallets/internal/repository"

	"gorm.io/gorm"
)

type WalletsService interface {
	CreateWallet() (*models.Wallet, error)
	GetLast(n int) ([]models.Transaction, error)
	GetBalance(address string) (string, error)
	CreateTransaction(from string, to string, amount string) error
	InitializeWallets() error
}

func NewWalletService(repo repository.WalletsRepository, db *gorm.DB) WalletsService {
	return &walletsService{repo: repo, db: db}
}

type walletsService struct {
	repo repository.WalletsRepository
	db   *gorm.DB
}

// Генерирует рандомный 64-символьный адрес кошелька в 16ричном формате
func generateAddress() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Создает новый кошелек с уникальным адресом и начальным балансом 100.00
func (s *walletsService) CreateWallet() (*models.Wallet, error) {
	address, err := generateAddress()
	if err != nil {
		return nil, err
	}

	wallet := &models.Wallet{
		Address: address,
		Balance: "100.00",
	}
	if err := s.repo.CreateWallet(wallet); err != nil {
		return nil, err
	}
	log.Println("Wallet created:", wallet)
	return wallet, nil
}

// получение списка последних n транзакций
// Если n <= 0, возвращает все доступные транзакции
func (s *walletsService) GetLast(n int) ([]models.Transaction, error) {
	return s.repo.GetLast(n)
}

// Получение баланса кошелька по его адресу
func (s *walletsService) GetBalance(address string) (string, error) {
	wallet, err := s.repo.GetWallet(address)
	if err != nil {
		return "", err
	}
	return wallet.Balance, nil
}

// Создание новой транзакции между кошельками
// Проверяет наличие достаточного баланса у отправителя и существование кошельков
func (s *walletsService) CreateTransaction(from string, to string, amount string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Получаем кошельки
		fromWallet, err := s.repo.GetWallet(from)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("кошелек отправителя не найден")
			}
			return err
		}

		toWallet, err := s.repo.GetWallet(to)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("кошелек получателя не найден")
			}
			return err
		}

		// Конвертируем суммы
		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return fmt.Errorf("некорректный формат суммы перевода")
		}

		fromBalance, _ := strconv.ParseFloat(fromWallet.Balance, 64)
		toBalance, _ := strconv.ParseFloat(toWallet.Balance, 64)

		// Проверка средств
		if fromBalance < amountFloat {
			return fmt.Errorf("недостаточно средств: доступно %.2f, требуется %.2f", fromBalance, amountFloat)
		}

		// Обновляем балансы
		fromWallet.Balance = strconv.FormatFloat(fromBalance-amountFloat, 'f', 2, 64)
		toWallet.Balance = strconv.FormatFloat(toBalance+amountFloat, 'f', 2, 64)

		// Сохраняем
		if err := s.repo.UpdateWallet(fromWallet); err != nil {
			return err
		}
		if err := s.repo.UpdateWallet(toWallet); err != nil {
			return err
		}

		// Создаем запись о транзакции
		txn := &models.Transaction{
			From:   from,
			To:     to,
			Amount: amount,
		}
		if err := s.repo.CreateTransaction(txn); err != nil {
			return err
		}

		return nil
	})
}

// Инициализация 10 тестовых кошельков с балансом 100.00
// Используется для первоначального заполнения базы данных
func (s *walletsService) InitializeWallets() error {
	for i := 0; i < 10; i++ {
		if _, err := s.CreateWallet(); err != nil {
			return err
		}
	}
	return nil
}
