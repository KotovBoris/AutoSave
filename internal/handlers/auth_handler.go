package handlers

import (
    "net/http"
    
    "github.com/autosave/backend/internal/middleware"
    "github.com/autosave/backend/internal/models"
    "github.com/autosave/backend/internal/services"
    "github.com/autosave/backend/pkg/validator"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req models.UserRegistration
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid request body",
                "details": err.Error(),
            },
        })
        return
    }
    
    if err := validator.Validate(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Validation failed",
                "details": err.Error(),
            },
        })
        return
    }
    
    user, token, err := h.authService.Register(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusConflict, gin.H{
            "error": gin.H{
                "code":    "REGISTRATION_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "user":  user.ToResponse(),
        "token": token,
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.UserLogin
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid request body",
            },
        })
        return
    }
    
    if err := validator.Validate(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Validation failed",
                "details": err.Error(),
            },
        })
        return
    }
    
    user, token, err := h.authService.Login(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "LOGIN_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "user":  user.ToResponse(),
        "token": token,
    })
}

func (h *AuthHandler) GetMe(c *gin.Context) {
    userID, exists := middleware.GetUserID(c)
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "UNAUTHORIZED",
                "message": "User not authenticated",
            },
        })
        return
    }
    
    user, err := h.authService.GetUser(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{
                "code":    "NOT_FOUND",
                "message": "User not found",
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, user.ToResponse())
}
