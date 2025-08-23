package main

import (
	"log"
	"net/http"
	"test_wallets/internal/db"
	"test_wallets/internal/handlers"
	"test_wallets/internal/walletsService"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Инициализация базы данных
	database, err := db.InitDB()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Автоматическая миграция схемы БД
	if err := database.AutoMigrate(&walletsService.Wallet{}, &walletsService.Transaction{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	walletRepository := walletsService.NewWalletRepository(database)
	walletService := walletsService.NewWalletService(walletRepository)

	// Проверяем, есть ли кошельки в базе
	var count int64
	if err := database.Model(&walletsService.Wallet{}).Count(&count).Error; err != nil {
		log.Fatalf("Failed to count wallets: %v", err)
	}

	// Если кошельков нет, создаем 10 тестовых кошельков
	if count == 0 {
		log.Println("No wallets found, creating 10 test wallets...")
		if err := walletService.InitializeWallets(); err != nil {
			log.Fatalf("Failed to initialize wallets: %v", err)
		}
	}

	walletHandler := handlers.NewWalletsHandler(walletService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/create_wallet", walletHandler.CreateWallet)
	e.GET("/api/wallet/:address/balance", walletHandler.GetBalanceByAddress)
	e.POST("/api/send", walletHandler.TransferFunds)
	e.GET("/api/transactions", walletHandler.GetLastTransactions)

	// TODO: Удалить этот эндпоинт перед продакшеном - только для тестирования
	e.GET("/api/test/wallets", func(c echo.Context) error {
		// Временный хендлер для получения всех кошельков
		var wallets []walletsService.Wallet
		if err := database.Find(&wallets).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch wallets",
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"count":   len(wallets),
			"wallets": wallets,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
