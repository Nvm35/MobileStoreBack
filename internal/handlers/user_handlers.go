package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetProfile(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		
		user, err := userService.GetByID(userID.(string))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func UpdateProfile(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		
		var req struct {
			FirstName *string `json:"first_name" validate:"omitempty,min=2"`
			LastName  *string `json:"last_name" validate:"omitempty,min=2"`
			Phone     *string `json:"phone" validate:"omitempty,e164"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.UpdateProfile(userID.(string), req.FirstName, req.LastName, req.Phone)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func GetUsers(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		users, err := userService.List(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func GetUser(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		user, err := userService.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req struct {
			FirstName *string `json:"first_name" validate:"omitempty,min=2"`
			LastName  *string `json:"last_name" validate:"omitempty,min=2"`
			Phone     *string `json:"phone" validate:"omitempty,e164"`
			IsActive  *bool   `json:"is_active"`
			IsAdmin   *bool   `json:"is_admin"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Валидация
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.Update(id, req.FirstName, req.LastName, req.Phone, req.IsActive, req.IsAdmin)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		err := userService.Delete(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
