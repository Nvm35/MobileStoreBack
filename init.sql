-- =============================================
-- ПОЛНАЯ СХЕМА БАЗЫ ДАННЫХ ДЛЯ МАГАЗИНА МОБИЛЬНЫХ АКСЕССУАРОВ
-- =============================================

-- Создание расширений
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =============================================
-- ОСНОВНЫЕ ТАБЛИЦЫ (в правильном порядке зависимостей)
-- =============================================

-- 1. Создание таблицы категорий (самореференс)
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    slug VARCHAR(255) NOT NULL UNIQUE,
    image_url TEXT, -- URL изображения категории
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    is_admin BOOLEAN DEFAULT false,
    email_verified BOOLEAN DEFAULT false,
    email_verification_token VARCHAR(255),
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP,
    last_login TIMESTAMP,
    -- Адрес пользователя (вместо отдельной таблицы addresses)
    address_title VARCHAR(255),
    address_first_name VARCHAR(255),
    address_last_name VARCHAR(255),
    address_street TEXT,
    address_city VARCHAR(255),
    address_state VARCHAR(255),
    address_postal_code VARCHAR(20),
    address_country VARCHAR(255),
    address_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 3. Создание таблицы продуктов (зависит от categories)
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE, -- URL-friendly slug для товара
    description TEXT,
    base_price DECIMAL(10,2) NOT NULL, -- базовая цена товара
    sku VARCHAR(255) NOT NULL UNIQUE,
    stock INTEGER NOT NULL DEFAULT 0, -- общий остаток товара
    is_active BOOLEAN DEFAULT true,
    brand VARCHAR(255) NOT NULL, -- бренд как строка
    model VARCHAR(255),
    material VARCHAR(255),
    category_id UUID NOT NULL REFERENCES categories(id),
    tags TEXT[], -- массив тегов
    view_count INTEGER DEFAULT 0, -- количество просмотров
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Создание таблицы вариантов товаров (зависит от products)
CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(255) NOT NULL UNIQUE, -- уникальный SKU для варианта
    name VARCHAR(255) NOT NULL, -- название варианта (например, "Красный, L")
    color VARCHAR(100), -- цвет варианта
    size VARCHAR(50), -- размер варианта
    price DECIMAL(10,2) NOT NULL, -- цена варианта (может отличаться от базовой)
    stock INTEGER NOT NULL DEFAULT 0, -- остаток варианта
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Создание таблицы изображений (зависит от products)
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    cloudinary_public_id VARCHAR(255) NOT NULL, -- ID изображения в Cloudinary
    url TEXT NOT NULL, -- полный URL изображения
    is_primary BOOLEAN DEFAULT false, -- главное изображение
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6. Создание таблицы корзины (зависит от users, products)
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(255), -- для неавторизованных пользователей
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price DECIMAL(10,2) NOT NULL, -- цена на момент добавления
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP, -- срок действия корзины
    UNIQUE(user_id, product_id), -- один товар на пользователя
    UNIQUE(session_id, product_id) -- один товар на сессию
);

-- 7. Создание таблицы избранного (зависит от users, products)
CREATE TABLE IF NOT EXISTS wishlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id) -- один товар на пользователя
);

