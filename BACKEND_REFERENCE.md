# 🚀 Backend Reference Guide

> Быстрый справочник по бэкенду для разработки фронтенда

## 📁 Структура проекта

```
MobileStoreBack/
├── main.go                          # Точка входа
├── docker-compose.yml               # Docker конфигурация
├── init.sql                         # База данных + тестовые данные
├── .env                             # Переменные окружения
├── internal/
│   ├── config/                      # Конфигурация
│   ├── database/                    # Подключение к БД
│   ├── models/                      # Модели данных
│   ├── repository/                  # Слой данных
│   ├── services/                    # Бизнес-логика
│   ├── handlers/                    # HTTP обработчики
│   └── middleware/                  # Middleware
└── API_ENDPOINTS.md                 # Полная документация API
```

## 🗄️ База данных (12 таблиц)

### Основные таблицы:

- `users` - пользователи
- `categories` - категории товаров
- `products` - товары
- `product_variants` - варианты товаров (цвет, размер)
- `warehouses` - склады/филиалы
- `warehouse_stocks` - остатки товаров по складам

### Вспомогательные таблицы:

- `images` - изображения товаров
- `cart_items` - корзина
- `wishlist_items` - избранное
- `orders` - заказы
- `order_items` - элементы заказов
- `reviews` - отзывы

## 🔗 API Endpoints (48 штук)

### Публичные (14 endpoints):

```
GET  /health                                    # Health check
GET  /api/categories                           # Список категорий
GET  /api/categories/:id                       # Категория по ID
GET  /api/categories/:id/products              # Товары категории
GET  /api/products                             # Список товаров (с поиском и фильтрацией)
GET  /api/products/:slug                       # Товар по slug
GET  /api/products/:slug/reviews               # Отзывы товара
GET  /api/products/:slug/variants              # Варианты товара
GET  /api/warehouses                           # Список складов
GET  /api/warehouses/main                      # Главный склад
GET  /api/warehouses/:slug                     # Склад по slug
GET  /api/warehouses/city/:city                # Склады по городу
GET  /api/stocks/warehouse/:warehouse_slug     # Остатки по складу
GET  /api/stocks/variant/:sku                  # Остатки по варианту
GET  /api/images/product/:slug                 # Изображения товара
```

### Защищенные (15 endpoints):

```
# Пользователи
GET  /api/users/profile                        # Профиль
PUT  /api/users/profile                        # Обновить профиль

# Заказы
POST /api/orders                               # Создать заказ
GET  /api/orders                               # Мои заказы
GET  /api/orders/:id                           # Заказ по ID
PUT  /api/orders/:id                           # Обновить заказ

# Корзина
GET    /api/cart                               # Содержимое корзины
POST   /api/cart                               # Добавить в корзину
PUT    /api/cart/:id                           # Обновить элемент
DELETE /api/cart/:id                           # Удалить из корзины
DELETE /api/cart                               # Очистить корзину
GET    /api/cart/count                         # Количество товаров

# Избранное
GET    /api/wishlist                           # Избранные товары
POST   /api/wishlist                           # Добавить в избранное
DELETE /api/wishlist/:id                       # Удалить из избранного
DELETE /api/wishlist                           # Очистить избранное
GET    /api/wishlist/check/:product_id         # Проверить наличие

# Отзывы
POST   /api/reviews                            # Создать отзыв
GET    /api/reviews/my                         # Мои отзывы
PUT    /api/reviews/:id                        # Обновить отзыв
DELETE /api/reviews/:id                        # Удалить отзыв
POST   /api/reviews/:id/vote                   # Проголосовать
```

### Админские (25 endpoints):

```
# Пользователи
GET    /api/admin/users                        # Список пользователей
GET    /api/admin/users/:id                    # Пользователь по ID
PUT    /api/admin/users/:id                    # Обновить пользователя
DELETE /api/admin/users/:id                    # Удалить пользователя

# Каталог
POST   /api/admin/products                     # Создать товар
PUT    /api/admin/products/:id                 # Обновить товар
DELETE /api/admin/products/:id                 # Удалить товар
POST   /api/admin/product-variants             # Создать вариант
GET    /api/admin/product-variants/:id         # Вариант по ID
PUT    /api/admin/product-variants/:id         # Обновить вариант
DELETE /api/admin/product-variants/:id         # Удалить вариант
POST   /api/admin/categories                   # Создать категорию
GET    /api/admin/categories/:id               # Категория по ID
PUT    /api/admin/categories/:id               # Обновить категорию
DELETE /api/admin/categories/:id               # Удалить категорию

# Склады
POST   /api/admin/warehouses                   # Создать склад
GET    /api/admin/warehouses/:id               # Склад по ID
PUT    /api/admin/warehouses/:id               # Обновить склад
DELETE /api/admin/warehouses/:id               # Удалить склад

# Остатки
POST   /api/admin/warehouse-stocks             # Создать остаток
PUT    /api/admin/warehouse-stocks/:id         # Обновить остаток
DELETE /api/admin/warehouse-stocks/:id         # Удалить остаток

# Изображения
POST   /api/admin/images/product/:id           # Загрузить изображение
DELETE /api/admin/images/:id                   # Удалить изображение
PUT    /api/admin/images/:id/primary           # Установить главное

# Заказы
GET    /api/admin/orders                       # Все заказы
PUT    /api/admin/orders/:id/status            # Обновить статус

# Контент
GET    /api/admin/reviews                      # Все отзывы
PUT    /api/admin/reviews/:id/approve          # Модерация отзыва
```

