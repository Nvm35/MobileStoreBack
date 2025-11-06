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
    role VARCHAR(20) DEFAULT 'customer' CHECK (role IN ('admin', 'manager', 'customer')),
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
    is_active BOOLEAN DEFAULT true,
    feature BOOLEAN DEFAULT false, -- флаг особенного товара для витрины
    brand VARCHAR(255) NOT NULL, -- бренд как строка
    model VARCHAR(255),
    material VARCHAR(255),
    category_id UUID NOT NULL REFERENCES categories(id),
    tags TEXT[], -- массив тегов
    video_url TEXT, -- ссылка на видео товара
    view_count INTEGER DEFAULT 0, -- количество просмотров
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Создание таблицы складов/филиалов
CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL, -- название склада/филиала
    slug VARCHAR(255) NOT NULL UNIQUE, -- slug для URL
    address TEXT NOT NULL, -- адрес склада
    city VARCHAR(255) NOT NULL, -- город
    phone VARCHAR(20), -- телефон склада
    email VARCHAR(255), -- email склада
    is_active BOOLEAN DEFAULT true,
    is_main BOOLEAN DEFAULT false, -- главный склад
    manager_name VARCHAR(255), -- имя менеджера склада
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Создание таблицы вариантов товаров (зависит от products)
CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(255) NOT NULL UNIQUE, -- уникальный SKU для варианта
    name VARCHAR(255) NOT NULL, -- название варианта (например, "Красный, L")
    color VARCHAR(100), -- цвет варианта
    size VARCHAR(50), -- размер варианта
    price DECIMAL(10,2) NOT NULL, -- цена варианта (может отличаться от базовой)
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6. Создание таблицы остатков товаров по складам
CREATE TABLE IF NOT EXISTS warehouse_stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    product_variant_id UUID NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    stock INTEGER NOT NULL DEFAULT 0, -- остаток товара на складе
    reserved_stock INTEGER NOT NULL DEFAULT 0, -- зарезервированный товар
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(warehouse_id, product_variant_id) -- уникальная комбинация склада и варианта товара
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
-- Мы не храним товары без логина, поэтому все записи должны иметь user_id
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(255), -- опциональное поле для объединения корзины после логина
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price DECIMAL(10,2) NOT NULL, -- цена на момент добавления
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP, -- срок действия корзины
    CONSTRAINT cart_items_user_product_unique UNIQUE(user_id, product_id) -- один товар на пользователя
);

-- 7. Создание таблицы избранного (зависит от users, products)
CREATE TABLE IF NOT EXISTS wishlist_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id) -- один товар на пользователя
);

-- 8. Создание таблицы заказов (зависит от users, warehouses)
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    warehouse_id UUID REFERENCES warehouses(id), -- склад, с которого выполняется заказ
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
CREATE INDEX IF NOT EXISTS idx_products_feature ON products(feature);
CREATE INDEX IF NOT EXISTS idx_products_feature_active ON products(feature, is_active);
CREATE INDEX IF NOT EXISTS idx_warehouses_city ON warehouses(city);
CREATE INDEX IF NOT EXISTS idx_warehouses_active ON warehouses(is_active);
CREATE INDEX IF NOT EXISTS idx_warehouses_main ON warehouses(is_main);
CREATE INDEX IF NOT EXISTS idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX IF NOT EXISTS idx_product_variants_sku ON product_variants(sku);
CREATE INDEX IF NOT EXISTS idx_product_variants_active ON product_variants(is_active);
CREATE INDEX IF NOT EXISTS idx_warehouse_stocks_warehouse_id ON warehouse_stocks(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_stocks_variant_id ON warehouse_stocks(product_variant_id);
CREATE INDEX IF NOT EXISTS idx_warehouse_stocks_warehouse_variant ON warehouse_stocks(warehouse_id, product_variant_id);
CREATE INDEX IF NOT EXISTS idx_images_product_id ON images(product_id);
CREATE INDEX IF NOT EXISTS idx_images_primary ON images(is_primary);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_warehouse_id ON orders(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);
CREATE INDEX IF NOT EXISTS idx_order_items_variant_id ON order_items(product_variant_id);


-- Индексы для корзины
CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
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
CREATE TRIGGER update_warehouses_updated_at BEFORE UPDATE ON warehouses FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_product_variants_updated_at BEFORE UPDATE ON product_variants FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_warehouse_stocks_updated_at BEFORE UPDATE ON warehouse_stocks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_cart_items_updated_at BEFORE UPDATE ON cart_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- ТЕСТОВЫЕ ДАННЫЕ
-- =============================================

-- Создание тестового пользователя
INSERT INTO users (email, password, first_name, last_name, phone, is_active, role) VALUES 
('admin@shop.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Админ', 'Админов', '+7 (999) 123-45-67', true, 'admin'),
('manager@shop.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Менеджер', 'Менеджеров', '+7 (999) 111-22-33', true, 'manager'),
('user@shop.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Иван', 'Петров', '+7 (999) 765-43-21', true, 'customer')
ON CONFLICT (email) DO NOTHING;

-- Создание тестовых складов/филиалов
INSERT INTO warehouses (name, slug, address, city, phone, email, is_active, is_main, manager_name) VALUES 
('Главный склад', 'main-warehouse', 'ул. Промышленная, 15', 'Москва', '+7 (495) 123-45-67', 'main@shop.com', true, true, 'Иванов Иван Иванович'),
('Филиал "Центр"', 'center-branch', 'ул. Тверская, 25', 'Москва', '+7 (495) 234-56-78', 'center@shop.com', true, false, 'Петров Петр Петрович'),
('Филиал "Север"', 'north-branch', 'пр. Мира, 100', 'Москва', '+7 (495) 345-67-89', 'north@shop.com', true, false, 'Сидоров Сидор Сидорович'),
('Филиал "Юг"', 'south-branch', 'ул. Ленина, 50', 'Санкт-Петербург', '+7 (812) 456-78-90', 'south@shop.com', true, false, 'Козлов Козел Козлович')
ON CONFLICT DO NOTHING;

-- Создание тестовых товаров с разными брендами
INSERT INTO products (name, slug, description, base_price, sku, is_active, feature, brand, model, material, category_id, tags, video_url) VALUES 
-- Чехлы для iPhone
('Чехол Apple Silicone Case для iPhone 15 Pro', 'apple-silicone-case-iphone-15-pro', 'Официальный силиконовый чехол Apple с мягкой внутренней поверхностью и внешней поверхностью из силикона.', 4990.00, 'APPLE-CASE-IP15P', true, false, 'Apple', 'Silicone Case', 'Силикон', (SELECT id FROM categories WHERE slug = 'cases'), ARRAY['iPhone 15 Pro', 'официальный', 'силикон'], 'https://www.youtube.com/watch?v=example1'),

('Чехол Spigen Ultra Hybrid для iPhone 15', 'spigen-ultra-hybrid-iphone-15', 'Прозрачный чехол Spigen с защитой от падений и поддержкой беспроводной зарядки.', 1990.00, 'SPIGEN-UH-IP15', true, false, 'Spigen', 'Ultra Hybrid', 'TPU + Поликарбонат', (SELECT id FROM categories WHERE slug = 'cases'), ARRAY['iPhone 15', 'прозрачный', 'защита']),

-- Зарядные устройства
('Зарядное устройство Apple 20W USB-C', 'apple-20w-usb-c-charger', 'Официальное зарядное устройство Apple мощностью 20W с разъемом USB-C.', 2990.00, 'APPLE-20W-USB-C', true, 'Apple', '20W USB-C Power Adapter', 'Пластик', (SELECT id FROM categories WHERE slug = 'chargers'), ARRAY['iPhone', 'iPad', 'быстрая зарядка']),

('Беспроводная зарядка Anker PowerWave 7.5W', 'anker-powervave-7-5w-wireless-charger', 'Беспроводная зарядная станция Anker с поддержкой быстрой зарядки до 7.5W.', 2490.00, 'ANKER-PW-7.5W', true, 'Anker', 'PowerWave 7.5W', 'Пластик + Силикон', (SELECT id FROM categories WHERE slug = 'wireless-chargers'), ARRAY['беспроводная зарядка', 'Qi', 'быстрая зарядка']),

-- Наушники
('Наушники Apple AirPods Pro 2', 'apple-airpods-pro-2', 'Беспроводные наушники Apple AirPods Pro 2-го поколения с активным шумоподавлением.', 19990.00, 'APPLE-AIRPODS-PRO-2', true, 'Apple', 'AirPods Pro 2', 'Пластик', (SELECT id FROM categories WHERE slug = 'headphones'), ARRAY['беспроводные', 'шумоподавление', 'пространственное аудио'], 'https://www.youtube.com/watch?v=example2'),

('Наушники JBL Tune 760NC', 'jbl-tune-760nc', 'Беспроводные наушники JBL с активным шумоподавлением и 50-часовым временем работы.', 8990.00, 'JBL-TUNE-760NC', true, 'JBL', 'Tune 760NC', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'headphones'), ARRAY['беспроводные', 'шумоподавление', 'долгая работа']),

