package db

import (
	"fmt"
	"log"
	"os"
	"test_wallets/internal/walletsService"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load("/Users/ppnggn/dev/test_wallets/.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Автомиграция схем
	err = db.AutoMigrate(
		&walletsService.Wallet{},
		&walletsService.Transaction{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db, nil

}
