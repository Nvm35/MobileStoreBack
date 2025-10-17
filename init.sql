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
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    image_url TEXT, -- URL изображения категории
    meta_title VARCHAR(255), -- для SEO
    meta_description TEXT, -- для SEO
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    gender VARCHAR(10) CHECK (gender IS NULL OR gender = '' OR gender IN ('male', 'female')),
    is_active BOOLEAN DEFAULT true,
    is_admin BOOLEAN DEFAULT false,
    email_verified BOOLEAN DEFAULT false,
    email_verification_token VARCHAR(255),
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP,
    last_login TIMESTAMP,
    -- Уведомления (вместо отдельной таблицы)
    notifications JSONB DEFAULT '[]', -- [{"type": "order", "title": "Заказ отправлен", "message": "Ваш заказ #12345 отправлен", "is_read": false, "created_at": "2024-01-01T10:00:00Z"}]
    -- Адрес пользователя (вместо отдельной таблицы addresses)
    address_title VARCHAR(255),
    address_first_name VARCHAR(255),
    address_last_name VARCHAR(255),
    address_company VARCHAR(255),
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

-- 5. Создание таблицы продуктов (зависит от categories)
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE, -- URL-friendly slug для товара
    description TEXT,
    short_description VARCHAR(500),
    price DECIMAL(10,2) NOT NULL,
    compare_price DECIMAL(10,2), -- цена для сравнения (зачеркнутая цена)
    sku VARCHAR(255) NOT NULL UNIQUE,
    stock INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false, -- товар в избранном
    is_new BOOLEAN DEFAULT false, -- новинка
    weight DECIMAL(8,2),
    dimensions VARCHAR(100), -- "10x5x2 cm"
    brand VARCHAR(255) NOT NULL, -- бренд как строка
    model VARCHAR(255),
    color VARCHAR(100),
    material VARCHAR(255),
    category_id UUID NOT NULL REFERENCES categories(id),
    tags TEXT[], -- массив тегов
    meta_title VARCHAR(255), -- для SEO
    meta_description TEXT, -- для SEO
    view_count INTEGER DEFAULT 0, -- количество просмотров
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 6. Создание таблицы изображений (зависит от products)
CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    cloudinary_public_id VARCHAR(255) NOT NULL, -- ID изображения в Cloudinary
    url TEXT NOT NULL, -- полный URL изображения
    alt VARCHAR(255) NOT NULL, -- описание для SEO
    is_primary BOOLEAN DEFAULT false, -- главное изображение
    sort_order INTEGER DEFAULT 0, -- порядок отображения
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 7. Создание таблицы корзины (зависит от users, products)
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

-- 8. Создание таблицы избранного (зависит от users, products)
CREATE TABLE IF NOT EXISTS wishlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id) -- один товар на пользователя
);

-- 9. Создание таблицы промокодов
CREATE TABLE IF NOT EXISTS coupons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('percentage', 'fixed')),
    value DECIMAL(10,2) NOT NULL, -- процент или фиксированная сумма
    minimum_amount DECIMAL(10,2) DEFAULT 0, -- минимальная сумма заказа
    maximum_discount DECIMAL(10,2), -- максимальная скидка
    usage_limit INTEGER, -- лимит использований
    used_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    starts_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. Создание таблицы заказов (зависит от users)
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    order_number VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    shipping_cost DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
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

-- 12. Создание таблицы элементов заказа (зависит от orders, products)
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 13. Создание таблицы отзывов (зависит от users, products, orders)
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

-- 14. Создание таблицы использования промокодов (зависит от coupons, users, orders)
CREATE TABLE IF NOT EXISTS coupon_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coupon_id UUID NOT NULL REFERENCES coupons(id),
    user_id UUID NOT NULL REFERENCES users(id),
    order_id UUID NOT NULL REFERENCES orders(id),
    discount_amount DECIMAL(10,2) NOT NULL,
    used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
CREATE INDEX IF NOT EXISTS idx_products_featured ON products(is_featured);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);

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

