package main

import (
	"log"
	"os"
	"wallets/internal/db"
	"wallets/internal/handlers"
	"wallets/internal/models"
	"wallets/internal/repository"
	"wallets/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Инициализация бд
	database, err := db.InitDB()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Автоматическая миграция схемы БД
	if err := database.AutoMigrate(&models.Wallet{}, &models.Transaction{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	walletRepository := repository.NewWalletRepository(database)
	walletService := service.NewWalletService(walletRepository, database)

	// Проверка наличия кошельков в бд
	count, err := walletRepository.CountWallets()
	if err != nil {
		log.Fatalf("Failed to count wallets: %v", err)
	}

	// Создание кошельков, если их 0
	if count == 0 {
		log.Println("First run, creating 10 wallets")
		if err := walletService.InitializeWallets(); err != nil {
			log.Fatalf("Failed to initialize wallets: %v", err)
		}
	}

	walletHandler := handlers.NewWalletsHandler(walletService)

	e := echo.New()
	// Логирование http запросов
	e.Use(middleware.Logger())
	// Обработка паник
	e.Use(middleware.Recover())

	e.POST("/api/create_wallet", walletHandler.CreateWallet)
	
	e.GET("/api/wallet/:address/balance", walletHandler.GetBalance)
	e.POST("/api/send", walletHandler.CreateTransaction)
	e.GET("/api/transactions", walletHandler.GetLast)

	appPort := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.Start(appPort))
}
