package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/eikuma/stockle/backend/internal/middleware"
	"github.com/eikuma/stockle/backend/internal/services"
)

type AuthController struct {
	authService *services.AuthService
	validator   *validator.Validate
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validator:   validator.New(),
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_REQUEST",
		})
		return
	}

	if err := ac.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	user, err := ac.authService.Register(req)
	if err != nil {
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
				"code":  "EMAIL_ALREADY_EXISTS",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
				"code":  "USER_CREATION_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user.ToResponse(),
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_REQUEST",
		})
		return
	}

	if err := ac.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	tokens, user, err := ac.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
			"code":  "INVALID_CREDENTIALS",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
		"user":   user.ToResponse(),
	})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req services.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_REQUEST",
		})
		return
	}

	if err := ac.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	tokens, err := ac.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
			"code":  "INVALID_REFRESH_TOKEN",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func (ac *AuthController) Logout(c *gin.Context) {
	var req services.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_REQUEST",
		})
		return
	}

	if err := ac.authService.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "LOGOUT_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ac *AuthController) GoogleAuth(c *gin.Context) {
	fmt.Println("🔥 [GoogleAuth] Request received from:", c.ClientIP())
	fmt.Println("🔥 [GoogleAuth] Request headers:", c.Request.Header)
	
	var req services.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("🔥 [GoogleAuth] JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_REQUEST",
		})
		return
	}

	fmt.Printf("🔥 [GoogleAuth] Request data: %+v\n", req)

	if err := ac.validator.Struct(req); err != nil {
		fmt.Printf("🔥 [GoogleAuth] Validation error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	fmt.Println("🔥 [GoogleAuth] Calling AuthService.GoogleAuth...")
	tokens, user, err := ac.authService.GoogleAuth(req)
	if err != nil {
		fmt.Printf("🔥 [GoogleAuth] AuthService error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to authenticate with Google",
			"code":  "GOOGLE_AUTH_ERROR",
		})
		return
	}

	fmt.Printf("🔥 [GoogleAuth] Success! User: %+v\n", user.ToResponse())
	fmt.Printf("🔥 [GoogleAuth] Tokens generated: hasAccessToken=%t, hasRefreshToken=%t\n", 
		tokens.AccessToken != "", tokens.RefreshToken != "")

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
		"user":   user.ToResponse(),
	})
}

func (ac *AuthController) Me(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
		return
	}

	userEmail, _ := middleware.GetUserEmail(c)

	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"user_email": userEmail,
		"message":    "Authenticated user",
	})
}