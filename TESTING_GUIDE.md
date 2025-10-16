# 🧪 Руководство по тестированию API

## 📋 Обзор

Это руководство содержит все необходимые запросы для тестирования API магазина мобильных аксессуаров через Postman.

## 🔧 Настройка Postman

### Переменные окружения:

- `{{base_url}}` = `http://localhost:8080/api`
- `{{token}}` = токен пользователя (получается после логина)
- `{{admin_token}}` = токен администратора

### Headers по умолчанию:

- `Content-Type: application/json` (для POST/PUT запросов)
- `Authorization: Bearer {{token}}` (для защищенных endpoints)

## 👤 Тестовые пользователи

Из скрипта `init.sql`:

### Администратор:

- **Email:** `admin@shop.com`
- **Password:** `password` (хэшированный пароль в БД)
- **Роль:** Администратор

### Обычный пользователь:

- **Email:** `user@shop.com`
- **Password:** `password` (хэшированный пароль в БД)
- **Роль:** Пользователь

## 🔐 1. Аутентификация

### Регистрация нового пользователя

```http
POST {{base_url}}/auth/register
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "gender": "male"
}
```

**Примечание:** Поля `phone` и `gender` являются необязательными. Для `gender` допустимы значения: `male`, `female` или можно не указывать (будет NULL).

### Логин пользователя

```http
POST {{base_url}}/auth/login
Content-Type: application/json

{
  "email": "user@shop.com",
  "password": "password"
}
```

### Логин администратора

```http
POST {{base_url}}/auth/login
Content-Type: application/json

{
  "email": "admin@shop.com",
  "password": "password"
}
```

**Сохраните токен из ответа в переменную `{{token}}` или `{{admin_token}}`**

## 👤 2. Управление пользователями

### Получить профиль

```http
GET {{base_url}}/profile
Authorization: Bearer {{token}}
```

### Обновить профиль

```http
PUT {{base_url}}/profile
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "date_of_birth": "1990-01-01",
  "gender": "male"
}
```

### Обновить адрес пользователя (встроенный)

```http
PUT {{base_url}}/profile
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "address_title": "Дом",
  "address_first_name": "Иван",
  "address_last_name": "Петров",
  "address_company": "ООО Компания",
  "address_street": "ул. Пушкина, д. 10, кв. 5",
  "address_city": "Москва",
  "address_state": "Московская область",
  "address_postal_code": "123456",
  "address_country": "Россия",
  "address_phone": "+7 (999) 123-45-67"
}
```

**Примечание:** Адрес пользователя теперь встроен в профиль. Все поля адреса являются необязательными.

### Получить всех пользователей (только админ)

```http
GET {{base_url}}/admin/users?limit=10&offset=0
Authorization: Bearer {{admin_token}}
```

### Получить пользователя по ID (только админ)

```http
GET {{base_url}}/admin/users/{user_id}
Authorization: Bearer {{admin_token}}
```

### Обновить пользователя (только админ)

```http
PUT {{base_url}}/admin/users/{user_id}
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith",
  "is_active": true,
  "is_admin": false
}
```

### Удалить пользователя (только админ)

```http
DELETE {{base_url}}/admin/users/{user_id}
Authorization: Bearer {{admin_token}}
```

## 📦 3. Управление товарами

### Получить список товаров

```http
GET {{base_url}}/products?limit=10&offset=0
```

### Получить товар по ID

```http
GET {{base_url}}/products/{product_id}
```

### Поиск товаров

```http
GET {{base_url}}/products/search?q=iPhone&limit=10&offset=0
```

### Получить товары по категории

```http
GET {{base_url}}/products/category/{category_id}?limit=10&offset=0
```

### Создать товар (только админ)

```http
POST {{base_url}}/admin/products
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "name": "iPhone 15 Pro Case",
  "description": "Premium case for iPhone 15 Pro with drop protection and wireless charging support",
  "short_description": "Premium case для iPhone 15 Pro",
  "price": 4990.00,
  "compare_price": 5990.00,
  "sku": "IPH15P-CASE-001",
  "stock": 25,
  "is_active": true,
  "is_featured": true,
  "is_new": true,
  "weight": 0.05,
  "dimensions": "15.5x7.8x1.2 cm",
  "brand": "Apple",
  "model": "Silicone Case",
  "color": "Синий",
  "material": "Силикон",
  "category_id": "category-uuid-here",
  "tags": ["iPhone 15 Pro", "официальный", "силикон"],
  "meta_title": "Чехол Apple для iPhone 15 Pro - Синий",
  "meta_description": "Официальный силиконовый чехол Apple для iPhone 15 Pro с защитой от падений"
}
```

