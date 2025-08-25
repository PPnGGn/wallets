package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"wallets/internal/models"
	"wallets/internal/service"

	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	walletService service.WalletsService
}

func NewWalletsHandler(walletService service.WalletsService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// Отправка ответа с ошибкой
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// Создание нового кошелька
func (h *WalletHandler) CreateWallet(c echo.Context) error {
	// Создаем новый кошелек
	wallet, err := h.walletService.CreateWallet()
	if err != nil {
		respondWithError(c.Response(), http.StatusInternalServerError, err.Error())
		return err
	}

	c.JSON(http.StatusCreated, wallet)
	return nil
}

// Обработка запроса на получение баланса кошелька
func (h *WalletHandler) GetBalance(c echo.Context) error {
	// Парсинг адреса из пути запроса
	address := c.Param("address")
	log.Printf("Received address: %s", address)
	if address == "" {
		respondWithError(c.Response(), http.StatusBadRequest, "no address provided")
		return echo.ErrBadRequest
	}

	balance, err := h.walletService.GetBalance(address)
	if err != nil {
		respondWithError(c.Response(), http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"balance": balance,
	})
}

// Создание транзакции между кошельками
func (h *WalletHandler) CreateTransaction(c echo.Context) error {
	var transaction models.Transaction
	if err := json.NewDecoder(c.Request().Body).Decode(&transaction); err != nil {
		respondWithError(c.Response(), http.StatusBadRequest, "неверный формат JSON")
		return echo.ErrBadRequest
	}

	// Валидация входных данных
	if transaction.From == "" || transaction.To == "" {
		respondWithError(c.Response(), http.StatusBadRequest, "не указан адрес отправителя или получателя")
		return echo.ErrBadRequest
	}

	if transaction.Amount == "" {
		respondWithError(c.Response(), http.StatusBadRequest, "не указана сумма перевода")
		return echo.ErrBadRequest
	}

	amount, err := strconv.ParseFloat(transaction.Amount, 64)
	if err != nil || amount <= 0 {
		respondWithError(c.Response(), http.StatusBadRequest, "некорректная сумма перевода")
		return echo.ErrBadRequest
	}

	if transaction.From == transaction.To {
		respondWithError(c.Response(), http.StatusBadRequest, "нельзя перевести средства на тот же кошелек")
		return echo.ErrBadRequest
	}

	// Выполняем перевод
	err = h.walletService.CreateTransaction(transaction.From, transaction.To, transaction.Amount)
	if err != nil {
		respondWithError(c.Response(), http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "Успешно",
		"message": "Средства переведены успешно",
	})
}

func (h *WalletHandler) GetLast(c echo.Context) error {
	// Получаем параметр count из query-параметров
	countStr := c.QueryParam("count")
	count := 10 // Значение по умолчанию

	// Если параметр count передан, парсим его
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			errMsg := "параметр count должен быть положительным числом"
			if count < 0 {
				errMsg = "количество транзакций не может быть отрицательным"
			}
			respondWithError(c.Response(), http.StatusBadRequest, errMsg)
			return echo.ErrBadRequest
		}
	}

	transactions, err := h.walletService.GetLast(count)
	if err != nil {
		respondWithError(c.Response(), http.StatusInternalServerError, "не удалось получить список транзакций")
		return err
	}

	// Возвращаем список транзакций
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count":        len(transactions),
		"transactions": transactions,
	})
}
