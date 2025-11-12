package services

import (
	"errors"
	"mobile-store-back/internal/config"
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TokenError представляет различные типы ошибок токена
type TokenError struct {
	Type    string // "expired", "invalid", "malformed"
	Message string
}

func (e *TokenError) Error() string {
	return e.Message
}

var (
	ErrTokenExpired   = &TokenError{Type: "expired", Message: "Token has expired"}
	ErrTokenInvalid   = &TokenError{Type: "invalid", Message: "Invalid token"}
	ErrTokenMalformed = &TokenError{Type: "malformed", Message: "Malformed token"}
)

type AuthService struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	FirstName string  `json:"first_name" validate:"required,min=2"`
	LastName  string  `json:"last_name" validate:"required,min=2"`
	Phone     string  `json:"phone" validate:"omitempty,e164"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Проверяем, существует ли пользователь
	existingUser, _ := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:   req.LastName,
		Phone:    req.Phone,
		IsActive: true,
		Role:     "customer",
	}
	

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Генерируем токен
	token, err := s.generateToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Получаем пользователя
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Проверяем, активен ли пользователь
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Обновляем время последнего входа
	now := time.Now()
	user.LastLogin = &now
	if err := s.repo.UpdateUser(user); err != nil {
		// Логируем ошибку, но не прерываем процесс входа
		// Можно добавить логирование здесь
	}

	// Генерируем токен
	token, err := s.generateToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		// Проверяем, истек ли токен (jwt v5)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}
		// Проверяем другие типы ошибок
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return "", ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return "", ErrTokenInvalid
		}
		// Для других ошибок возвращаем общую ошибку
		return "", ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", ErrTokenInvalid
		}
		return userID, nil
	}

	return "", ErrTokenInvalid
}

// ValidateTokenWithClaims возвращает userID и claims токена, даже если токен истек
// Используется для refresh токена
func (s *AuthService) ValidateTokenWithClaims(tokenString string) (string, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	// Даже если токен истек, мы можем извлечь claims
	var claims jwt.MapClaims
	if token != nil {
		if c, ok := token.Claims.(jwt.MapClaims); ok {
			claims = c
		}
	}

	if err != nil {
		// Если токен истек, но claims валидны, возвращаем их
		if errors.Is(err, jwt.ErrTokenExpired) || (token != nil && claims != nil) {
			if userID, ok := claims["user_id"].(string); ok {
				return userID, claims, ErrTokenExpired
			}
		}
		return "", nil, ErrTokenInvalid
	}

	if claims != nil && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", nil, ErrTokenInvalid
		}
		return userID, claims, nil
	}

	return "", nil, ErrTokenInvalid
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

// RefreshToken обновляет токен, если старый токен валиден или истек недавно (в пределах grace period)
// Grace period по умолчанию 7 дней после истечения токена
func (s *AuthService) RefreshToken(tokenString string) (*AuthResponse, error) {
	userID, claims, err := s.ValidateTokenWithClaims(tokenString)
	if err != nil && err != ErrTokenExpired {
		return nil, errors.New("invalid token")
	}

	// Проверяем grace period (7 дней после истечения)
	if err == ErrTokenExpired && claims != nil {
		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			gracePeriod := 7 * 24 * time.Hour // 7 дней
			if time.Since(expTime) > gracePeriod {
				return nil, errors.New("token expired too long ago, please login again")
			}
		}
	}

	// Проверяем, что пользователь существует и активен
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Генерируем новый токен
	newToken, err := s.generateToken(userID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: newToken,
		User:  *user,
	}, nil
}

func (s *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.cfg.JWT.ExpireHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}