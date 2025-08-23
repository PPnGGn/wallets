package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"test_wallets/internal/walletsService"

	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	walletService walletsService.WalletsService
}

func NewWalletsHandler(walletService walletsService.WalletsService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// respondWithError отправляет ответ с ошибкой
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// respondWithJSON отправляет JSON-ответ
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// parseQueryParam извлекает и проверяет строковый параметр из URL
func parseQueryParam(r *http.Request, paramName string) (string, bool) {
	value := r.URL.Query().Get(paramName)
	if value == "" {
		return "", false
	}
	return value, true
}

// parseIntParam извлекает и преобразует числовой параметр из URL
func parseIntParam(r *http.Request, paramName string) (int, bool, string) {
	strValue, ok := parseQueryParam(r, paramName)
	if !ok {
		return 0, false, ""
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, false, "invalid integer value"
	}

	return intValue, true, ""
}

func (h *WalletHandler) CreateWallet(c echo.Context) error {
	// Создаем пустой кошелек, так как все данные будут сгенерированы в сервисе
	var wallet walletsService.Wallet

	if err := h.walletService.CreateWallet(&wallet); err != nil {
		respondWithError(c.Response(), http.StatusInternalServerError, err.Error())
		return err
	}

	c.JSON(http.StatusCreated, wallet)
	return nil
}

func (h *WalletHandler) GetBalanceByAddress(c echo.Context) error {
	// Получаем адрес из параметров пути
	address := c.Param("address")
	log.Printf("Received address: %s", address)
	if address == "" {
		respondWithError(c.Response(), http.StatusBadRequest, "no address provided")
		return echo.ErrBadRequest
	}

	balance, err := h.walletService.GetBalanceByAddress(address)
	if err != nil {
		respondWithError(c.Response(), http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"balance": balance,
	})
}

func (h *WalletHandler) TransferFunds(c echo.Context) error {
	var transaction walletsService.Transaction
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
	err = h.walletService.TransferFunds(transaction.From, transaction.To, transaction.Amount)
	if err != nil {
		respondWithError(c.Response(), http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Funds transferred successfully",
	})
}

func (h *WalletHandler) GetLastTransactions(c echo.Context) error {
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

		// Устанавливаем максимальный лимит на количество возвращаемых транзакций
		if count > 1000 {
			count = 1000
		}
	}

	transactions, err := h.walletService.GetLastTransactions(count)
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