-- Защитные стекла
('Защитное стекло Belkin ScreenForce для iPhone 15', 'belkin-screenforce-iphone-15', 'Защитное стекло Belkin с технологией ScreenForce для iPhone 15.', 1990.00, 'BELKIN-SF-IP15', true, 'Belkin', 'ScreenForce', 'Закаленное стекло', (SELECT id FROM categories WHERE slug = 'screen-protectors'), ARRAY['iPhone 15', 'защитное стекло', '9H']),

-- Кабели
('USB-C кабель UGREEN 1м', 'ugreen-usb-c-cable-1m', 'USB-C кабель UGREEN длиной 1 метр с поддержкой быстрой зарядки и передачи данных.', 590.00, 'UGREEN-USB-C-1M', true, 'UGREEN', 'USB-C Cable', 'Медь + Пластик', (SELECT id FROM categories WHERE slug = 'cables-adapters'), ARRAY['USB-C', 'быстрая зарядка', 'передача данных']),

-- Power Bank
('Портативный аккумулятор Anker PowerCore 10000', 'anker-powercore-10000', 'Портативный аккумулятор Anker PowerCore емкостью 10000 мАч с быстрой зарядкой.', 2990.00, 'ANKER-PC-10000', true, 'Anker', 'PowerCore 10000', 'Пластик', (SELECT id FROM categories WHERE slug = 'power-banks'), ARRAY['10000 мАч', 'быстрая зарядка', 'компактный']),

-- Подставки
('Подставка для телефона Baseus Metal Stand', 'baseus-metal-stand-phone', 'Металлическая подставка Baseus для телефона с регулируемым углом наклона.', 1290.00, 'BASEUS-MS-PHONE', true, 'Baseus', 'Metal Stand', 'Алюминий', (SELECT id FROM categories WHERE slug = 'stands-holders'), ARRAY['подставка', 'металлическая', 'регулируемая']),

-- Дополнительные товары для разных категорий
-- Чехлы
('Чехол OtterBox Defender для iPhone 15', 'otterbox-defender-iphone-15', 'Защитный чехол OtterBox Defender с тройной защитой для iPhone 15.', 3990.00, 'OTTERBOX-DEF-IP15', true, 'OtterBox', 'Defender', 'Пластик + Силикон', (SELECT id FROM categories WHERE slug = 'cases'), ARRAY['iPhone 15', 'защитный', 'тройная защита']),
('Чехол Casetify для iPhone 15 Pro Max', 'casetify-iphone-15-pro-max', 'Стильный чехол Casetify с возможностью кастомизации для iPhone 15 Pro Max.', 2990.00, 'CASETIFY-IP15PM', true, 'Casetify', 'Custom Case', 'Пластик', (SELECT id FROM categories WHERE slug = 'cases'), ARRAY['iPhone 15 Pro Max', 'кастомизация', 'стильный']),

-- Зарядные устройства
('Зарядное устройство Samsung 25W', 'samsung-25w-charger', 'Быстрое зарядное устройство Samsung мощностью 25W с разъемом USB-C.', 2490.00, 'SAMSUNG-25W', true, 'Samsung', '25W Fast Charger', 'Пластик', (SELECT id FROM categories WHERE slug = 'chargers'), ARRAY['Samsung', 'быстрая зарядка', '25W']),
('Зарядное устройство Xiaomi 67W', 'xiaomi-67w-charger', 'Мощное зарядное устройство Xiaomi 67W для быстрой зарядки смартфонов.', 1990.00, 'XIAOMI-67W', true, 'Xiaomi', '67W Charger', 'Пластик', (SELECT id FROM categories WHERE slug = 'chargers'), ARRAY['Xiaomi', '67W', 'быстрая зарядка']),

-- Наушники
('Наушники Sony WH-1000XM5', 'sony-wh-1000xm5', 'Премиальные беспроводные наушники Sony с лучшим в мире шумоподавлением.', 24990.00, 'SONY-WH-1000XM5', true, 'Sony', 'WH-1000XM5', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'headphones'), ARRAY['Sony', 'шумоподавление', 'премиум']),
('Наушники Sennheiser HD 660S', 'sennheiser-hd-660s', 'Студийные наушники Sennheiser HD 660S для профессионального прослушивания.', 18990.00, 'SENNHEISER-HD-660S', true, 'Sennheiser', 'HD 660S', 'Металл + Кожа', (SELECT id FROM categories WHERE slug = 'headphones'), ARRAY['Sennheiser', 'студийные', 'профессиональные']),

