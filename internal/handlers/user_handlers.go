package handlers

import (
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"
	"net/http"

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
		userID, _ := c.Get("user_id")
		
		var req struct {
			FirstName         *string `json:"first_name" validate:"omitempty,min=2"`
			LastName          *string `json:"last_name" validate:"omitempty,min=2"`
			Phone             *string `json:"phone" validate:"omitempty,e164"`
			AddressStreet     *string `json:"address_street"`
			AddressCity       *string `json:"address_city"`
			AddressState      *string `json:"address_state"`
			AddressPostalCode *string `json:"address_postal_code"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		user, err := userService.UpdateProfile(
			userID.(string),
			req.FirstName,
			req.LastName,
			req.Phone,
			req.AddressStreet,
			req.AddressCity,
			req.AddressState,
			req.AddressPostalCode,
		)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func GetUsers(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userService.List()
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func GetUser(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		user, err := userService.GetByID(id)
		utils.HandleNotFound(c, err, "User not found")
		if err != nil {
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
			Role      *string `json:"role" validate:"omitempty,oneof=admin manager customer"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		user, err := userService.Update(id, req.FirstName, req.LastName, req.Phone, req.IsActive, req.Role)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		err := userService.Delete(id)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