-- Индексы для промокодов
CREATE INDEX IF NOT EXISTS idx_coupons_code ON coupons(code);
CREATE INDEX IF NOT EXISTS idx_coupons_active ON coupons(is_active);
CREATE INDEX IF NOT EXISTS idx_coupon_usage_coupon_id ON coupon_usage(coupon_id);
CREATE INDEX IF NOT EXISTS idx_coupon_usage_user_id ON coupon_usage(user_id);

-- Индексы для брендов (теперь в products)
CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand);

-- Индексы для JSONB полей
CREATE INDEX IF NOT EXISTS idx_users_notifications_gin ON users USING gin(notifications);
CREATE INDEX IF NOT EXISTS idx_reviews_helpful_votes_gin ON reviews USING gin(helpful_votes);

-- Индексы для поиска по брендам в товарах
CREATE INDEX IF NOT EXISTS idx_products_brand_category ON products(brand, category_id);

-- Полнотекстовый поиск по товарам
CREATE INDEX IF NOT EXISTS idx_products_search ON products USING gin(to_tsvector('russian', name || ' ' || COALESCE(description, '') || ' ' || COALESCE(short_description, '')));

-- Индексы для JSON полей удалены - поля больше не существуют

-- =============================================
-- НАЧАЛЬНЫЕ ДАННЫЕ
-- =============================================

-- Бренды теперь хранятся как строки в таблице products

-- Вставка категорий для всех типов товаров
INSERT INTO categories (name, description, slug, sort_order) VALUES 
-- Мобильные аксессуары
('Чехлы', 'Чехлы для мобильных телефонов', 'cases', 1),
('Зарядные устройства', 'Зарядные устройства и кабели', 'chargers', 2),
('Наушники', 'Наушники и гарнитуры', 'headphones', 3),
('Защитные стекла', 'Защитные стекла и пленки', 'screen-protectors', 4),
('Беспроводные зарядки', 'Беспроводные зарядные устройства', 'wireless-chargers', 5),
('Кабели и адаптеры', 'USB кабели и адаптеры', 'cables-adapters', 6),
('Портативные аккумуляторы', 'Power Bank и внешние батареи', 'power-banks', 7),
('Подставки и держатели', 'Подставки для телефонов', 'stands-holders', 8),
('Стилусы', 'Стилусы для сенсорных экранов', 'styluses', 9),
-- Компьютерные аксессуары
('Клавиатуры', 'Механические и мембранные клавиатуры', 'keyboards', 10),
('Мыши', 'Компьютерные мыши и трекболы', 'mice', 11),
('Коврики для мыши', 'Игровые и офисные коврики', 'mouse-pads', 12),
('Веб-камеры', 'Веб-камеры для стриминга и конференций', 'webcams', 13),
('Микрофоны', 'Микрофоны для стриминга и записи', 'microphones', 14),
-- Фото и видео
('Штативы', 'Штативы для камер и телефонов', 'tripods', 15),
('Освещение', 'Световое оборудование для фото и видео', 'lighting', 16),
('Карты памяти', 'SD карты и флеш-накопители', 'memory-cards', 17),
-- Автомобильные аксессуары
('Автомобильные держатели', 'Держатели для телефонов в авто', 'car-holders', 18),
('Автомобильные зарядки', 'Зарядные устройства для авто', 'car-chargers', 19),
-- Игровые аксессуары
('Игровые контроллеры', 'Геймпады и контроллеры', 'game-controllers', 20),
('Игровые кресла', 'Кресла для геймеров', 'gaming-chairs', 21)
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
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Триггер для addresses удален - таблица больше не существует
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_cart_items_updated_at BEFORE UPDATE ON cart_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_coupons_updated_at BEFORE UPDATE ON coupons FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Триггер для shipping_methods удален - таблица больше не существует
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- ТЕСТОВЫЕ ДАННЫЕ
-- =============================================

