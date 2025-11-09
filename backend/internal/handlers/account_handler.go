package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/KotovBoris/AutoSave/backend/internal/middleware"
    "github.com/KotovBoris/AutoSave/backend/internal/services"
    "github.com/gin-gonic/gin"
)

type AccountHandler struct {
    accountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
    return &AccountHandler{
        accountService: accountService,
    }
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    accounts, err := h.accountService.GetUserAccounts(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) GetAccountTransactions(c *gin.Context) {
    accountID, err := strconv.Atoi(c.Param("accountId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid account ID",
            },
        })
        return
    }
    
    limit := 50
    if limitParam := c.Query("limit"); limitParam != "" {
        if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }
    
    transactions, err := h.accountService.GetAccountTransactions(c.Request.Context(), accountID, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "accountId":    accountID,
        "transactions": transactions,
    })
}