-- Защитные стекла
('Защитное стекло ESR для iPhone 15 Pro', 'esr-screen-protector-iphone-15-pro', 'Защитное стекло ESR с технологией Easy Installation для iPhone 15 Pro.', 1290.00, 'ESR-SP-IP15P', true, 'ESR', 'Easy Installation', 'Закаленное стекло', (SELECT id FROM categories WHERE slug = 'screen-protectors'), ARRAY['iPhone 15 Pro', 'легкая установка', 'защитное стекло']),
('Защитное стекло ZAGG для Samsung Galaxy S24', 'zagg-screen-protector-galaxy-s24', 'Защитное стекло ZAGG InvisibleShield для Samsung Galaxy S24.', 1990.00, 'ZAGG-SP-S24', true, 'ZAGG', 'InvisibleShield', 'Закаленное стекло', (SELECT id FROM categories WHERE slug = 'screen-protectors'), ARRAY['Samsung Galaxy S24', 'InvisibleShield', 'защита']),

-- Беспроводные зарядки
('Беспроводная зарядка Samsung 15W', 'samsung-15w-wireless-charger', 'Беспроводная зарядка Samsung мощностью 15W с быстрой зарядкой.', 2990.00, 'SAMSUNG-15W-WC', true, 'Samsung', '15W Wireless Charger', 'Пластик + Силикон', (SELECT id FROM categories WHERE slug = 'wireless-chargers'), ARRAY['Samsung', '15W', 'быстрая зарядка']),
('Беспроводная зарядка Belkin 10W', 'belkin-10w-wireless-charger', 'Беспроводная зарядка Belkin с поддержкой Qi стандарта.', 2490.00, 'BELKIN-10W-WC', true, 'Belkin', '10W Wireless Charger', 'Пластик', (SELECT id FROM categories WHERE slug = 'wireless-chargers'), ARRAY['Belkin', '10W', 'Qi стандарт']),

-- Кабели
('Lightning кабель Apple 1м', 'apple-lightning-cable-1m', 'Официальный Lightning кабель Apple длиной 1 метр.', 1990.00, 'APPLE-LIGHTNING-1M', true, 'Apple', 'Lightning Cable', 'Медь + Пластик', (SELECT id FROM categories WHERE slug = 'cables-adapters'), ARRAY['Apple', 'Lightning', 'официальный']),
('USB-C кабель Anker PowerLine III', 'anker-powerline-iii-usb-c', 'Прочный USB-C кабель Anker PowerLine III с пожизненной гарантией.', 1290.00, 'ANKER-PL3-USB-C', true, 'Anker', 'PowerLine III', 'Нейлон + Металл', (SELECT id FROM categories WHERE slug = 'cables-adapters'), ARRAY['Anker', 'PowerLine III', 'прочный']),

-- Power Bank
('Портативный аккумулятор Xiaomi Power Bank 3', 'xiaomi-power-bank-3', 'Портативный аккумулятор Xiaomi Power Bank 3 емкостью 20000 мАч.', 3990.00, 'XIAOMI-PB-20000', true, 'Xiaomi', 'Power Bank 3', 'Пластик', (SELECT id FROM categories WHERE slug = 'power-banks'), ARRAY['Xiaomi', '20000 мАч', 'Power Bank 3']),
('Портативный аккумулятор RAVPower 20000', 'ravpower-20000-power-bank', 'Портативный аккумулятор RAVPower емкостью 20000 мАч с быстрой зарядкой.', 3490.00, 'RAVPOWER-20000', true, 'RAVPower', '20000 mAh', 'Пластик', (SELECT id FROM categories WHERE slug = 'power-banks'), ARRAY['RAVPower', '20000 мАч', 'быстрая зарядка']),

-- Подставки
('Подставка для телефона Lamicall Adjustable', 'lamicall-adjustable-stand', 'Регулируемая подставка Lamicall для телефона с углом наклона 0-180°.', 1590.00, 'LAMICALL-ADJUSTABLE', true, 'Lamicall', 'Adjustable Stand', 'Алюминий', (SELECT id FROM categories WHERE slug = 'stands-holders'), ARRAY['Lamicall', 'регулируемая', '0-180°']),
('Подставка для телефона Spigen ArcStation', 'spigen-arcstation-stand', 'Компактная подставка Spigen ArcStation для телефона.', 1290.00, 'SPIGEN-ARCSTATION', true, 'Spigen', 'ArcStation', 'Пластик', (SELECT id FROM categories WHERE slug = 'stands-holders'), ARRAY['Spigen', 'компактная', 'ArcStation']),

-- Стилусы
('Apple Pencil 2', 'apple-pencil-2', 'Apple Pencil 2-го поколения для iPad с магнитным креплением и беспроводной зарядкой.', 8990.00, 'APPLE-PENCIL-2', true, 'Apple', 'Pencil 2', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'styluses'), ARRAY['Apple', 'iPad', 'магнитное крепление']),
('Samsung S Pen для Galaxy Tab', 'samsung-s-pen-galaxy-tab', 'Samsung S Pen для Galaxy Tab с улучшенной чувствительностью.', 2990.00, 'SAMSUNG-S-PEN', true, 'Samsung', 'S Pen', 'Пластик', (SELECT id FROM categories WHERE slug = 'styluses'), ARRAY['Samsung', 'Galaxy Tab', 'чувствительность']),

-- Клавиатуры
('Клавиатура Logitech MX Keys', 'logitech-mx-keys', 'Беспроводная клавиатура Logitech MX Keys с подсветкой и мультиустройством.', 8990.00, 'LOGITECH-MX-KEYS', true, 'Logitech', 'MX Keys', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'keyboards'), ARRAY['Logitech', 'беспроводная', 'подсветка']),
('Механическая клавиатура Keychron K2', 'keychron-k2-mechanical', 'Компактная механическая клавиатура Keychron K2 с переключателями Gateron.', 5990.00, 'KEYCHRON-K2', true, 'Keychron', 'K2', 'Пластик + Алюминий', (SELECT id FROM categories WHERE slug = 'keyboards'), ARRAY['Keychron', 'механическая', 'Gateron']),

-- Мыши
('Мышь Logitech MX Master 3', 'logitech-mx-master-3', 'Беспроводная мышь Logitech MX Master 3 с эргономичным дизайном.', 6990.00, 'LOGITECH-MX-MASTER-3', true, 'Logitech', 'MX Master 3', 'Пластик + Резина', (SELECT id FROM categories WHERE slug = 'mice'), ARRAY['Logitech', 'беспроводная', 'эргономичная']),
('Игровая мышь Razer DeathAdder V3', 'razer-deathadder-v3', 'Игровая мышь Razer DeathAdder V3 с оптическим сенсором Focus Pro 30K.', 7990.00, 'RAZER-DEATHADDER-V3', true, 'Razer', 'DeathAdder V3', 'Пластик + Резина', (SELECT id FROM categories WHERE slug = 'mice'), ARRAY['Razer', 'игровая', 'Focus Pro 30K']),

-- Коврики для мыши
('Игровой коврик SteelSeries QcK', 'steelseries-qck-gaming-pad', 'Игровой коврик SteelSeries QcK с оптимальной поверхностью для игр.', 1290.00, 'STEELSERIES-QCK', true, 'SteelSeries', 'QcK', 'Ткань', (SELECT id FROM categories WHERE slug = 'mouse-pads'), ARRAY['SteelSeries', 'игровой', 'QcK']),
('Коврик для мыши Corsair MM300', 'corsair-mm300-mouse-pad', 'Большой коврик для мыши Corsair MM300 с прошитыми краями.', 1590.00, 'CORSAIR-MM300', true, 'Corsair', 'MM300', 'Ткань', (SELECT id FROM categories WHERE slug = 'mouse-pads'), ARRAY['Corsair', 'большой', 'прошитые края']),