-- Создание тестового пользователя
INSERT INTO users (email, password, first_name, last_name, phone, gender, is_active, is_admin) VALUES 
('admin@shop.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Админ', 'Админов', '+7 (999) 123-45-67', 'male', true, true),
('user@shop.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Иван', 'Петров', '+7 (999) 765-43-21', 'male', true, false)
ON CONFLICT (email) DO NOTHING;

-- Создание тестовых товаров с разными брендами
INSERT INTO products (name, slug, description, short_description, price, compare_price, sku, stock, is_active, is_featured, is_new, weight, dimensions, brand, model, color, material, category_id, tags, meta_title, meta_description) VALUES 
-- Чехлы для iPhone
('Чехол Apple Silicone Case для iPhone 15 Pro', 'apple-silicone-case-iphone-15-pro', 'Официальный силиконовый чехол Apple с мягкой внутренней поверхностью и внешней поверхностью из силикона.', 'Официальный чехол Apple для iPhone 15 Pro', 4990.00, 5990.00, 'APPLE-CASE-IP15P-BLUE', 25, true, true, true, 0.05, '15.5x7.8x1.2 cm', 'Apple', 'Silicone Case', 'Синий', 'Силикон', (SELECT id FROM categories WHERE slug = 'cases'), ARRAY['iPhone 15 Pro', 'официальный', 'силикон'], 'Чехол Apple для iPhone 15 Pro - Синий', 'Официальный силиконовый чехол Apple для iPhone 15 Pro с защитой от падений и поддержкой MagSafe'),

('Чехол Spigen Ultra Hybrid для iPhone 15', 'spigen-ultra-hybrid-iphone-15', 'Прозрачный чехол Spigen с защитой от падений и поддержкой беспроводной зарядки.', 'Прозрачный чехол Spigen для iPhone 15', 1990.00, 2490.00, 'SPIGEN-UH-IP15-CLEAR', 50, true, false, false, 0.08, '15.0x7.6x1.0 cm', 'Spigen', 'Ultra Hybrid', 'Прозрачный', 'TPU + Поликарбонат', (SELECT id FROM categories WHERE slug = 'cases'), ARRAY['iPhone 15', 'прозрачный', 'защита'], 'Чехол Spigen для iPhone 15 - Прозрачный', 'Прозрачный чехол Spigen Ultra Hybrid для iPhone 15 с защитой от падений'),

-- Зарядные устройства
('Зарядное устройство Apple 20W USB-C', 'apple-20w-usb-c-charger', 'Официальное зарядное устройство Apple мощностью 20W с разъемом USB-C.', 'Зарядка Apple 20W USB-C', 2990.00, 3990.00, 'APPLE-20W-USB-C', 30, true, true, false, 0.12, '5.6x4.2x2.7 cm', 'Apple', '20W USB-C Power Adapter', 'Белый', 'Пластик', (SELECT id FROM categories WHERE slug = 'chargers'), ARRAY['iPhone', 'iPad', 'быстрая зарядка'], 'Зарядка Apple 20W USB-C', 'Официальное зарядное устройство Apple 20W с разъемом USB-C для быстрой зарядки'),

('Беспроводная зарядка Anker PowerWave 7.5W', 'anker-powervave-7-5w-wireless-charger', 'Беспроводная зарядная станция Anker с поддержкой быстрой зарядки до 7.5W.', 'Беспроводная зарядка Anker 7.5W', 2490.00, 2990.00, 'ANKER-PW-7.5W', 20, true, false, true, 0.15, '10.0x10.0x1.5 cm', 'Anker', 'PowerWave 7.5W', 'Черный', 'Пластик + Силикон', (SELECT id FROM categories WHERE slug = 'wireless-chargers'), ARRAY['беспроводная зарядка', 'Qi', 'быстрая зарядка'], 'Беспроводная зарядка Anker 7.5W', 'Беспроводная зарядная станция Anker PowerWave 7.5W с поддержкой Qi стандарта'),

