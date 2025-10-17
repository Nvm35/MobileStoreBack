-- =============================================
-- МИГРАЦИЯ: Добавление поля slug к товарам
-- =============================================

-- Добавляем поле slug к таблице products
ALTER TABLE products ADD COLUMN IF NOT EXISTS slug VARCHAR(255);

-- Создаем функцию для генерации slug из названия товара
CREATE OR REPLACE FUNCTION generate_slug_from_name(product_name TEXT)
RETURNS TEXT AS $$
DECLARE
    result TEXT;
BEGIN
    -- Транслитерация кириллицы в латиницу
    result := product_name;
    
    -- Заменяем кириллические символы на латинские
    result := regexp_replace(result, '[аА]', 'a', 'g');
    result := regexp_replace(result, '[бБ]', 'b', 'g');
    result := regexp_replace(result, '[вВ]', 'v', 'g');
    result := regexp_replace(result, '[гГ]', 'g', 'g');
    result := regexp_replace(result, '[дД]', 'd', 'g');
    result := regexp_replace(result, '[еЕёЁ]', 'e', 'g');
    result := regexp_replace(result, '[жЖ]', 'zh', 'g');
    result := regexp_replace(result, '[зЗ]', 'z', 'g');
    result := regexp_replace(result, '[иИ]', 'i', 'g');
    result := regexp_replace(result, '[йЙ]', 'y', 'g');
    result := regexp_replace(result, '[кК]', 'k', 'g');
    result := regexp_replace(result, '[лЛ]', 'l', 'g');
    result := regexp_replace(result, '[мМ]', 'm', 'g');
    result := regexp_replace(result, '[нН]', 'n', 'g');
    result := regexp_replace(result, '[оО]', 'o', 'g');
    result := regexp_replace(result, '[пП]', 'p', 'g');
    result := regexp_replace(result, '[рР]', 'r', 'g');
    result := regexp_replace(result, '[сС]', 's', 'g');
    result := regexp_replace(result, '[тТ]', 't', 'g');
    result := regexp_replace(result, '[уУ]', 'u', 'g');
    result := regexp_replace(result, '[фФ]', 'f', 'g');
    result := regexp_replace(result, '[хХ]', 'h', 'g');
    result := regexp_replace(result, '[цЦ]', 'ts', 'g');
    result := regexp_replace(result, '[чЧ]', 'ch', 'g');
    result := regexp_replace(result, '[шШ]', 'sh', 'g');
    result := regexp_replace(result, '[щЩ]', 'sch', 'g');
    result := regexp_replace(result, '[ъЪ]', '', 'g');
    result := regexp_replace(result, '[ыЫ]', 'y', 'g');
    result := regexp_replace(result, '[ьЬ]', '', 'g');
    result := regexp_replace(result, '[эЭ]', 'e', 'g');
    result := regexp_replace(result, '[юЮ]', 'yu', 'g');
    result := regexp_replace(result, '[яЯ]', 'ya', 'g');
    
    -- Приводим к нижнему регистру
    result := lower(result);
    
    -- Заменяем пробелы и специальные символы на дефисы
    result := regexp_replace(result, '[^a-z0-9]+', '-', 'g');
    
    -- Убираем дефисы в начале и конце
    result := trim(both '-' from result);
    
    -- Ограничиваем длину
    IF length(result) > 100 THEN
        result := left(result, 100);
        result := trim(both '-' from result);
    END IF;
    
    -- Если slug пустой, генерируем случайный
    IF result = '' THEN
        result := 'product-' || substr(md5(random()::text), 1, 8);
    END IF;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Создаем функцию для генерации уникального slug
CREATE OR REPLACE FUNCTION generate_unique_slug(base_slug TEXT, product_id UUID)
RETURNS TEXT AS $$
DECLARE
    result TEXT;
    counter INTEGER := 1;
    exists_count INTEGER;
BEGIN
    result := base_slug;
    
    LOOP
        -- Проверяем, существует ли такой slug
        SELECT COUNT(*) INTO exists_count 
        FROM products 
        WHERE slug = result AND id != product_id;
        
        -- Если slug уникален, выходим из цикла
        IF exists_count = 0 THEN
            EXIT;
        END IF;
        
        -- Добавляем суффикс
        result := base_slug || '-' || counter::text;
        counter := counter + 1;
        
        -- Защита от бесконечного цикла
        IF counter > 1000 THEN
            result := base_slug || '-' || substr(md5(random()::text), 1, 8);
            EXIT;
        END IF;
    END LOOP;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Обновляем существующие товары, добавляя slug
UPDATE products 
SET slug = generate_unique_slug(
    generate_slug_from_name(name), 
    id
)
WHERE slug IS NULL OR slug = '';

-- Добавляем ограничение NOT NULL для slug
ALTER TABLE products ALTER COLUMN slug SET NOT NULL;

-- Добавляем уникальный индекс для slug
CREATE UNIQUE INDEX IF NOT EXISTS idx_products_slug_unique ON products(slug);

-- Удаляем временные функции
DROP FUNCTION IF EXISTS generate_slug_from_name(TEXT);
DROP FUNCTION IF EXISTS generate_unique_slug(TEXT, UUID);