-- Веб-камеры
('Веб-камера Logitech C920 HD Pro', 'logitech-c920-hd-pro', 'Веб-камера Logitech C920 HD Pro с разрешением 1080p и автофокусом.', 6990.00, 'LOGITECH-C920', true, 'Logitech', 'C920 HD Pro', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'webcams'), ARRAY['Logitech', '1080p', 'автофокус']),
('Веб-камера Razer Kiyo Pro', 'razer-kiyo-pro-webcam', 'Веб-камера Razer Kiyo Pro с HDR и встроенной подсветкой.', 9990.00, 'RAZER-KIYO-PRO', true, 'Razer', 'Kiyo Pro', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'webcams'), ARRAY['Razer', 'HDR', 'подсветка']),

-- Микрофоны
('Микрофон Blue Yeti USB', 'blue-yeti-usb-microphone', 'USB микрофон Blue Yeti с тремя капсулами и четырьмя режимами записи.', 8990.00, 'BLUE-YETI-USB', true, 'Blue', 'Yeti USB', 'Металл + Пластик', (SELECT id FROM categories WHERE slug = 'microphones'), ARRAY['Blue', 'USB', 'три капсулы']),
('Микрофон Shure SM7B', 'shure-sm7b-microphone', 'Динамический микрофон Shure SM7B для профессиональной записи.', 24990.00, 'SHURE-SM7B', true, 'Shure', 'SM7B', 'Металл', (SELECT id FROM categories WHERE slug = 'microphones'), ARRAY['Shure', 'динамический', 'профессиональный']),

-- Штативы
('Штатив Manfrotto Compact Action', 'manfrotto-compact-action-tripod', 'Компактный штатив Manfrotto Compact Action для камер и телефонов.', 2990.00, 'MANFROTTO-COMPACT', true, 'Manfrotto', 'Compact Action', 'Алюминий', (SELECT id FROM categories WHERE slug = 'tripods'), ARRAY['Manfrotto', 'компактный', 'алюминий']),
('Штатив Joby GripTight PRO', 'joby-griptight-pro-tripod', 'Гибкий штатив Joby GripTight PRO для телефонов и экшн-камер.', 1990.00, 'JOBY-GRIPTIGHT-PRO', true, 'Joby', 'GripTight PRO', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'tripods'), ARRAY['Joby', 'гибкий', 'GripTight']),

-- Освещение
('Кольцевая лампа Neewer 18"', 'neewer-18-inch-ring-light', 'Кольцевая лампа Neewer 18 дюймов с регулируемой яркостью и цветовой температурой.', 3990.00, 'NEEWER-18-RING', true, 'Neewer', '18" Ring Light', 'Металл + Пластик', (SELECT id FROM categories WHERE slug = 'lighting'), ARRAY['Neewer', 'кольцевая лампа', '18 дюймов']),
('LED панель Godox SL-60W', 'godox-sl-60w-led-panel', 'LED панель Godox SL-60W с высокой цветопередачей для фото и видео.', 5990.00, 'GODOX-SL-60W', true, 'Godox', 'SL-60W', 'Металл + Пластик', (SELECT id FROM categories WHERE slug = 'lighting'), ARRAY['Godox', 'LED панель', '60W']),

-- Карты памяти
('SD карта SanDisk Extreme Pro 128GB', 'sandisk-extreme-pro-128gb', 'SD карта SanDisk Extreme Pro 128GB с высокой скоростью записи.', 2990.00, 'SANDISK-EXTREME-128GB', true, 'SanDisk', 'Extreme Pro', 'Пластик', (SELECT id FROM categories WHERE slug = 'memory-cards'), ARRAY['SanDisk', '128GB', 'высокая скорость']),
('MicroSD карта Samsung EVO Plus 256GB', 'samsung-evo-plus-256gb', 'MicroSD карта Samsung EVO Plus 256GB с адаптером.', 3990.00, 'SAMSUNG-EVO-256GB', true, 'Samsung', 'EVO Plus', 'Пластик', (SELECT id FROM categories WHERE slug = 'memory-cards'), ARRAY['Samsung', '256GB', 'MicroSD']),

-- Автомобильные держатели
('Автомобильный держатель iOttie Easy One Touch', 'iottie-easy-one-touch-holder', 'Автомобильный держатель iOttie Easy One Touch с одной рукой.', 2490.00, 'IOTTIE-EASY-TOUCH', true, 'iOttie', 'Easy One Touch', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'car-holders'), ARRAY['iOttie', 'автомобильный', 'одна рука']),
('Магнитный держатель Scosche MagicMount', 'scosche-magicmount-holder', 'Магнитный автомобильный держатель Scosche MagicMount.', 1990.00, 'SCOSCHE-MAGICMOUNT', true, 'Scosche', 'MagicMount', 'Металл + Магнит', (SELECT id FROM categories WHERE slug = 'car-holders'), ARRAY['Scosche', 'магнитный', 'MagicMount']),

-- Автомобильные зарядки
('Автомобильная зарядка Anker PowerDrive 2', 'anker-powerdrive-2-car-charger', 'Автомобильная зарядка Anker PowerDrive 2 с двумя USB портами.', 1290.00, 'ANKER-POWERDRIVE-2', true, 'Anker', 'PowerDrive 2', 'Пластик + Металл', (SELECT id FROM categories WHERE slug = 'car-chargers'), ARRAY['Anker', 'автомобильная', 'два USB']),
('Беспроводная автомобильная зарядка Belkin', 'belkin-wireless-car-charger', 'Беспроводная автомобильная зарядка Belkin с креплением на вентиляцию.', 2990.00, 'BELKIN-WIRELESS-CAR', true, 'Belkin', 'Wireless Car Charger', 'Пластик + Силикон', (SELECT id FROM categories WHERE slug = 'car-chargers'), ARRAY['Belkin', 'беспроводная', 'на вентиляцию']),

-- Игровые контроллеры
('Геймпад Xbox Wireless Controller', 'xbox-wireless-controller', 'Беспроводной геймпад Xbox с улучшенной эргономикой и тактильной обратной связью.', 4990.00, 'XBOX-WIRELESS-CTRL', true, 'Microsoft', 'Xbox Wireless Controller', 'Пластик + Резина', (SELECT id FROM categories WHERE slug = 'game-controllers'), ARRAY['Microsoft', 'Xbox', 'беспроводной']),
('Геймпад Sony DualSense', 'sony-dualsense-controller', 'Геймпад Sony DualSense для PlayStation 5 с адаптивными триггерами.', 5990.00, 'SONY-DUALSENSE', true, 'Sony', 'DualSense', 'Пластик + Резина', (SELECT id FROM categories WHERE slug = 'game-controllers'), ARRAY['Sony', 'PlayStation 5', 'адаптивные триггеры']),

