# ?? AutoSave Backend

Автоматическое управление накоплениями и кредитами через мультибанковскую интеграцию.

## ?? Быстрый старт

### Требования

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (опционально)

### Установка

```bash
# Клонирование репозитория
git clone https://github.com/autosave/backend.git
cd autosave-backend

# Установка зависимостей
go mod download

# Копирование конфигурации
cp .env.example .env

# Запуск миграций
make migrate-up

# Запуск сервера
make run
```

### Docker

```bash
# Запуск всех сервисов
make docker-up

# Остановка
make docker-down
```

## ?? Структура проекта

```
+-- cmd/api/           # Entry point
+-- internal/          # Внутренний код приложения
¦   +-- bankadapter/   # Адаптеры для банков (VBank, ABank, SBank)
¦   +-- config/        # Конфигурация
¦   +-- handlers/      # HTTP handlers
¦   +-- middleware/    # Middleware
¦   +-- models/        # Модели данных
¦   +-- repository/    # Работа с БД
¦   +-- router/        # Маршрутизация
¦   L-- services/      # Бизнес-логика
+-- migrations/        # SQL миграции
+-- pkg/              # Переиспользуемые пакеты
L-- tests/            # Тесты
```

## ?? Поддерживаемые банки

- VBank (Virtual Bank)
- ABank (Awesome Bank)
- SBank (Smart Bank)

## ?? API Endpoints

### Authentication
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход
- `GET /api/auth/me` - Текущий пользователь

### Banks
- `POST /api/banks/:bankId/connect` - Подключить банк
- `GET /api/banks/connected` - Список подключенных банков
- `POST /api/banks/sync` - Синхронизация данных

### Accounts
- `GET /api/accounts` - Все счета
- `GET /api/accounts/:id/transactions` - Транзакции по счету

### Goals
- `GET /api/goals` - Список целей
- `POST /api/goals` - Создать цель
- `PUT /api/goals/:id` - Обновить цель
- `DELETE /api/goals/:id` - Удалить цель

### Operations
- `POST /api/operations/emergency-withdraw` - Экстренное снятие
- `GET /api/operations` - История операций

## ?? Тестирование

```bash
# Запуск всех тестов
make test

# С покрытием
make test-coverage

# Только unit тесты
go test ./internal/...

# Только integration тесты
go test ./tests/...
```

## ?? Команды Make

```bash
make run            # Запустить приложение
make build          # Собрать бинарник
make test           # Запустить тесты
make lint           # Проверить код
make migrate-up     # Применить миграции
make migrate-down   # Откатить миграции
make docker-up      # Запустить в Docker
make docker-down    # Остановить Docker
```

## ?? Лицензия

MIT

## ?? Команда

Hackathon Team 242
