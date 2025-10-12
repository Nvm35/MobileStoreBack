#!/bin/bash

# Скрипт для тестирования API

BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Тестирование Mobile Store API..."

# Тест 1: Проверка здоровья API
echo "1️⃣ Проверка доступности API..."
curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/products" || echo "❌ API недоступен"

# Тест 2: Регистрация пользователя
echo -e "\n2️⃣ Тестирование регистрации..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }')

echo "Ответ регистрации: $REGISTER_RESPONSE"

# Извлекаем токен из ответа
TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo "✅ Регистрация успешна, токен получен"
    
    # Тест 3: Получение профиля
    echo -e "\n3️⃣ Тестирование получения профиля..."
    PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/users/profile" \
      -H "Authorization: Bearer $TOKEN")
    
    echo "Профиль: $PROFILE_RESPONSE"
    
    # Тест 4: Получение списка продуктов
    echo -e "\n4️⃣ Тестирование получения продуктов..."
    PRODUCTS_RESPONSE=$(curl -s -X GET "$BASE_URL/products")
    
    echo "Продукты: $PRODUCTS_RESPONSE"
    
else
    echo "❌ Не удалось получить токен"
fi

echo -e "\n✅ Тестирование завершено!"