-- Игровые кресла
('Игровое кресло DXRacer Formula Series', 'dxracer-formula-gaming-chair', 'Игровое кресло DXRacer Formula Series с эргономичной поддержкой спины.', 19990.00, 'DXRACER-FORMULA', true, 'DXRacer', 'Formula Series', 'Кожа + Пенополиуретан', (SELECT id FROM categories WHERE slug = 'gaming-chairs'), ARRAY['DXRacer', 'игровое кресло', 'эргономичное']),
('Игровое кресло Secretlab Titan', 'secretlab-titan-gaming-chair', 'Премиальное игровое кресло Secretlab Titan с регулируемой поясничной поддержкой.', 29990.00, 'SECRETLAB-TITAN', true, 'Secretlab', 'Titan', 'Кожа + Пенополиуретан', (SELECT id FROM categories WHERE slug = 'gaming-chairs'), ARRAY['Secretlab', 'премиум', 'поясничная поддержка'])
ON CONFLICT (sku) DO NOTHING;

-- Создание тестовых вариантов товаров
INSERT INTO product_variants (product_id, sku, name, color, size, price, is_active) VALUES 
-- Варианты чехла Apple
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P'), 'APPLE-CASE-IP15P-BLUE', 'Apple Silicone Case - Синий', 'Синий', NULL, 4990.00, true),
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P'), 'APPLE-CASE-IP15P-BLACK', 'Apple Silicone Case - Черный', 'Черный', NULL, 4990.00, true),
((SELECT id FROM products WHERE sku = 'APPLE-CASE-IP15P'), 'APPLE-CASE-IP15P-WHITE', 'Apple Silicone Case - Белый', 'Белый', NULL, 4990.00, true),

-- Варианты чехла Spigen
((SELECT id FROM products WHERE sku = 'SPIGEN-UH-IP15'), 'SPIGEN-UH-IP15-CLEAR', 'Spigen Ultra Hybrid - Прозрачный', 'Прозрачный', NULL, 1990.00, true),
((SELECT id FROM products WHERE sku = 'SPIGEN-UH-IP15'), 'SPIGEN-UH-IP15-BLACK', 'Spigen Ultra Hybrid - Черный', 'Черный', NULL, 1990.00, true),

-- Варианты наушников JBL (размеры)
((SELECT id FROM products WHERE sku = 'JBL-TUNE-760NC'), 'JBL-TUNE-760NC-BLACK', 'JBL Tune 760NC - Черный', 'Черный', 'One Size', 8990.00, true),
((SELECT id FROM products WHERE sku = 'JBL-TUNE-760NC'), 'JBL-TUNE-760NC-BLUE', 'JBL Tune 760NC - Синий', 'Синий', 'One Size', 8990.00, true),

-- Варианты кабеля (длины)
((SELECT id FROM products WHERE sku = 'UGREEN-USB-C-1M'), 'UGREEN-USB-C-1M', 'UGREEN USB-C кабель 1м', 'Черный', '1м', 590.00, true),
((SELECT id FROM products WHERE sku = 'UGREEN-USB-C-1M'), 'UGREEN-USB-C-2M', 'UGREEN USB-C кабель 2м', 'Черный', '2м', 790.00, true),
((SELECT id FROM products WHERE sku = 'UGREEN-USB-C-1M'), 'UGREEN-USB-C-3M', 'UGREEN USB-C кабель 3м', 'Черный', '3м', 990.00, true),

-- Варианты для новых товаров
-- Варианты чехла OtterBox
((SELECT id FROM products WHERE sku = 'OTTERBOX-DEF-IP15'), 'OTTERBOX-DEF-IP15-BLACK', 'OtterBox Defender - Черный', 'Черный', NULL, 3990.00, true),
((SELECT id FROM products WHERE sku = 'OTTERBOX-DEF-IP15'), 'OTTERBOX-DEF-IP15-BLUE', 'OtterBox Defender - Синий', 'Синий', NULL, 3990.00, true),

-- Варианты чехла Casetify
((SELECT id FROM products WHERE sku = 'CASETIFY-IP15PM'), 'CASETIFY-IP15PM-CLEAR', 'Casetify Custom Case - Прозрачный', 'Прозрачный', NULL, 2990.00, true),
((SELECT id FROM products WHERE sku = 'CASETIFY-IP15PM'), 'CASETIFY-IP15PM-BLACK', 'Casetify Custom Case - Черный', 'Черный', NULL, 2990.00, true),

-- Варианты наушников Sony (цвета)
((SELECT id FROM products WHERE sku = 'SONY-WH-1000XM5'), 'SONY-WH-1000XM5-BLACK', 'Sony WH-1000XM5 - Черный', 'Черный', 'One Size', 24990.00, true),
((SELECT id FROM products WHERE sku = 'SONY-WH-1000XM5'), 'SONY-WH-1000XM5-SILVER', 'Sony WH-1000XM5 - Серебристый', 'Серебристый', 'One Size', 24990.00, true),

-- Варианты наушников Sennheiser
((SELECT id FROM products WHERE sku = 'SENNHEISER-HD-660S'), 'SENNHEISER-HD-660S-BLACK', 'Sennheiser HD 660S - Черный', 'Черный', 'One Size', 18990.00, true),

-- Варианты защитных стекол
((SELECT id FROM products WHERE sku = 'ESR-SP-IP15P'), 'ESR-SP-IP15P-CLEAR', 'ESR Screen Protector - Прозрачный', 'Прозрачный', 'iPhone 15 Pro', 1290.00, true),
((SELECT id FROM products WHERE sku = 'ZAGG-SP-S24'), 'ZAGG-SP-S24-CLEAR', 'ZAGG InvisibleShield - Прозрачный', 'Прозрачный', 'Galaxy S24', 1990.00, true),

-- Варианты беспроводных зарядок
((SELECT id FROM products WHERE sku = 'SAMSUNG-15W-WC'), 'SAMSUNG-15W-WC-BLACK', 'Samsung 15W Wireless Charger - Черный', 'Черный', 'One Size', 2990.00, true),
((SELECT id FROM products WHERE sku = 'SAMSUNG-15W-WC'), 'SAMSUNG-15W-WC-WHITE', 'Samsung 15W Wireless Charger - Белый', 'Белый', 'One Size', 2990.00, true),
((SELECT id FROM products WHERE sku = 'BELKIN-10W-WC'), 'BELKIN-10W-WC-BLACK', 'Belkin 10W Wireless Charger - Черный', 'Черный', 'One Size', 2490.00, true),

-- Варианты кабелей (длины)
((SELECT id FROM products WHERE sku = 'APPLE-LIGHTNING-1M'), 'APPLE-LIGHTNING-1M', 'Apple Lightning кабель 1м', 'Белый', '1м', 1990.00, true),
((SELECT id FROM products WHERE sku = 'APPLE-LIGHTNING-1M'), 'APPLE-LIGHTNING-2M', 'Apple Lightning кабель 2м', 'Белый', '2м', 2490.00, true),
((SELECT id FROM products WHERE sku = 'ANKER-PL3-USB-C'), 'ANKER-PL3-USB-C-1M', 'Anker PowerLine III USB-C 1м', 'Черный', '1м', 1290.00, true),
((SELECT id FROM products WHERE sku = 'ANKER-PL3-USB-C'), 'ANKER-PL3-USB-C-2M', 'Anker PowerLine III USB-C 2м', 'Черный', '2м', 1590.00, true),