-- 8. Создание таблицы заказов (зависит от users)
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    order_number VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL, -- общая сумма заказа
    payment_method VARCHAR(50) NOT NULL,
    payment_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- Способ доставки
    shipping_method VARCHAR(50) NOT NULL DEFAULT 'delivery', -- 'delivery', 'pickup'
    -- Адрес доставки (если нужен другой адрес, чем у пользователя)
    shipping_address TEXT, -- полный адрес доставки в текстовом виде
    -- Пункт самовывоза (если выбран pickup)
    pickup_point TEXT, -- название и адрес пункта самовывоза
    tracking_number VARCHAR(255),
    notes TEXT,
    customer_notes TEXT, -- заметки клиента
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 9. Создание таблицы элементов заказа (зависит от orders, products)
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    product_variant_id UUID REFERENCES product_variants(id), -- ссылка на вариант товара
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL, -- цена на момент заказа
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. Создание таблицы отзывов (зависит от users, products, orders)
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    order_id UUID REFERENCES orders(id), -- связь с заказом для проверки покупки
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    comment TEXT,
    is_verified BOOLEAN DEFAULT false, -- подтвержденная покупка
    is_approved BOOLEAN DEFAULT true, -- модерация отзыва
    helpful_count INTEGER DEFAULT 0, -- количество лайков
    unhelpful_count INTEGER DEFAULT 0, -- количество дизлайков
    -- Оценки полезности (вместо отдельной таблицы)
    helpful_votes JSONB DEFAULT '[]', -- [{"user_id": "uuid", "helpful": true, "created_at": "2024-01-01T10:00:00Z"}]
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id) -- один отзыв на товар от пользователя
);


-- =============================================
-- ИНДЕКСЫ ДЛЯ ОПТИМИЗАЦИИ
-- =============================================

-- Основные индексы
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_product_variants_sku ON product_variants(sku);
CREATE INDEX IF NOT EXISTS idx_product_variants_active ON product_variants(is_active);
CREATE INDEX IF NOT EXISTS idx_images_product_id ON images(product_id);
CREATE INDEX IF NOT EXISTS idx_images_primary ON images(is_primary);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);
CREATE INDEX IF NOT EXISTS idx_order_items_variant_id ON order_items(product_variant_id);

