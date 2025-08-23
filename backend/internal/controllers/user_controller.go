package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/eikuma/stockle/backend/internal/middleware"
	"github.com/eikuma/stockle/backend/internal/models"
	"github.com/eikuma/stockle/backend/internal/repositories"
)

type UserController struct {
	userRepo  repositories.UserRepository
	validator *validator.Validate
}

func NewUserController(userRepo repositories.UserRepository) *UserController {
	return &UserController{
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// GetProfile - GET /users/me
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
		return
	}

	// string to uint conversion
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	user, err := uc.userRepo.GetByID(uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
			"code":  "USER_NOT_FOUND",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateProfile - PUT /users/me
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
		return
	}

	// string to uint conversion
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_REQUEST",
		})
		return
	}

	if err := uc.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	user, err := uc.userRepo.GetByID(uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
			"code":  "USER_NOT_FOUND",
		})
		return
	}

	// Update user fields
	if req.Name != "" {
		user.Name = req.Name
		user.DisplayName = req.Name // sync display name with name
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := uc.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
			"code":  "UPDATE_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}