-- Варианты Power Bank (цвета)
((SELECT id FROM products WHERE sku = 'XIAOMI-PB-20000'), 'XIAOMI-PB-20000-BLACK', 'Xiaomi Power Bank 3 - Черный', 'Черный', '20000 мАч', 3990.00, true),
((SELECT id FROM products WHERE sku = 'XIAOMI-PB-20000'), 'XIAOMI-PB-20000-WHITE', 'Xiaomi Power Bank 3 - Белый', 'Белый', '20000 мАч', 3990.00, true),
((SELECT id FROM products WHERE sku = 'RAVPOWER-20000'), 'RAVPOWER-20000-BLACK', 'RAVPower 20000 - Черный', 'Черный', '20000 мАч', 3490.00, true),

-- Варианты подставок (цвета)
((SELECT id FROM products WHERE sku = 'LAMICALL-ADJUSTABLE'), 'LAMICALL-ADJUSTABLE-SILVER', 'Lamicall Adjustable Stand - Серебристый', 'Серебристый', 'One Size', 1590.00, true),
((SELECT id FROM products WHERE sku = 'LAMICALL-ADJUSTABLE'), 'LAMICALL-ADJUSTABLE-BLACK', 'Lamicall Adjustable Stand - Черный', 'Черный', 'One Size', 1590.00, true),
((SELECT id FROM products WHERE sku = 'SPIGEN-ARCSTATION'), 'SPIGEN-ARCSTATION-BLACK', 'Spigen ArcStation - Черный', 'Черный', 'One Size', 1290.00, true),

-- Варианты стилусов
((SELECT id FROM products WHERE sku = 'APPLE-PENCIL-2'), 'APPLE-PENCIL-2-WHITE', 'Apple Pencil 2 - Белый', 'Белый', 'One Size', 8990.00, true),
((SELECT id FROM products WHERE sku = 'SAMSUNG-S-PEN'), 'SAMSUNG-S-PEN-BLACK', 'Samsung S Pen - Черный', 'Черный', 'One Size', 2990.00, true),
((SELECT id FROM products WHERE sku = 'SAMSUNG-S-PEN'), 'SAMSUNG-S-PEN-WHITE', 'Samsung S Pen - Белый', 'Белый', 'One Size', 2990.00, true),

-- Варианты клавиатур (цвета)
((SELECT id FROM products WHERE sku = 'LOGITECH-MX-KEYS'), 'LOGITECH-MX-KEYS-BLACK', 'Logitech MX Keys - Черный', 'Черный', 'One Size', 8990.00, true),
((SELECT id FROM products WHERE sku = 'LOGITECH-MX-KEYS'), 'LOGITECH-MX-KEYS-WHITE', 'Logitech MX Keys - Белый', 'Белый', 'One Size', 8990.00, true),
((SELECT id FROM products WHERE sku = 'KEYCHRON-K2'), 'KEYCHRON-K2-BLACK', 'Keychron K2 - Черный', 'Черный', 'One Size', 5990.00, true),
((SELECT id FROM products WHERE sku = 'KEYCHRON-K2'), 'KEYCHRON-K2-WHITE', 'Keychron K2 - Белый', 'Белый', 'One Size', 5990.00, true),

-- Варианты мышей (цвета)
((SELECT id FROM products WHERE sku = 'LOGITECH-MX-MASTER-3'), 'LOGITECH-MX-MASTER-3-BLACK', 'Logitech MX Master 3 - Черный', 'Черный', 'One Size', 6990.00, true),
((SELECT id FROM products WHERE sku = 'LOGITECH-MX-MASTER-3'), 'LOGITECH-MX-MASTER-3-GRAY', 'Logitech MX Master 3 - Серый', 'Серый', 'One Size', 6990.00, true),
((SELECT id FROM products WHERE sku = 'RAZER-DEATHADDER-V3'), 'RAZER-DEATHADDER-V3-BLACK', 'Razer DeathAdder V3 - Черный', 'Черный', 'One Size', 7990.00, true),

-- Варианты ковриков для мыши (размеры)
((SELECT id FROM products WHERE sku = 'STEELSERIES-QCK'), 'STEELSERIES-QCK-SMALL', 'SteelSeries QcK - Маленький', 'Черный', 'S', 1290.00, true),
((SELECT id FROM products WHERE sku = 'STEELSERIES-QCK'), 'STEELSERIES-QCK-LARGE', 'SteelSeries QcK - Большой', 'Черный', 'L', 1990.00, true),
((SELECT id FROM products WHERE sku = 'CORSAIR-MM300'), 'CORSAIR-MM300-LARGE', 'Corsair MM300 - Большой', 'Черный', 'L', 1590.00, true),

-- Варианты веб-камер
((SELECT id FROM products WHERE sku = 'LOGITECH-C920'), 'LOGITECH-C920-BLACK', 'Logitech C920 HD Pro - Черный', 'Черный', 'One Size', 6990.00, true),
((SELECT id FROM products WHERE sku = 'RAZER-KIYO-PRO'), 'RAZER-KIYO-PRO-BLACK', 'Razer Kiyo Pro - Черный', 'Черный', 'One Size', 9990.00, true),

-- Варианты микрофонов
((SELECT id FROM products WHERE sku = 'BLUE-YETI-USB'), 'BLUE-YETI-USB-BLACK', 'Blue Yeti USB - Черный', 'Черный', 'One Size', 8990.00, true),
((SELECT id FROM products WHERE sku = 'BLUE-YETI-USB'), 'BLUE-YETI-USB-SILVER', 'Blue Yeti USB - Серебристый', 'Серебристый', 'One Size', 8990.00, true),
((SELECT id FROM products WHERE sku = 'SHURE-SM7B'), 'SHURE-SM7B-BLACK', 'Shure SM7B - Черный', 'Черный', 'One Size', 24990.00, true),

-- Варианты штативов (высота)
((SELECT id FROM products WHERE sku = 'MANFROTTO-COMPACT'), 'MANFROTTO-COMPACT-SMALL', 'Manfrotto Compact Action - Маленький', 'Черный', 'S', 2990.00, true),
((SELECT id FROM products WHERE sku = 'MANFROTTO-COMPACT'), 'MANFROTTO-COMPACT-LARGE', 'Manfrotto Compact Action - Большой', 'Черный', 'L', 3990.00, true),
((SELECT id FROM products WHERE sku = 'JOBY-GRIPTIGHT-PRO'), 'JOBY-GRIPTIGHT-PRO-BLACK', 'Joby GripTight PRO - Черный', 'Черный', 'One Size', 1990.00, true),