## 🔐 Аутентификация

### JWT токены:

- Заголовок: `Authorization: Bearer <token>`
- Токен получается через `/api/auth/login`
- Токен нужен для защищенных и админских API

### Роли:

- **Публичные API** - без токена
- **Защищенные API** - нужен токен пользователя
- **Админские API** - нужен токен + роль admin

## 📊 Модели данных

### Product (товар):

```json
{
  "id": "uuid",
  "name": "Чехол Apple iPhone 15 Pro",
  "slug": "chehol-apple-iphone-15-pro",
  "description": "Описание товара",
  "base_price": 4990.0,
  "sku": "APPLE-CASE-IP15P",
  "is_active": true,
  "brand": "Apple",
  "model": "iPhone 15 Pro",
  "material": "Силикон",
  "category_id": "uuid",
  "tags": ["чехол", "apple", "iphone"],
  "view_count": 0,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### ProductVariant (вариант товара):

```json
{
  "id": "uuid",
  "product_id": "uuid",
  "sku": "APPLE-CASE-IP15P-BLACK",
  "name": "Черный, L",
  "color": "Черный",
  "size": "L",
  "price": 4990.0,
  "is_active": true
}
```

### Category (категория):

```json
{
  "id": "uuid",
  "name": "Чехлы для телефонов",
  "slug": "chehly-dlya-telefonov",
  "description": "Описание категории",
  "image_url": "https://example.com/category.jpg"
}
```

### Warehouse (склад):

```json
{
  "id": "uuid",
  "name": "Главный склад",
  "slug": "main-warehouse",
  "address": "ул. Промышленная, 15",
  "city": "Москва",
  "phone": "+7 (495) 123-45-67",
  "email": "main@shop.com",
  "is_active": true,
  "is_main": true,
  "manager_name": "Иванов Иван Иванович"
}
```

### WarehouseStock (остаток):

```json
{
  "id": "uuid",
  "warehouse_id": "uuid",
  "product_variant_id": "uuid",
  "stock": 50,
  "reserved_stock": 5,
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Image (изображение):

```json
{
  "id": "uuid",
  "product_id": "uuid",
  "cloudinary_public_id": "product-image-123",
  "url": "https://res.cloudinary.com/image.jpg",
  "is_primary": true,
  "created_at": "2024-01-15T10:30:00Z"
}
```

## 🚀 Запуск проекта

### Docker (рекомендуется):

```bash
docker-compose up -d
```

### Локально:

```bash
# 1. Установить PostgreSQL и Redis
# 2. Создать БД и выполнить init.sql
# 3. Настроить .env
go run main.go
```

### Проверка:

```bash
curl http://localhost:8080/health
```

## 🔧 Полезные команды

### Проверка API:

```bash
# Health check
curl http://localhost:8080/health

# Категории
curl http://localhost:8080/api/categories

# Товары
curl http://localhost:8080/api/products

# Поиск
curl "http://localhost:8080/api/products/search?q=iphone"
```

### Тестовые данные:

- 4 категории товаров
- 20+ товаров с вариантами
- 4 склада с остатками
- Тестовые изображения
- Пользователи и заказы

## 📝 Важные особенности

### URL-friendly:

- Публичные API используют slug вместо UUID
- `/api/products/chehol-apple-iphone-15-pro` вместо `/api/products/uuid`

### Мультискладовая система:

- Товары могут быть на разных складах
- Остатки управляются через `warehouse_stocks`
- Резервирование товаров при заказе

### Фильтрация товаров:

- `?category_id=uuid` - по категории
- `?brand=Apple` - по бренду
- `?min_price=1000&max_price=5000` - по цене
- `?is_active=true` - только активные
- `?limit=20&offset=0` - пагинация

### Поиск:

- `?q=iphone` - поиск по названию и описанию
- `?category=chehly` - поиск в категории
- `?brand=Apple` - поиск по бренду

## 🎯 Для фронтенда

### Начните с:

1. **Health check** - проверка сервера
2. **Категории** - навигация
3. **Товары** - каталог
4. **Авторизация** - вход/регистрация
5. **Корзина** - покупки

### Критичные API:

```bash
GET /health                           # Проверка сервера
GET /api/categories                   # Меню категорий
GET /api/products                     # Каталог товаров
GET /api/products/:slug               # Страница товара
POST /api/auth/login                  # Авторизация
GET /api/cart                         # Корзина
POST /api/cart                        # Добавить в корзину
```

### Примеры запросов:

```bash
# Получить товар
GET /api/products/chehol-apple-iphone-15-pro

# Получить варианты товара
GET /api/products/chehol-apple-iphone-15-pro/variants

# Получить изображения товара
GET /api/images/product/chehol-apple-iphone-15-pro

# Поиск товаров
GET /api/products/search?q=iphone&brand=Apple

# Остатки товара
GET /api/stocks/variant/APPLE-CASE-IP15P-BLACK
```

## 📞 Поддержка

- **API документация**: `API_ENDPOINTS.md`
- **База данных**: `init.sql`
- **Конфигурация**: `.env`
- **Docker**: `docker-compose.yml`

---

_Обновлено: $(date)_
_Версия API: 1.1_
_Всего endpoints: 48_