-- Индексы для корзины
CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_session_id ON cart_items(session_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items(product_id);

-- Индексы для избранного
CREATE INDEX IF NOT EXISTS idx_wishlist_items_user_id ON wishlist_items(user_id);
CREATE INDEX IF NOT EXISTS idx_wishlist_items_product_id ON wishlist_items(product_id);

-- Индексы для отзывов
CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews(product_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
CREATE INDEX IF NOT EXISTS idx_reviews_approved ON reviews(is_approved);


-- Индексы для брендов (теперь в products)
CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand);

CREATE INDEX IF NOT EXISTS idx_reviews_helpful_votes_gin ON reviews USING gin(helpful_votes);

-- Индексы для поиска по брендам в товарах
CREATE INDEX IF NOT EXISTS idx_products_brand_category ON products(brand, category_id);

-- Полнотекстовый поиск по товарам
CREATE INDEX IF NOT EXISTS idx_products_search ON products USING gin(to_tsvector('russian', name || ' ' || COALESCE(description, '')));

-- Индексы для JSON полей удалены - поля больше не существуют

-- =============================================
-- НАЧАЛЬНЫЕ ДАННЫЕ
-- =============================================

-- Бренды теперь хранятся как строки в таблице products

-- Вставка категорий для всех типов товаров
INSERT INTO categories (name, description, slug) VALUES 
-- Мобильные аксессуары
('Чехлы', 'Чехлы для мобильных телефонов', 'cases'),
('Зарядные устройства', 'Зарядные устройства и кабели', 'chargers'),
('Наушники', 'Наушники и гарнитуры', 'headphones'),
('Защитные стекла', 'Защитные стекла и пленки', 'screen-protectors'),
('Беспроводные зарядки', 'Беспроводные зарядные устройства', 'wireless-chargers'),
('Кабели и адаптеры', 'USB кабели и адаптеры', 'cables-adapters'),
('Портативные аккумуляторы', 'Power Bank и внешние батареи', 'power-banks'),
('Подставки и держатели', 'Подставки для телефонов', 'stands-holders'),
('Стилусы', 'Стилусы для сенсорных экранов', 'styluses'),
-- Компьютерные аксессуары
('Клавиатуры', 'Механические и мембранные клавиатуры', 'keyboards'),
('Мыши', 'Компьютерные мыши и трекболы', 'mice'),
('Коврики для мыши', 'Игровые и офисные коврики', 'mouse-pads'),
('Веб-камеры', 'Веб-камеры для стриминга и конференций', 'webcams'),
('Микрофоны', 'Микрофоны для стриминга и записи', 'microphones'),
-- Фото и видео
('Штативы', 'Штативы для камер и телефонов', 'tripods'),
('Освещение', 'Световое оборудование для фото и видео', 'lighting'),
('Карты памяти', 'SD карты и флеш-накопители', 'memory-cards'),
-- Автомобильные аксессуары
('Автомобильные держатели', 'Держатели для телефонов в авто', 'car-holders'),
('Автомобильные зарядки', 'Зарядные устройства для авто', 'car-chargers'),
-- Игровые аксессуары
('Игровые контроллеры', 'Геймпады и контроллеры', 'game-controllers'),
('Игровые кресла', 'Кресла для геймеров', 'gaming-chairs')
ON CONFLICT (slug) DO NOTHING;

-- Способы доставки теперь встроены в orders

-- =============================================
-- ТРИГГЕРЫ ДЛЯ АВТОМАТИЧЕСКОГО ОБНОВЛЕНИЯ
-- =============================================

-- Создание триггеров для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Применение триггеров к таблицам
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Триггер для addresses удален - таблица больше не существует
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_product_variants_updated_at BEFORE UPDATE ON product_variants FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_cart_items_updated_at BEFORE UPDATE ON cart_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Триггер для coupons удален - таблица больше не существует
-- Триггер для shipping_methods удален - таблица больше не существует
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- ТЕСТОВЫЕ ДАННЫЕ
-- =============================================

-- Создание тестового пользователя
INSERT INTO users (email, password, first_name, last_name, phone, is_active, is_admin) VALUES 
('admin@shop.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Админ', 'Админов', '+7 (999) 123-45-67', true, true),
('user@shop.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Иван', 'Петров', '+7 (999) 765-43-21', true, false)
ON CONFLICT (email) DO NOTHING;

-- Создание тестовых товаров с разными брендами
INSERT INTO products (name, slug, description, base_price, sku, stock, is_active, brand, model, material, category_id, tags) VALUES 
-- Чехлы для iPhone
('Чехол Apple Silicone Case для iPhone 15 Pro', 'apple-silicone-case-iphone-15-pro', 'Официальный силиконовый чехол Apple с мягкой внутренней поверхностью и внешней поверхностью из силикона.', 4990.00, 'APPLE-CASE-IP15P', 50, true, 'Apple', 'Silicone Case', 'Силикон', (SELECT id FROM categories WHERE slug = 'cases'), ARRAY['iPhone 15 Pro', 'официальный', 'силикон']),

('Чехол Spigen Ultra Hybrid для iPhone 15', 'spigen-ultra-hybrid-iphone-15', 'Прозрачный чехол Spigen с защитой от падений и поддержкой беспроводной зарядки.', 1990.00, 'SPIGEN-UH-IP15', 30, true, 'Spigen', 'Ultra Hybrid', 'TPU + Поликарбонат', (SELECT id FROM categories WHERE slug = 'cases'), ARRAY['iPhone 15', 'прозрачный', 'защита']),

-- Зарядные устройства
('Зарядное устройство Apple 20W USB-C', 'apple-20w-usb-c-charger', 'Официальное зарядное устройство Apple мощностью 20W с разъемом USB-C.', 2990.00, 'APPLE-20W-USB-C', 25, true, 'Apple', '20W USB-C Power Adapter', 'Пластик', (SELECT id FROM categories WHERE slug = 'chargers'), ARRAY['iPhone', 'iPad', 'быстрая зарядка']),

('Беспроводная зарядка Anker PowerWave 7.5W', 'anker-powervave-7-5w-wireless-charger', 'Беспроводная зарядная станция Anker с поддержкой быстрой зарядки до 7.5W.', 2490.00, 'ANKER-PW-7.5W', 15, true, 'Anker', 'PowerWave 7.5W', 'Пластик + Силикон', (SELECT id FROM categories WHERE slug = 'wireless-chargers'), ARRAY['беспроводная зарядка', 'Qi', 'быстрая зарядка']),

-- Наушники
('Наушники Apple AirPods Pro 2', 'apple-airpods-pro-2', 'Беспроводные наушники Apple AirPods Pro 2-го поколения с активным шумоподавлением.', 19990.00, 'APPLE-AIRPODS-PRO-2', 10, true, 'Apple', 'AirPods Pro 2', 'Пластик', (SELECT id FROM categories WHERE slug = 'headphones'), ARRAY['беспроводные', 'шумоподавление', 'пространственное аудио']),

('Наушники JBL Tune 760NC', 'jbl-tune-760nc', 'Беспроводные наушники JBL с активным шумоподавлением и 50-часовым временем работы.', 8990.00, 'JBL-TUNE-760NC', 20, true, 'JBL', 'Tune 760NC', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'headphones'), ARRAY['беспроводные', 'шумоподавление', 'долгая работа']),