-- Варианты освещения (размеры)
((SELECT id FROM products WHERE sku = 'NEEWER-18-RING'), 'NEEWER-18-RING-BLACK', 'Neewer 18" Ring Light - Черный', 'Черный', '18"', 3990.00, true),
((SELECT id FROM products WHERE sku = 'GODOX-SL-60W'), 'GODOX-SL-60W-BLACK', 'Godox SL-60W LED Panel - Черный', 'Черный', 'One Size', 5990.00, true),

-- Варианты карт памяти (емкость)
((SELECT id FROM products WHERE sku = 'SANDISK-EXTREME-128GB'), 'SANDISK-EXTREME-128GB', 'SanDisk Extreme Pro 128GB', 'Черный', '128GB', 2990.00, true),
((SELECT id FROM products WHERE sku = 'SANDISK-EXTREME-128GB'), 'SANDISK-EXTREME-256GB', 'SanDisk Extreme Pro 256GB', 'Черный', '256GB', 4990.00, true),
((SELECT id FROM products WHERE sku = 'SAMSUNG-EVO-256GB'), 'SAMSUNG-EVO-256GB', 'Samsung EVO Plus 256GB', 'Черный', '256GB', 3990.00, true),
((SELECT id FROM products WHERE sku = 'SAMSUNG-EVO-256GB'), 'SAMSUNG-EVO-512GB', 'Samsung EVO Plus 512GB', 'Черный', '512GB', 6990.00, true),

-- Варианты автомобильных держателей
((SELECT id FROM products WHERE sku = 'IOTTIE-EASY-TOUCH'), 'IOTTIE-EASY-TOUCH-BLACK', 'iOttie Easy One Touch - Черный', 'Черный', 'One Size', 2490.00, true),
((SELECT id FROM products WHERE sku = 'SCOSCHE-MAGICMOUNT'), 'SCOSCHE-MAGICMOUNT-BLACK', 'Scosche MagicMount - Черный', 'Черный', 'One Size', 1990.00, true),

-- Варианты автомобильных зарядок
((SELECT id FROM products WHERE sku = 'ANKER-POWERDRIVE-2'), 'ANKER-POWERDRIVE-2-BLACK', 'Anker PowerDrive 2 - Черный', 'Черный', 'One Size', 1290.00, true),
((SELECT id FROM products WHERE sku = 'BELKIN-WIRELESS-CAR'), 'BELKIN-WIRELESS-CAR-BLACK', 'Belkin Wireless Car Charger - Черный', 'Черный', 'One Size', 2990.00, true),

-- Варианты игровых контроллеров (цвета)
((SELECT id FROM products WHERE sku = 'XBOX-WIRELESS-CTRL'), 'XBOX-WIRELESS-CTRL-BLACK', 'Xbox Wireless Controller - Черный', 'Черный', 'One Size', 4990.00, true),
((SELECT id FROM products WHERE sku = 'XBOX-WIRELESS-CTRL'), 'XBOX-WIRELESS-CTRL-WHITE', 'Xbox Wireless Controller - Белый', 'Белый', 'One Size', 4990.00, true),
((SELECT id FROM products WHERE sku = 'SONY-DUALSENSE'), 'SONY-DUALSENSE-WHITE', 'Sony DualSense - Белый', 'Белый', 'One Size', 5990.00, true),
((SELECT id FROM products WHERE sku = 'SONY-DUALSENSE'), 'SONY-DUALSENSE-BLACK', 'Sony DualSense - Черный', 'Черный', 'One Size', 5990.00, true),

-- Варианты игровых кресел (размеры)
((SELECT id FROM products WHERE sku = 'DXRACER-FORMULA'), 'DXRACER-FORMULA-SMALL', 'DXRacer Formula Series - Маленький', 'Черный', 'S', 19990.00, true),
((SELECT id FROM products WHERE sku = 'DXRACER-FORMULA'), 'DXRACER-FORMULA-LARGE', 'DXRacer Formula Series - Большой', 'Черный', 'L', 21990.00, true),
((SELECT id FROM products WHERE sku = 'SECRETLAB-TITAN'), 'SECRETLAB-TITAN-SMALL', 'Secretlab Titan - Маленький', 'Черный', 'S', 29990.00, true),
((SELECT id FROM products WHERE sku = 'SECRETLAB-TITAN'), 'SECRETLAB-TITAN-LARGE', 'Secretlab Titan - Большой', 'Черный', 'L', 31990.00, true)
ON CONFLICT (sku) DO NOTHING;

-- Создание остатков товаров по складам
INSERT INTO warehouse_stocks (warehouse_id, product_variant_id, stock, reserved_stock) VALUES 
-- Остатки для главного склада
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'APPLE-CASE-IP15P-BLUE'), 50, 5),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'APPLE-CASE-IP15P-BLACK'), 30, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'APPLE-CASE-IP15P-WHITE'), 20, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SPIGEN-UH-IP15-CLEAR'), 40, 3),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'JBL-TUNE-760NC-BLACK'), 25, 2),

-- Остатки для филиала "Центр"
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-CASE-IP15P-BLUE'), 15, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-CASE-IP15P-BLACK'), 20, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'SPIGEN-UH-IP15-CLEAR'), 25, 2),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'JBL-TUNE-760NC-BLACK'), 10, 1),

-- Остатки для филиала "Север"
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-CASE-IP15P-WHITE'), 15, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SPIGEN-UH-IP15-BLACK'), 20, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'JBL-TUNE-760NC-BLUE'), 8, 0),

-- Остатки для филиала "Юг" (Санкт-Петербург)
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-CASE-IP15P-BLUE'), 10, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-CASE-IP15P-BLACK'), 12, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'SPIGEN-UH-IP15-CLEAR'), 15, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'JBL-TUNE-760NC-BLACK'), 5, 0),