-- Наушники
('Наушники Apple AirPods Pro 2', 'apple-airpods-pro-2', 'Беспроводные наушники Apple AirPods Pro 2-го поколения с активным шумоподавлением.', 'AirPods Pro 2 с шумоподавлением', 19990.00, 22990.00, 'APPLE-AIRPODS-PRO-2', 15, true, true, true, 0.06, '6.0x4.5x2.1 cm', 'Apple', 'AirPods Pro 2', 'Белый', 'Пластик', (SELECT id FROM categories WHERE slug = 'headphones'), ARRAY['беспроводные', 'шумоподавление', 'пространственное аудио'], 'Apple AirPods Pro 2', 'Беспроводные наушники Apple AirPods Pro 2 с активным шумоподавлением и пространственным аудио'),

('Наушники JBL Tune 760NC', 'jbl-tune-760nc', 'Беспроводные наушники JBL с активным шумоподавлением и 50-часовым временем работы.', 'JBL Tune 760NC с шумоподавлением', 8990.00, 10990.00, 'JBL-TUNE-760NC', 25, true, false, false, 0.25, '20.0x18.0x8.0 cm', 'JBL', 'Tune 760NC', 'Черный', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'headphones'), ARRAY['беспроводные', 'шумоподавление', 'долгая работа'], 'JBL Tune 760NC', 'Беспроводные наушники JBL Tune 760NC с активным шумоподавлением и 50-часовым временем работы'),

-- Защитные стекла
('Защитное стекло Belkin ScreenForce для iPhone 15', 'belkin-screenforce-iphone-15', 'Защитное стекло Belkin с технологией ScreenForce для iPhone 15.', 'Защитное стекло Belkin для iPhone 15', 1990.00, 2490.00, 'BELKIN-SF-IP15', 40, true, false, true, 0.02, '15.0x7.6x0.3 cm', 'Belkin', 'ScreenForce', 'Прозрачный', 'Закаленное стекло', (SELECT id FROM categories WHERE slug = 'screen-protectors'), ARRAY['iPhone 15', 'защитное стекло', '9H'], 'Защитное стекло Belkin для iPhone 15', 'Защитное стекло Belkin ScreenForce для iPhone 15 с твердостью 9H'),

-- Кабели
('USB-C кабель UGREEN 1м', 'ugreen-usb-c-cable-1m', 'USB-C кабель UGREEN длиной 1 метр с поддержкой быстрой зарядки и передачи данных.', 'USB-C кабель UGREEN 1м', 590.00, 790.00, 'UGREEN-USB-C-1M', 100, true, false, false, 0.05, '100.0x0.5x0.3 cm', 'UGREEN', 'USB-C Cable', 'Черный', 'Медь + Пластик', (SELECT id FROM categories WHERE slug = 'cables-adapters'), ARRAY['USB-C', 'быстрая зарядка', 'передача данных'], 'USB-C кабель UGREEN 1м', 'USB-C кабель UGREEN длиной 1 метр с поддержкой быстрой зарядки'),

-- Power Bank
('Портативный аккумулятор Anker PowerCore 10000', 'anker-powercore-10000', 'Портативный аккумулятор Anker PowerCore емкостью 10000 мАч с быстрой зарядкой.', 'PowerBank Anker 10000 мАч', 2990.00, 3490.00, 'ANKER-PC-10000', 35, true, true, false, 0.2, '9.0x6.0x2.2 cm', 'Anker', 'PowerCore 10000', 'Черный', 'Пластик', (SELECT id FROM categories WHERE slug = 'power-banks'), ARRAY['10000 мАч', 'быстрая зарядка', 'компактный'], 'PowerBank Anker 10000 мАч', 'Портативный аккумулятор Anker PowerCore емкостью 10000 мАч с быстрой зарядкой'),

