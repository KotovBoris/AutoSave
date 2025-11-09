package handlers

import (
    "net/http"
    
    "github.com/KotovBoris/AutoSave/backend/internal/middleware"
    "github.com/KotovBoris/AutoSave/backend/internal/models"
    "github.com/KotovBoris/AutoSave/backend/internal/services"
    "github.com/gin-gonic/gin"
)

type BankHandler struct {
    bankService *services.BankService
}

func NewBankHandler(bankService *services.BankService) *BankHandler {
    return &BankHandler{
        bankService: bankService,
    }
}

func (h *BankHandler) GetBanks(c *gin.Context) {
    banks, err := h.bankService.GetAllBanks(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, banks)
}

func (h *BankHandler) ConnectBank(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    var req models.ConnectBankRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid request body",
            },
        })
        return
    }
    
    connection, err := h.bankService.ConnectBank(c.Request.Context(), userID, req.BankID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "CONNECTION_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, connection)
}

func (h *BankHandler) GetConnectedBanks(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    banks, err := h.bankService.GetConnectedBanks(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, banks)
}

func (h *BankHandler) SyncBanks(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    result, err := h.bankService.SyncBanks(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "SYNC_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, result)
}

func (h *BankHandler) DisconnectBank(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    bankID := c.Param("bankId")
    
    if err := h.bankService.DisconnectBank(c.Request.Context(), userID, bankID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "DISCONNECT_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusNoContent, nil)
}