**Примечание:**

- `brand` теперь строка, а не ссылка на отдельную таблицу
- Убраны поля `barcode`, `min_stock`, `attributes`, `specifications`, `compatible_with`
- `tags` - массив строк
- Цены в рублях (DECIMAL)

### Обновить товар (только админ)

```http
PUT {{base_url}}/admin/products/{product_id}
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "name": "iPhone 15 Pro Case - Updated",
  "price": 899.99,
  "is_featured": true
}
```

### Удалить товар (только админ)

```http
DELETE {{base_url}}/admin/products/{product_id}
Authorization: Bearer {{admin_token}}
```

## 🛒 4. Корзина

### Получить корзину

```http
GET {{base_url}}/cart
Authorization: Bearer {{token}}
```

### Добавить товар в корзину

```http
POST {{base_url}}/cart/add
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "product_id": "product-uuid-here",
  "quantity": 2
}
```

### Обновить количество товара в корзине

```http
PUT {{base_url}}/cart/items/{item_id}
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "quantity": 3
}
```

### Удалить товар из корзины

```http
DELETE {{base_url}}/cart/items/{item_id}
Authorization: Bearer {{token}}
```

### Очистить корзину

```http
DELETE {{base_url}}/cart/clear
Authorization: Bearer {{token}}
```

### Получить количество товаров в корзине

```http
GET {{base_url}}/cart/count
Authorization: Bearer {{token}}
```

## ❤️ 5. Избранное

### Получить список избранного

```http
GET {{base_url}}/wishlist?limit=10&offset=0
Authorization: Bearer {{token}}
```

### Добавить товар в избранное

```http
POST {{base_url}}/wishlist/add
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "product_id": "product-uuid-here"
}
```

### Удалить товар из избранного

```http
DELETE {{base_url}}/wishlist/items/{item_id}
Authorization: Bearer {{token}}
```

### Проверить, есть ли товар в избранном

```http
GET {{base_url}}/wishlist/check/{product_id}
Authorization: Bearer {{token}}
```

### Очистить избранное

```http
DELETE {{base_url}}/wishlist/clear
Authorization: Bearer {{token}}
```

## 📝 6. Отзывы

### Получить отзывы товара

```http
GET {{base_url}}/products/{product_id}/reviews?limit=10&offset=0
```

### Создать отзыв

```http
POST {{base_url}}/reviews
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "product_id": "product-uuid-here",
  "order_id": "order-uuid-here",
  "rating": 5,
  "title": "Great product!",
  "comment": "Really happy with this purchase"
}
```

### Обновить отзыв

```http
PUT {{base_url}}/reviews/{review_id}
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "rating": 4,
  "title": "Updated review",
  "comment": "Changed my mind"
}
```

### Удалить отзыв

```http
DELETE {{base_url}}/reviews/{review_id}
Authorization: Bearer {{token}}
```

### Оценить полезность отзыва

```http
POST {{base_url}}/reviews/{review_id}/vote
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "helpful": true
}
```

### Получить мои отзывы

```http
GET {{base_url}}/reviews/my?limit=10&offset=0
Authorization: Bearer {{token}}
```

### Получить все отзывы (только админ)

```http
GET {{base_url}}/admin/reviews?limit=10&offset=0
Authorization: Bearer {{admin_token}}
```

### Одобрить отзыв (только админ)

```http
PUT {{base_url}}/admin/reviews/{review_id}/approve
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "approved": true
}
```

## 🎫 7. Промокоды

### Получить все промокоды

```http
GET {{base_url}}/coupons?limit=10&offset=0
```

### Получить промокод по ID

```http
GET {{base_url}}/coupons/{coupon_id}
```

### Валидировать промокод

```http
POST {{base_url}}/coupons/validate
Content-Type: application/json

{
  "code": "SAVE20",
  "user_id": "user-uuid-here",
  "order_amount": 100.00
}
```

### Создать промокод (только админ)

