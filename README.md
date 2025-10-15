# Wallets API Service

Сервис для управления криптовалютными кошельками с REST API интерфейсом.

## Возможности
- Создание 10 кошельков при запуске
- Получение баланса кошелька
- Перевод средств между кошельками
- История транзакций

## Технологии

- **Язык программирования**: Go 1.25+
- **Веб-фреймворк**: Echo v4
- **База данных**: PostgreSQL
- **ORM**: GORM
- **Контейнеризация**: Docker

## Быстрый старт

### Требования

- Docker и Docker Compose
- Go 1.25+
- PostgreSQL

### Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/wallets.git
   cd wallets

2. Запустите приложение с помощью Docker Compose:
   ```bash
   docker-compose up --build -d
   ```

## API Endpoints

- `POST /api/create_wallet` - Создать новый кошелек
- `GET /api/wallet/{address}/balance` - Получить баланс кошелька
- `POST /api/send` - Перевести средства
- `GET /api/transactions` - Получить историю транзакций

## Настройка

Настройки приложения задаются через переменные окружения в файле `.env`:

DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=wallets
SERVER_PORT=8080
