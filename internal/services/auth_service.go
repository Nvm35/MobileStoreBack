package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"mobile-store-back/internal/config"
	"mobile-store-back/internal/models"
	"mobile-store-back/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	ErrTokenExpired          = &TokenError{Type: "expired", Message: "Token has expired"}
	ErrTokenInvalid          = &TokenError{Type: "invalid", Message: "Invalid token"}
	ErrTokenMalformed        = &TokenError{Type: "malformed", Message: "Malformed token"}
	ErrRefreshTokenMalformed = errors.New("refresh token malformed")
	ErrRefreshTokenInvalid   = errors.New("refresh token invalid")
	ErrRefreshTokenExpired   = errors.New("refresh token expired")
	ErrRefreshSessionRevoked = errors.New("refresh session revoked")
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
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
}

type SessionMetadata struct {
	UserAgent string
	IPAddress string
}

type AuthResponse struct {
	Token            string      `json:"token"`
	User             models.User `json:"user"`
	RefreshToken     string      `json:"-"`
	RefreshExpiresAt time.Time   `json:"-"`
}

func (s *AuthService) Register(req *RegisterRequest, meta *SessionMetadata) (*AuthResponse, error) {
	existingUser, _ := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		IsActive:  true,
		Role:      "customer",
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return s.issueTokens(user, meta)
}

func (s *AuthService) Login(req *LoginRequest, meta *SessionMetadata) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	now := time.Now()
	user.LastLogin = &now
	if err := s.repo.UpdateUser(user); err != nil {
		// Ошибку обновления логина игнорируем, чтобы не блокировать вход
	}

	return s.issueTokens(user, meta)
}

func (s *AuthService) RefreshSession(refreshToken string, meta *SessionMetadata) (*AuthResponse, error) {
	sessionID, secret, err := parseRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrRefreshTokenMalformed
	}

	session, err := s.repo.GetSessionByID(sessionID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenInvalid
		}
		return nil, err
	}

	if session.RevokedAt != nil {
		return nil, ErrRefreshSessionRevoked
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	if hashTokenSecret(secret) != session.TokenHash {
		return nil, ErrRefreshTokenInvalid
	}

	user, err := s.repo.GetUserByID(session.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	newSecret, err := generateRefreshSecret()
	if err != nil {
		return nil, err
	}

	session.TokenHash = hashTokenSecret(newSecret)
	session.ExpiresAt = time.Now().Add(s.refreshTokenDuration())
	if meta != nil {
		session.UserAgent = meta.UserAgent
		session.IPAddress = meta.IPAddress
	}

	if err := s.repo.UpdateSession(session); err != nil {
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:            accessToken,
		User:             *user,
		RefreshToken:     composeRefreshToken(session.ID, newSecret),
		RefreshExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	sessionID, secret, err := parseRefreshToken(refreshToken)
	if err != nil {
		return nil
	}

	session, err := s.repo.GetSessionByID(sessionID.String())
	if err != nil {
		return nil
	}

	if hashTokenSecret(secret) != session.TokenHash {
		return nil
	}

	return s.repo.DeleteSessionByID(sessionID.String())
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return "", ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return "", ErrTokenInvalid
		}
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

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *AuthService) issueTokens(user *models.User, meta *SessionMetadata) (*AuthResponse, error) {
	accessToken, err := s.generateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := s.createSession(user.ID, meta)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:            accessToken,
		User:             *user,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) createSession(userID uuid.UUID, meta *SessionMetadata) (string, time.Time, error) {
	// Ограничиваем количество активных сессий пользователя (оставляем последние 4 + новая = 5)
	// Ошибку логируем, но не блокируем создание новой сессии
	if err := s.repo.DeleteOldSessionsForUser(userID.String(), 4); err != nil {
		// TODO: добавить логирование
		// log.Printf("Failed to cleanup old sessions for user %s: %v", userID, err)
	}

	secret, err := generateRefreshSecret()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(s.refreshTokenDuration())
	session := &models.Session{
		UserID:    userID,
		TokenHash: hashTokenSecret(secret),
		ExpiresAt: expiresAt,
	}

	if meta != nil {
		session.UserAgent = meta.UserAgent
		session.IPAddress = meta.IPAddress
	}

	if err := s.repo.CreateSession(session); err != nil {
		return "", time.Time{}, err
	}

	return composeRefreshToken(session.ID, secret), expiresAt, nil
}

func (s *AuthService) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.accessTokenDuration()).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *AuthService) accessTokenDuration() time.Duration {
	minutes := s.cfg.JWT.AccessTokenMinutes
	if minutes <= 0 {
		minutes = 15
	}
	return time.Duration(minutes) * time.Minute
}

func (s *AuthService) refreshTokenDuration() time.Duration {
	days := s.cfg.Auth.RefreshTokenDays
	if days <= 0 {
		days = 30
	}
	return time.Duration(days) * 24 * time.Hour
}

func generateRefreshSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hashTokenSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

func composeRefreshToken(sessionID uuid.UUID, secret string) string {
	return sessionID.String() + "." + secret
}

func parseRefreshToken(raw string) (uuid.UUID, string, error) {
	parts := strings.SplitN(raw, ".", 2)
	if len(parts) != 2 {
		return uuid.Nil, "", ErrRefreshTokenMalformed
	}

	sessionID, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, "", ErrRefreshTokenMalformed
	}

	secret := parts[1]
	if secret == "" {
		return uuid.Nil, "", ErrRefreshTokenMalformed
	}

	return sessionID, secret, nil
}
