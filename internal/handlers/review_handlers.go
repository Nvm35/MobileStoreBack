package handlers

import (
	"mobile-store-back/internal/services"
	"mobile-store-back/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetProductReviews - получение отзывов о товаре
func GetProductReviews(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := c.Param("slug") // Может быть как slug, так и ID

		reviews, err := reviewService.GetByProductSlugOrID(identifier)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"reviews": reviews})
	}
}

// CreateReview - создание отзыва
func CreateReview(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		
		var req struct {
			ProductID uuid.UUID `json:"product_id" validate:"required"`
			OrderID   *uuid.UUID `json:"order_id"`
			Rating    int       `json:"rating" validate:"required,min=1,max=5"`
			Title     string    `json:"title" validate:"required,min=2"`
			Comment   string    `json:"comment" validate:"required,min=10"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		var orderIDStr *string
		if req.OrderID != nil {
			s := req.OrderID.String()
			orderIDStr = &s
		}
		
		review, err := reviewService.Create(userID.(string), req.ProductID.String(), orderIDStr, req.Rating, req.Title, req.Comment)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, review)
	}
}

// UpdateReview - обновление отзыва
func UpdateReview(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := c.Get("user_id")
		
		var req struct {
			Rating  *int    `json:"rating" validate:"omitempty,min=1,max=5"`
			Title   *string `json:"title" validate:"omitempty,min=2"`
			Comment *string `json:"comment" validate:"omitempty,min=10"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		review, err := reviewService.Update(id, userID.(string), req.Rating, req.Title, req.Comment)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, review)
	}
}

// DeleteReview - удаление отзыва
func DeleteReview(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := c.Get("user_id")
		
		err := reviewService.Delete(id, userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
	}
}

// VoteReview - голосование за полезность отзыва
func VoteReview(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := c.Get("user_id")
		
		var req struct {
			Helpful bool `json:"helpful" validate:"required"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		err := reviewService.Vote(id, userID.(string), req.Helpful)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Vote recorded"})
	}
}

// GetUserReviews - получение отзывов пользователя
func GetUserReviews(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		reviews, err := reviewService.GetByUserID(userID.(string))
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"reviews": reviews})
	}
}

// GetAllReviews - получение всех отзывов (админ)
func GetAllReviews(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviews, err := reviewService.GetAll()
		utils.HandleInternalError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"reviews": reviews})
	}
}

// ApproveReview - одобрение отзыва (админ)
func ApproveReview(reviewService *services.ReviewService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		
		var req struct {
			Approved bool `json:"approved" validate:"required"`
		}

		if !utils.ValidateRequest(c, &req) {
			return
		}

		err := reviewService.Approve(id, req.Approved)
		utils.HandleError(c, err)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Review approval updated"})
	}
}