-- Подставки
('Подставка для телефона Baseus Metal Stand', 'baseus-metal-stand-phone', 'Металлическая подставка Baseus для телефона с регулируемым углом наклона.', 'Металлическая подставка Baseus', 1290.00, 1590.00, 'BASEUS-MS-PHONE', 60, true, false, false, 0.15, '15.0x8.0x2.0 cm', 'Baseus', 'Metal Stand', 'Серебристый', 'Алюминий', (SELECT id FROM categories WHERE slug = 'stands-holders'), ARRAY['подставка', 'металлическая', 'регулируемая'], 'Подставка Baseus для телефона', 'Металлическая подставка Baseus для телефона с регулируемым углом наклона')
ON CONFLICT (sku) DO NOTHING;

-- Создание тестовых изображений для товаров
INSERT INTO images (product_id, cloudinary_public_id, url, alt, is_primary, sort_order) VALUES 
-- Изображения для чехла Apple
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P-BLUE'), 'apple-case-ip15p-blue-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-case-ip15p-blue-1.jpg', 'Чехол Apple для iPhone 15 Pro - Синий - Вид спереди', true, 1),
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P-BLUE'), 'apple-case-ip15p-blue-2', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-case-ip15p-blue-2.jpg', 'Чехол Apple для iPhone 15 Pro - Синий - Вид сбоку', false, 2),

-- Изображения для чехла Spigen
((SELECT id FROM products WHERE sku = 'SPIGEN-UH-IP15-CLEAR'), 'spigen-uh-ip15-clear-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/spigen-uh-ip15-clear-1.jpg', 'Чехол Spigen Ultra Hybrid для iPhone 15 - Прозрачный', true, 1),

-- Изображения для зарядки Apple
((SELECT id FROM products WHERE sku = 'APPLE-20W-USB-C'), 'apple-20w-usb-c-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-20w-usb-c-1.jpg', 'Зарядное устройство Apple 20W USB-C', true, 1),

-- Изображения для беспроводной зарядки Anker
((SELECT id FROM products WHERE sku = 'ANKER-PW-7.5W'), 'anker-pw-7.5w-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/anker-pw-7.5w-1.jpg', 'Беспроводная зарядка Anker PowerWave 7.5W', true, 1),

-- Изображения для AirPods Pro
((SELECT id FROM products WHERE sku = 'APPLE-AIRPODS-PRO-2'), 'apple-airpods-pro-2-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-airpods-pro-2-1.jpg', 'Apple AirPods Pro 2 - Белый', true, 1),
((SELECT id FROM products WHERE sku = 'APPLE-AIRPODS-PRO-2'), 'apple-airpods-pro-2-2', 'https://res.cloudinary.com/your-cloud/image/upload/v1/apple-airpods-pro-2-2.jpg', 'Apple AirPods Pro 2 - В кейсе', false, 2),

-- Изображения для JBL наушников
((SELECT id FROM products WHERE sku = 'JBL-TUNE-760NC'), 'jbl-tune-760nc-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/jbl-tune-760nc-1.jpg', 'JBL Tune 760NC - Черный', true, 1),

-- Изображения для защитного стекла Belkin
((SELECT id FROM products WHERE sku = 'BELKIN-SF-IP15'), 'belkin-sf-ip15-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/belkin-sf-ip15-1.jpg', 'Защитное стекло Belkin для iPhone 15', true, 1),

-- Изображения для USB-C кабеля
((SELECT id FROM products WHERE sku = 'UGREEN-USB-C-1M'), 'ugreen-usb-c-1m-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/ugreen-usb-c-1m-1.jpg', 'USB-C кабель UGREEN 1м', true, 1),

-- Изображения для PowerBank
((SELECT id FROM products WHERE sku = 'ANKER-PC-10000'), 'anker-pc-10000-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/anker-pc-10000-1.jpg', 'PowerBank Anker 10000 мАч', true, 1),

-- Изображения для подставки Baseus
((SELECT id FROM products WHERE sku = 'BASEUS-MS-PHONE'), 'baseus-ms-phone-1', 'https://res.cloudinary.com/your-cloud/image/upload/v1/baseus-ms-phone-1.jpg', 'Подставка Baseus для телефона', true, 1)
ON CONFLICT DO NOTHING;