-- Защитные стекла
('Защитное стекло Belkin ScreenForce для iPhone 15', 'belkin-screenforce-iphone-15', 'Защитное стекло Belkin с технологией ScreenForce для iPhone 15.', 1990.00, 'BELKIN-SF-IP15', 100, true, 'Belkin', 'ScreenForce', 'Закаленное стекло', (SELECT id FROM categories WHERE slug = 'screen-protectors'), ARRAY['iPhone 15', 'защитное стекло', '9H']),

-- Кабели
('USB-C кабель UGREEN 1м', 'ugreen-usb-c-cable-1m', 'USB-C кабель UGREEN длиной 1 метр с поддержкой быстрой зарядки и передачи данных.', 590.00, 'UGREEN-USB-C-1M', 200, true, 'UGREEN', 'USB-C Cable', 'Медь + Пластик', (SELECT id FROM categories WHERE slug = 'cables-adapters'), ARRAY['USB-C', 'быстрая зарядка', 'передача данных']),

-- Power Bank
('Портативный аккумулятор Anker PowerCore 10000', 'anker-powercore-10000', 'Портативный аккумулятор Anker PowerCore емкостью 10000 мАч с быстрой зарядкой.', 2990.00, 'ANKER-PC-10000', 40, true, 'Anker', 'PowerCore 10000', 'Пластик', (SELECT id FROM categories WHERE slug = 'power-banks'), ARRAY['10000 мАч', 'быстрая зарядка', 'компактный']),

-- Подставки
('Подставка для телефона Baseus Metal Stand', 'baseus-metal-stand-phone', 'Металлическая подставка Baseus для телефона с регулируемым углом наклона.', 1290.00, 'BASEUS-MS-PHONE', 60, true, 'Baseus', 'Metal Stand', 'Алюминий', (SELECT id FROM categories WHERE slug = 'stands-holders'), ARRAY['подставка', 'металлическая', 'регулируемая'])
ON CONFLICT (sku) DO NOTHING;

-- Создание тестовых вариантов товаров
INSERT INTO product_variants (product_id, sku, name, color, size, price, stock, is_active) VALUES 
-- Варианты чехла Apple
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P'), 'APPLE-CASE-IP15P-BLUE', 'Apple Silicone Case - Синий', 'Синий', NULL, 4990.00, 15, true),
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P'), 'APPLE-CASE-IP15P-BLACK', 'Apple Silicone Case - Черный', 'Черный', NULL, 4990.00, 20, true),
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P'), 'APPLE-CASE-IP15P-WHITE', 'Apple Silicone Case - Белый', 'Белый', NULL, 4990.00, 10, true),

-- Варианты чехла Spigen
((SELECT id FROM products WHERE sku = 'SPIGEN-UH-IP15'), 'SPIGEN-UH-IP15-CLEAR', 'Spigen Ultra Hybrid - Прозрачный', 'Прозрачный', NULL, 1990.00, 30, true),
((SELECT id FROM products WHERE sku = 'SPIGEN-UH-IP15'), 'SPIGEN-UH-IP15-BLACK', 'Spigen Ultra Hybrid - Черный', 'Черный', NULL, 1990.00, 25, true),

