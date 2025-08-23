package walletsService

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateAddress генерирует случайный адрес кошелька
func GenerateAddress() (string, error) {
	bytes := make([]byte, 32) // 256 бит для адреса
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(bytes)[:40], // Ограничиваем длину и добавляем префикс
		nil
}

type WalletsService interface {
	CreateWallet(wallet *Wallet) error
	GetLastTransactions(n int) ([]Transaction, error)
	GetBalanceByAddress(address string) (string, error)
	TransferFunds(from string, to string, amount string) error
	InitializeWallets() error
}

func NewWalletService(repo WalletsRepository) WalletsService {
	return &walletsService{repo: repo}
}

type walletsService struct {
	repo WalletsRepository
}

func (w *walletsService) CreateWallet(wallet *Wallet) error {
	// Генерируем новый уникальный адрес
	address, err := GenerateAddress()
	if err != nil {
		return err
	}
	wallet.Address = address
	wallet.Balance = "100.00" // Устанавливаем начальный баланс
	return w.repo.CreateWallet(wallet)
}

func (w *walletsService) GetLastTransactions(n int) ([]Transaction, error) {
	return w.repo.GetLastTransactions(n)
}

func (w *walletsService) GetBalanceByAddress(address string) (string, error) {
	return w.repo.GetBalanceByAddress(address)

}

func (w *walletsService) TransferFunds(from string, to string, amount string) error {
	return w.repo.TransferFunds(from, to, amount)
}

func (w *walletsService) InitializeWallets() error {
	for i := 0; i < 10; i++ {
		address, err := GenerateAddress()
		if err != nil {
			return err
		}
		wallet := &Wallet{
			Address: address,
			Balance: "100.00",
		}
		if err := w.repo.CreateWallet(wallet); err != nil {
			return err
		}
	}
	return nil
}
