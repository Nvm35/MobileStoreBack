package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"mobile-store-back/internal/config"
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"

	"github.com/gin-gonic/gin"
)

func Register(authService *services.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req services.RegisterRequest
		if !utils.ValidateRequest(c, &req) {
			return
		}

		response, err := authService.Register(&req, sessionMetadataFromContext(c))
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		setRefreshTokenCookie(c, cfg, response.RefreshToken, response.RefreshExpiresAt)

		c.JSON(http.StatusCreated, response)
	}
}

func Login(authService *services.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req services.LoginRequest
		if !utils.ValidateRequest(c, &req) {
			return
		}

		response, err := authService.Login(&req, sessionMetadataFromContext(c))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		setRefreshTokenCookie(c, cfg, response.RefreshToken, response.RefreshExpiresAt)

		// После успешного логина корзина синхронизируется на фронте
		// Фронт должен отправить товары из localStorage на бэк после логина
		// Используется JWT токен через заголовок Authorization: Bearer <token>

		c.JSON(http.StatusOK, response)
	}
}

// Refresh обновляет JWT токен
func Refresh(authService *services.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie(cfg.Auth.RefreshCookieName)
		if err != nil || refreshToken == "" {
			clearRefreshTokenCookie(c, cfg)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":    "refresh token missing",
				"code":     "REFRESH_TOKEN_MISSING",
				"redirect": true,
			})
			return
		}

		response, err := authService.RefreshSession(refreshToken, sessionMetadataFromContext(c))
		if err != nil {
			clearRefreshTokenCookie(c, cfg)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":    err.Error(),
				"code":     mapRefreshError(err),
				"redirect": true,
			})
			return
		}

		setRefreshTokenCookie(c, cfg, response.RefreshToken, response.RefreshExpiresAt)

		c.JSON(http.StatusOK, response)
	}
}

// Logout очищает refresh cookie и инвалидирует сессию
func Logout(authService *services.AuthService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if refreshToken, err := c.Cookie(cfg.Auth.RefreshCookieName); err == nil && refreshToken != "" {
			_ = authService.Logout(refreshToken)
		}

		clearRefreshTokenCookie(c, cfg)
		c.Header("Clear-Site-Data", "\"cache\", \"storage\", \"executionContexts\"")
		c.Status(http.StatusNoContent)
	}
}

func sessionMetadataFromContext(c *gin.Context) *services.SessionMetadata {
	return &services.SessionMetadata{
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
	}
}

func setRefreshTokenCookie(c *gin.Context, cfg *config.Config, token string, expiresAt time.Time) {
	if token == "" {
		return
	}

	cookie := &http.Cookie{
		Name:     cfg.Auth.RefreshCookieName,
		Value:    token,
		Path:     cookiePath(cfg),
		HttpOnly: true,
		Secure:   cfg.Auth.RefreshCookieSecure,
		SameSite: sameSiteFromConfig(cfg.Auth.RefreshCookieSameSite),
	}

	// Устанавливаем Domain только если он указан (для локальной разработки лучше не устанавливать)
	if cfg.Auth.RefreshCookieDomain != "" {
		cookie.Domain = cfg.Auth.RefreshCookieDomain
	}

	if !expiresAt.IsZero() {
		cookie.Expires = expiresAt
		cookie.MaxAge = int(time.Until(expiresAt).Seconds())
	}

	http.SetCookie(c.Writer, cookie)
}

func clearRefreshTokenCookie(c *gin.Context, cfg *config.Config) {
	cookie := &http.Cookie{
		Name:     cfg.Auth.RefreshCookieName,
		Value:    "",
		Path:     cookiePath(cfg),
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   cfg.Auth.RefreshCookieSecure,
		SameSite: sameSiteFromConfig(cfg.Auth.RefreshCookieSameSite),
	}

	// Устанавливаем Domain только если он указан (для локальной разработки лучше не устанавливать)
	if cfg.Auth.RefreshCookieDomain != "" {
		cookie.Domain = cfg.Auth.RefreshCookieDomain
	}

	http.SetCookie(c.Writer, cookie)
}

func sameSiteFromConfig(value string) http.SameSite {
	switch strings.ToLower(value) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func cookiePath(cfg *config.Config) string {
	if cfg.Auth.RefreshCookiePath == "" {
		return "/"
	}
	return cfg.Auth.RefreshCookiePath
}

func mapRefreshError(err error) string {
	switch {
	case errors.Is(err, services.ErrRefreshTokenMalformed):
		return "REFRESH_TOKEN_MALFORMED"
	case errors.Is(err, services.ErrRefreshTokenInvalid):
		return "REFRESH_TOKEN_INVALID"
	case errors.Is(err, services.ErrRefreshTokenExpired):
		return "REFRESH_TOKEN_EXPIRED"
	case errors.Is(err, services.ErrRefreshSessionRevoked):
		return "REFRESH_SESSION_REVOKED"
	default:
		return "REFRESH_FAILED"
	}
}
