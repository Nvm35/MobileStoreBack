package handlers

import (
	"mobile-store-back/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		// TODO: Implement profile update
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
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
		// TODO: Implement user update
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	}
}

func DeleteUser(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement user deletion
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
	}
}