-- Варианты наушников JBL (размеры)
((SELECT id FROM products WHERE sku = 'JBL-TUNE-760NC'), 'JBL-TUNE-760NC-BLACK', 'JBL Tune 760NC - Черный', 'Черный', 'One Size', 8990.00, 15, true),
((SELECT id FROM products WHERE sku = 'JBL-TUNE-760NC'), 'JBL-TUNE-760NC-BLUE', 'JBL Tune 760NC - Синий', 'Синий', 'One Size', 8990.00, 10, true),

-- Варианты кабеля (длины)
((SELECT id FROM products WHERE sku = 'UGREEN-USB-C-1M'), 'UGREEN-USB-C-1M', 'UGREEN USB-C кабель 1м', 'Черный', '1м', 590.00, 50, true),
((SELECT id FROM products WHERE sku = 'UGREEN-USB-C-1M'), 'UGREEN-USB-C-2M', 'UGREEN USB-C кабель 2м', 'Черный', '2м', 790.00, 30, true),
((SELECT id FROM products WHERE sku = 'UGREEN-USB-C-1M'), 'UGREEN-USB-C-3M', 'UGREEN USB-C кабель 3м', 'Черный', '3м', 990.00, 20, true)
ON CONFLICT (sku) DO NOTHING;

-- Создание тестовых изображений для товаров
INSERT INTO images (product_id, cloudinary_public_id, url, is_primary) VALUES 
-- Изображения для чехла Apple
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P'), 'apple-case-ip15p-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-case-ip15p-1.jpg', true),
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P'), 'apple-case-ip15p-2', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-case-ip15p-2.jpg', false),

-- Изображения для чехла Spigen
((SELECT id FROM products WHERE sku = 'SPIGEN-UH-IP15'), 'spigen-uh-ip15-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/spigen-uh-ip15-1.jpg', true),

-- Изображения для зарядки Apple
((SELECT id FROM products WHERE sku = 'APPLE-20W-USB-C'), 'apple-20w-usb-c-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-20w-usb-c-1.jpg', true),

-- Изображения для беспроводной зарядки Anker
((SELECT id FROM products WHERE sku = 'ANKER-PW-7.5W'), 'anker-pw-7.5w-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/anker-pw-7.5w-1.jpg', true),

-- Изображения для AirPods Pro
((SELECT id FROM products WHERE sku = 'APPLE-AIRPODS-PRO-2'), 'apple-airpods-pro-2-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-airpods-pro-2-1.jpg', true),
((SELECT id FROM products WHERE sku = 'APPLE-AIRPODS-PRO-2'), 'apple-airpods-pro-2-2', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-airpods-pro-2-2.jpg', false),

-- Изображения для JBL наушников
((SELECT id FROM products WHERE sku = 'JBL-TUNE-760NC'), 'jbl-tune-760nc-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/jbl-tune-760nc-1.jpg', true),

-- Изображения для защитного стекла Belkin
((SELECT id FROM products WHERE sku = 'BELKIN-SF-IP15'), 'belkin-sf-ip15-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/belkin-sf-ip15-1.jpg', true),

-- Изображения для USB-C кабеля
((SELECT id FROM products WHERE sku = 'UGREEN-USB-C-1M'), 'ugreen-usb-c-1m-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/ugreen-usb-c-1m-1.jpg', true),

-- Изображения для PowerBank
((SELECT id FROM products WHERE sku = 'ANKER-PC-10000'), 'anker-pc-10000-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/anker-pc-10000-1.jpg', true),

-- Изображения для подставки Baseus
((SELECT id FROM products WHERE sku = 'BASEUS-MS-PHONE'), 'baseus-ms-phone-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/baseus-ms-phone-1.jpg', true)
ON CONFLICT DO NOTHING;