-- Остатки для новых товаров на главном складе
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'OTTERBOX-DEF-IP15-BLACK'), 30, 3),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'OTTERBOX-DEF-IP15-BLUE'), 25, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'CASETIFY-IP15PM-CLEAR'), 20, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'CASETIFY-IP15PM-BLACK'), 18, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SONY-WH-1000XM5-BLACK'), 8, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SONY-WH-1000XM5-SILVER'), 6, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SENNHEISER-HD-660S-BLACK'), 5, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'ESR-SP-IP15P-CLEAR'), 50, 5),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'ZAGG-SP-S24-CLEAR'), 40, 3),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-15W-WC-BLACK'), 20, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-15W-WC-WHITE'), 15, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'BELKIN-10W-WC-BLACK'), 25, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'APPLE-LIGHTNING-1M'), 80, 5),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'APPLE-LIGHTNING-2M'), 60, 3),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'ANKER-PL3-USB-C-1M'), 50, 4),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'ANKER-PL3-USB-C-2M'), 40, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'XIAOMI-PB-20000-BLACK'), 15, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'XIAOMI-PB-20000-WHITE'), 12, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'RAVPOWER-20000-BLACK'), 10, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'LAMICALL-ADJUSTABLE-SILVER'), 25, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'LAMICALL-ADJUSTABLE-BLACK'), 20, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SPIGEN-ARCSTATION-BLACK'), 30, 3),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'APPLE-PENCIL-2-WHITE'), 8, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-S-PEN-BLACK'), 12, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-S-PEN-WHITE'), 10, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-KEYS-BLACK'), 6, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-KEYS-WHITE'), 5, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'KEYCHRON-K2-BLACK'), 4, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'KEYCHRON-K2-WHITE'), 3, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-MASTER-3-BLACK'), 8, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-MASTER-3-GRAY'), 6, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'RAZER-DEATHADDER-V3-BLACK'), 7, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'STEELSERIES-QCK-SMALL'), 20, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'STEELSERIES-QCK-LARGE'), 15, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'CORSAIR-MM300-LARGE'), 12, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-C920-BLACK'), 5, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'RAZER-KIYO-PRO-BLACK'), 3, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'BLUE-YETI-USB-BLACK'), 4, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'BLUE-YETI-USB-SILVER'), 3, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SHURE-SM7B-BLACK'), 2, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'MANFROTTO-COMPACT-SMALL'), 10, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'MANFROTTO-COMPACT-LARGE'), 8, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'JOBY-GRIPTIGHT-PRO-BLACK'), 12, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'NEEWER-18-RING-BLACK'), 6, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'GODOX-SL-60W-BLACK'), 4, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SANDISK-EXTREME-128GB'), 25, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SANDISK-EXTREME-256GB'), 15, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-EVO-256GB'), 20, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-EVO-512GB'), 10, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'IOTTIE-EASY-TOUCH-BLACK'), 15, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SCOSCHE-MAGICMOUNT-BLACK'), 18, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'ANKER-POWERDRIVE-2-BLACK'), 25, 2),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'BELKIN-WIRELESS-CAR-BLACK'), 10, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'XBOX-WIRELESS-CTRL-BLACK'), 8, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'XBOX-WIRELESS-CTRL-WHITE'), 6, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SONY-DUALSENSE-WHITE'), 7, 1),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SONY-DUALSENSE-BLACK'), 5, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'DXRACER-FORMULA-SMALL'), 3, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'DXRACER-FORMULA-LARGE'), 2, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SECRETLAB-TITAN-SMALL'), 2, 0),
((SELECT id FROM warehouses WHERE is_main = true), (SELECT id FROM product_variants WHERE sku = 'SECRETLAB-TITAN-LARGE'), 1, 0),

-- Остатки для филиала "Центр" (новые товары)
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'OTTERBOX-DEF-IP15-BLACK'), 10, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'CASETIFY-IP15PM-CLEAR'), 8, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'SONY-WH-1000XM5-BLACK'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'ESR-SP-IP15P-CLEAR'), 20, 2),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-15W-WC-BLACK'), 8, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-LIGHTNING-1M'), 30, 2),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'XIAOMI-PB-20000-BLACK'), 6, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'LAMICALL-ADJUSTABLE-SILVER'), 10, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-PENCIL-2-WHITE'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-KEYS-BLACK'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-MASTER-3-BLACK'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'STEELSERIES-QCK-SMALL'), 8, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-C920-BLACK'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'BLUE-YETI-USB-BLACK'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'MANFROTTO-COMPACT-SMALL'), 5, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'SANDISK-EXTREME-128GB'), 10, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'IOTTIE-EASY-TOUCH-BLACK'), 6, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'ANKER-POWERDRIVE-2-BLACK'), 10, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'XBOX-WIRELESS-CTRL-BLACK'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Центр"'), (SELECT id FROM product_variants WHERE sku = 'SONY-DUALSENSE-WHITE'), 2, 0),

-- Остатки для филиала "Север" (новые товары)
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'OTTERBOX-DEF-IP15-BLUE'), 8, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'CASETIFY-IP15PM-BLACK'), 6, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SONY-WH-1000XM5-SILVER'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'ZAGG-SP-S24-CLEAR'), 15, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-15W-WC-WHITE'), 5, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'BELKIN-10W-WC-BLACK'), 8, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-LIGHTNING-2M'), 20, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'ANKER-PL3-USB-C-1M'), 15, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'XIAOMI-PB-20000-WHITE'), 4, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'LAMICALL-ADJUSTABLE-BLACK'), 8, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SPIGEN-ARCSTATION-BLACK'), 12, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-S-PEN-BLACK'), 5, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'KEYCHRON-K2-BLACK'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-MASTER-3-GRAY'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'RAZER-DEATHADDER-V3-BLACK'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'STEELSERIES-QCK-LARGE'), 6, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'CORSAIR-MM300-LARGE'), 5, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'RAZER-KIYO-PRO-BLACK'), 1, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'BLUE-YETI-USB-SILVER'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'MANFROTTO-COMPACT-LARGE'), 4, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'JOBY-GRIPTIGHT-PRO-BLACK'), 6, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'NEEWER-18-RING-BLACK'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SANDISK-EXTREME-256GB'), 8, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-EVO-256GB'), 10, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SCOSCHE-MAGICMOUNT-BLACK'), 8, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'BELKIN-WIRELESS-CAR-BLACK'), 4, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'XBOX-WIRELESS-CTRL-WHITE'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Север"'), (SELECT id FROM product_variants WHERE sku = 'SONY-DUALSENSE-BLACK'), 2, 0),

-- Остатки для филиала "Юг" (новые товары)
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'OTTERBOX-DEF-IP15-BLACK'), 5, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'CASETIFY-IP15PM-CLEAR'), 4, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'SONY-WH-1000XM5-BLACK'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'ESR-SP-IP15P-CLEAR'), 12, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-15W-WC-BLACK'), 6, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'APPLE-LIGHTNING-1M'), 15, 1),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'XIAOMI-PB-20000-BLACK'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'LAMICALL-ADJUSTABLE-SILVER'), 5, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'SAMSUNG-S-PEN-WHITE'), 4, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-KEYS-WHITE'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-MX-MASTER-3-BLACK'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'STEELSERIES-QCK-SMALL'), 6, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'LOGITECH-C920-BLACK'), 1, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'BLUE-YETI-USB-BLACK'), 1, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'MANFROTTO-COMPACT-SMALL'), 3, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'SANDISK-EXTREME-128GB'), 8, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'IOTTIE-EASY-TOUCH-BLACK'), 4, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'ANKER-POWERDRIVE-2-BLACK'), 6, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'XBOX-WIRELESS-CTRL-BLACK'), 2, 0),
((SELECT id FROM warehouses WHERE name = 'Филиал "Юг"'), (SELECT id FROM product_variants WHERE sku = 'SONY-DUALSENSE-WHITE'), 1, 0)
ON CONFLICT (warehouse_id, product_variant_id) DO NOTHING;

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

-- Обновляем несколько товаров как особенные (для примера)
UPDATE products 
SET feature = true 
WHERE sku IN (
    'APPLE-AIRPODS-PRO-2',
    'SONY-WH-1000XM5', 
    'SENNHEISER-HD-660S',
    'SECRETLAB-TITAN',
    'DXRACER-FORMULA'
);
