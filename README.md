# Mobile Store Backend

Backend API для магазина мобильных аксессуаров, написанный на Go.

## Технологический стек

- **Go 1.21** - основной язык
- **Gin** - HTTP веб-фреймворк
- **PostgreSQL** - основная база данных
- **Redis** - кэширование и сессии
- **GORM** - ORM для работы с базой данных
- **JWT** - аутентификация
- **Docker** - контейнеризация

## Структура проекта

еба

```
├── internal/
│   ├── config/          # Конфигурация приложения
│   ├── database/        # Подключение к БД
│   ├── handlers/        # HTTP обработчики
│   ├── middleware/      # Middleware
│   ├── models/          # Модели данных
│   ├── repository/      # Слой доступа к данным
│   └── services/        # Бизнес-логика
├── main.go              # Точка входа
├── go.mod               # Зависимости Go
├── Dockerfile           # Docker образ
├── docker-compose.yml   # Docker Compose
└── init.sql            # Инициализация БД
```

## Быстрый старт

### 1. Клонирование и установка зависимостей

```bash
git clone <repository-url>
cd mobile-store-back
go mod tidy
```

### 2. Настройка окружения

Скопируйте `config.env.example` в `config.env` и настройте переменные:

```bash
cp config.env.example config.env
```

### 3. Запуск с Docker Compose

```bash
docker-compose up -d
```

Это запустит:

- PostgreSQL на порту 5432
- Redis на порту 6379
- PgAdmin на порту 5050
- API сервер на порту 8080

### 4. Или запуск локально

```bash
# Запуск PostgreSQL и Redis
docker-compose up -d postgres redis

# Запуск приложения
go run main.go
```

## API Endpoints

### Аутентификация

- `POST /api/v1/auth/register` - Регистрация
- `POST /api/v1/auth/login` - Вход

### Продукты (публичные)

- `GET /api/v1/products` - Список продуктов
- `GET /api/v1/products/:id` - Получить продукт
- `GET /api/v1/products/search?q=query` - Поиск продуктов
- `GET /api/v1/products/category/:category_id` - Продукты по категории

### Пользователи (требует аутентификации)

- `GET /api/v1/users/profile` - Профиль пользователя
- `PUT /api/v1/users/profile` - Обновить профиль

### Заказы (требует аутентификации)

- `POST /api/v1/orders` - Создать заказ
- `GET /api/v1/orders` - Мои заказы
- `GET /api/v1/orders/:id` - Получить заказ

### Админ (требует админских прав)

- `GET /api/v1/admin/users` - Список пользователей
- `POST /api/v1/admin/products` - Создать продукт
- `PUT /api/v1/admin/products/:id` - Обновить продукт
- `DELETE /api/v1/admin/products/:id` - Удалить продукт
- `GET /api/v1/admin/orders` - Все заказы

## Переменные окружения

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=mobile_store

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Server
SERVER_PORT=8080
SERVER_HOST=localhost

# JWT
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRE_HOURS=24

# Environment
ENV=development
```

## Разработка

### Структура кода

Проект следует принципам Clean Architecture:

1. **Handlers** - HTTP слой, обработка запросов
2. **Services** - Бизнес-логика
3. **Repository** - Слой доступа к данным
4. **Models** - Модели данных

### Добавление новых функций

1. Создайте модель в `internal/models/`
2. Добавьте методы в соответствующий репозиторий
3. Создайте сервис для бизнес-логики
4. Добавьте HTTP обработчики
5. Зарегистрируйте маршруты в `handlers.go`

## Тестирование

```bash
go test ./...
```

## Мониторинг

- Логирование через Zap
- Метрики через Prometheus (планируется)
- Трейсинг через Jaeger (планируется)

## Развертывание

### Production

1. Настройте переменные окружения
2. Используйте внешние PostgreSQL и Redis
3. Настройте reverse proxy (nginx)
4. Используйте HTTPS

### Docker

```bash
docker build -t mobile-store-back .
docker run -p 8080:8080 mobile-store-back
```
