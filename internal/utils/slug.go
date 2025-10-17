package utils

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// GenerateSlug создает URL-friendly slug из строки
func GenerateSlug(text string) string {
	// Транслитерация кириллицы в латиницу
	text = transliterate(text)
	
	// Приводим к нижнему регистру
	text = strings.ToLower(text)
	
	// Заменяем пробелы и специальные символы на дефисы
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	text = reg.ReplaceAllString(text, "-")
	
	// Убираем дефисы в начале и конце
	text = strings.Trim(text, "-")
	
	// Ограничиваем длину
	if len(text) > 100 {
		text = text[:100]
		text = strings.TrimSuffix(text, "-")
	}
	
	// Если slug пустой, генерируем случайный
	if text == "" {
		text = "product-" + uuid.New().String()[:8]
	}
	
	return text
}

// GenerateUniqueSlug создает уникальный slug, добавляя суффикс если нужно
func GenerateUniqueSlug(baseSlug string, checkUnique func(string) bool) string {
	slug := baseSlug
	counter := 1
	
	for !checkUnique(slug) {
		slug = baseSlug + "-" + string(rune('0'+counter))
		counter++
	}
	
	return slug
}

// transliterate транслитерирует кириллицу в латиницу
func transliterate(text string) string {
	transliterationMap := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
		'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo",
		'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M",
		'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U",
		'Ф': "F", 'Х': "H", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch",
		'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
	}
	
	var result strings.Builder
	for _, r := range text {
		if replacement, exists := transliterationMap[r]; exists {
			result.WriteString(replacement)
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}
	
	return result.String()
}

// ValidateSlug проверяет валидность slug
func ValidateSlug(slug string) bool {
	if slug == "" {
		return false
	}
	
	// Проверяем, что slug содержит только латинские буквы, цифры и дефисы
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	if !matched {
		return false
	}
	
	// Проверяем, что slug не начинается и не заканчивается дефисом
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return false
	}
	
	// Проверяем длину
	if len(slug) < 3 || len(slug) > 100 {
		return false
	}
	
	return true
}