```http
POST {{base_url}}/admin/coupons
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "code": "SAVE20",
  "name": "20% Off",
  "description": "Get 20% off your order",
  "type": "percentage",
  "value": 20.0,
  "minimum_amount": 50.0,
  "maximum_discount": 100.0,
  "usage_limit": 100,
  "starts_at": "2024-01-01T00:00:00Z",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

### Обновить промокод (только админ)

```http
PUT {{base_url}}/admin/coupons/{coupon_id}
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "name": "Updated Coupon",
  "is_active": true
}
```

### Удалить промокод (только админ)

```http
DELETE {{base_url}}/admin/coupons/{coupon_id}
Authorization: Bearer {{admin_token}}
```

### Получить использование промокода (только админ)

```http
GET {{base_url}}/admin/coupons/{coupon_id}/usage
Authorization: Bearer {{admin_token}}
```

## 📦 8. Заказы

### Создать заказ

```http
POST {{base_url}}/orders
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "items": [
    {
      "product_id": "product-uuid-here",
      "quantity": 2
    }
  ],
  "shipping_method": "delivery",
  "shipping_address": "ул. Пушкина, д. 10, кв. 5, Москва, 123456",
  "pickup_point": "Пункт самовывоза: ул. Ленина, д. 15, Москва",
  "payment_method": "card",
  "customer_notes": "Please deliver after 5 PM",
  "coupon_code": "SAVE20"
}
```

### Получить заказы пользователя

```http
GET {{base_url}}/orders?limit=10&offset=0
Authorization: Bearer {{token}}
```

### Получить заказ по ID

```http
GET {{base_url}}/orders/{order_id}
Authorization: Bearer {{token}}
```

### Обновить заказ

```http
PUT {{base_url}}/orders/{order_id}
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "status": "shipped",
  "payment_status": "paid",
  "tracking_number": "TRK123456789"
}
```

### Получить все заказы (только админ)

```http
GET {{base_url}}/admin/orders?limit=10&offset=0
Authorization: Bearer {{admin_token}}
```

### Обновить статус заказа (только админ)

```http
PUT {{base_url}}/admin/orders/{order_id}/status
Authorization: Bearer {{admin_token}}
Content-Type: application/json

{
  "status": "delivered",
  "tracking_number": "TRK123456789"
}
```

## 🧪 Последовательность тестирования

### 1. Базовая настройка

1. Запустите сервер
2. Выполните логин администратора
3. Сохраните токен в `{{admin_token}}`

### 2. Тестирование товаров

1. Получите список категорий
2. Создайте несколько товаров
3. Протестируйте поиск и фильтрацию

### 3. Тестирование пользователей

1. Зарегистрируйте нового пользователя
2. Выполните логин
3. Сохраните токен в `{{token}}`

### 4. Тестирование корзины и избранного

1. Добавьте товары в корзину
2. Добавьте товары в избранное
3. Протестируйте обновление и удаление

### 5. Тестирование заказов

1. Обновите адрес в профиле (если нужно)
2. Создайте заказ
3. Протестируйте обновление статуса

### 6. Тестирование отзывов

1. Создайте отзыв
2. Протестируйте голосование
3. Протестируйте модерацию

## 📊 Ожидаемые ответы

### Успешные ответы:

- **200 OK** - успешное выполнение
- **201 Created** - создание ресурса
- **204 No Content** - успешное удаление

### Ошибки:

- **400 Bad Request** - неверные данные
- **401 Unauthorized** - не авторизован
- **403 Forbidden** - нет прав доступа
- **404 Not Found** - ресурс не найден
- **500 Internal Server Error** - ошибка сервера

## 🔍 Полезные советы

1. **Всегда проверяйте токены** - они могут истекать
2. **Используйте реальные UUID** - скопируйте их из ответов
3. **Тестируйте валидацию** - отправляйте неверные данные
4. **Проверяйте права доступа** - тестируйте админские функции
5. **Тестируйте пагинацию** - используйте параметры limit/offset

## 🚀 Готовые тестовые данные

В базе данных уже есть:

- **2 пользователя** (admin@shop.com, user@shop.com)
  - Пароль для обоих: `password`
  - admin@shop.com - администратор
  - user@shop.com - обычный пользователь
- **21 категория** товаров (чехлы, зарядки, наушники, защитные стекла и др.)
- **10 товаров** с изображениями:
  - Чехлы Apple и Spigen для iPhone
  - Зарядные устройства Apple и Anker
  - Наушники Apple AirPods Pro 2 и JBL
  - Защитные стекла Belkin
  - USB-C кабели UGREEN
  - PowerBank Anker
  - Подставки Baseus
- **Способы доставки** встроены в заказы (delivery, pickup)

**Важно:**

- Поле `gender` в пользователях может быть NULL (исправлено ограничение)
- Адреса пользователей встроены в профиль
- Бренды хранятся как строки в таблице товаров

Используйте эти данные для быстрого старта тестирования!
