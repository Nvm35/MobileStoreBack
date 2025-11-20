package services

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"mobile-store-back/internal/config"
)

type CloudinaryService struct {
	config *config.CloudinaryConfig
	client *http.Client
}

func NewCloudinaryService(cfg *config.CloudinaryConfig) *CloudinaryService {
	return &CloudinaryService{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetConfig возвращает конфигурацию (для отладки)
func (s *CloudinaryService) GetConfig() *config.CloudinaryConfig {
	return s.config
}

// validateConfig проверяет, что конфигурация Cloudinary валидна
func (s *CloudinaryService) validateConfig() error {
	if s.config == nil {
		return fmt.Errorf("cloudinary config is nil")
	}
	if s.config.CloudName == "" {
		return fmt.Errorf("cloudinary cloud_name is not set")
	}
	if s.config.APIKey == "" {
		return fmt.Errorf("cloudinary api_key is not set")
	}
	if s.config.APISecret == "" {
		return fmt.Errorf("cloudinary api_secret is not set")
	}
	return nil
}

// CloudinaryImage представляет изображение из Cloudinary
type CloudinaryImage struct {
	PublicID     string    `json:"public_id"`
	Format       string    `json:"format"`
	Version      int       `json:"version"`
	ResourceType string    `json:"resource_type"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
	Bytes        int       `json:"bytes"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	URL          string    `json:"url"`
	SecureURL    string    `json:"secure_url"`
	Folder       string    `json:"folder,omitempty"`
}

// CloudinaryFolder представляет папку в Cloudinary
type CloudinaryFolder struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
}

// CloudinaryListResponse ответ от API Cloudinary для списка ресурсов
type CloudinaryListResponse struct {
	Resources []CloudinaryImage `json:"resources"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

// generateSignature генерирует подпись для Cloudinary Admin API
// Формат: отсортированные параметры (кроме api_key, file, cloud_name, resource_type) + api_secret
// Затем SHA1 хеш
func (s *CloudinaryService) generateSignature(params map[string]string) string {
	// Исключаемые ключи из подписи
	excludedKeys := map[string]bool{
		"api_key":      true,
		"file":         true,
		"cloud_name":   true,
		"resource_type": true,
		"signature":    true,
	}

	// Собираем ключи для подписи
	keys := make([]string, 0, len(params))
	for k := range params {
		if !excludedKeys[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// Формируем строку для подписи: key1=value1&key2=value2&...&api_secret
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signString := strings.Join(parts, "&") + s.config.APISecret

	// Вычисляем SHA1
	hash := sha1.Sum([]byte(signString))
	return fmt.Sprintf("%x", hash)
}

// GetImages получает список изображений из Cloudinary Admin API
// Использует Basic Authentication - простой и надежный способ
func (s *CloudinaryService) GetImages(folder string, maxResults int, nextCursor string) (*CloudinaryListResponse, error) {
	if err := s.validateConfig(); err != nil {
		return nil, fmt.Errorf("cloudinary configuration error: %w", err)
	}
	// Формируем query параметры
	query := url.Values{}
	if maxResults > 0 {
		query.Add("max_results", strconv.Itoa(maxResults))
	}
	if nextCursor != "" {
		query.Add("next_cursor", nextCursor)
	}
	if folder != "" {
		query.Add("prefix", folder)
	}

	// Формируем URL
	apiURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/resources/image/upload", s.config.CloudName)
	if len(query) > 0 {
		apiURL += "?" + query.Encode()
	}
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Используем Basic Authentication - простой и надежный способ
	req.SetBasicAuth(s.config.APIKey, s.config.APISecret)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("cloudinary API error: %s - %s", resp.Status, string(body))
		if cldError := resp.Header.Get("x-cld-error"); cldError != "" {
			errorMsg += fmt.Sprintf(" | x-cld-error: %s", cldError)
		}
		return nil, fmt.Errorf(errorMsg)
	}

	var result CloudinaryListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for i := range result.Resources {
		result.Resources[i].URL = s.buildImageURL(result.Resources[i].PublicID, result.Resources[i].Format)
		result.Resources[i].SecureURL = s.buildSecureImageURL(result.Resources[i].PublicID, result.Resources[i].Format)
	}

	return &result, nil
}

// GetFolders получает список папок из Cloudinary
// Извлекаем папки из списка ресурсов
func (s *CloudinaryService) GetFolders() ([]CloudinaryFolder, error) {
	if err := s.validateConfig(); err != nil {
		return nil, fmt.Errorf("cloudinary configuration error: %w", err)
	}

	// Получаем все ресурсы
	images, err := s.GetImages("", 500, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	// Извлекаем уникальные папки из public_id
	folderMap := make(map[string]CloudinaryFolder)
	
	for _, img := range images.Resources {
		if img.PublicID != "" {
			parts := strings.Split(img.PublicID, "/")
			if len(parts) > 1 {
				currentPath := ""
				for i := 0; i < len(parts)-1; i++ {
					if currentPath == "" {
						currentPath = parts[i]
					} else {
						currentPath = currentPath + "/" + parts[i]
					}
					
					if _, exists := folderMap[currentPath]; !exists {
						folderMap[currentPath] = CloudinaryFolder{
							Name:     parts[i],
							Path:     currentPath,
							FullPath: currentPath,
						}
					}
				}
			}
		}
	}

	// Преобразуем map в slice и сортируем
	folders := make([]CloudinaryFolder, 0, len(folderMap))
	for _, folder := range folderMap {
		folders = append(folders, folder)
	}
	sort.Slice(folders, func(i, j int) bool {
		return folders[i].FullPath < folders[j].FullPath
	})

	return folders, nil
}

// UploadImage загружает изображение в Cloudinary (Upload API)
func (s *CloudinaryService) UploadImage(file io.Reader, fileName string, folder string) (*CloudinaryImage, error) {
	if err := s.validateConfig(); err != nil {
		return nil, fmt.Errorf("cloudinary configuration error: %w", err)
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params := map[string]string{
		"api_key":   s.config.APIKey,
		"timestamp": timestamp,
	}

	if folder != "" {
		params["folder"] = folder
	}

	signature := s.generateSignature(params)
	params["signature"] = signature

	// Формируем multipart/form-data запрос
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(fileData); err != nil {
		return nil, err
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", s.config.CloudName)
	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("cloudinary API error: %s - %s", resp.Status, string(bodyBytes))
	}

	var result CloudinaryImage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	result.URL = s.buildImageURL(result.PublicID, result.Format)
	result.SecureURL = s.buildSecureImageURL(result.PublicID, result.Format)

	return &result, nil
}

// DeleteImage удаляет изображение из Cloudinary (Admin API)
// Использует Basic Authentication - простой и надежный способ
func (s *CloudinaryService) DeleteImage(publicID string) error {
	if err := s.validateConfig(); err != nil {
		return fmt.Errorf("cloudinary configuration error: %w", err)
	}

	// Формируем URL с public_id
	apiURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/resources/image/upload", s.config.CloudName)
	query := url.Values{}
	query.Add("public_ids[]", publicID)
	apiURL += "?" + query.Encode()
	
	req, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		return err
	}

	// Используем Basic Authentication
	req.SetBasicAuth(s.config.APIKey, s.config.APISecret)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cloudinary API error: %s - %s", resp.Status, string(body))
	}

	return nil
}

// buildImageURL строит URL изображения
func (s *CloudinaryService) buildImageURL(publicID, format string) string {
	return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s.%s", s.config.CloudName, publicID, format)
}

// buildSecureImageURL строит secure URL изображения
func (s *CloudinaryService) buildSecureImageURL(publicID, format string) string {
	return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s.%s", s.config.CloudName, publicID, format)
}
