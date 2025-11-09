package handlers

import (
    "net/http"
    
    "github.com/autosave/backend/internal/middleware"
    "github.com/autosave/backend/internal/models"
    "github.com/autosave/backend/internal/services"
    "github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
    analysisService *services.AnalysisService
}

func NewAnalysisHandler(analysisService *services.AnalysisService) *AnalysisHandler {
    return &AnalysisHandler{
        analysisService: analysisService,
    }
}

func (h *AnalysisHandler) DetectSalaries(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    detections, err := h.analysisService.DetectSalaries(c.Request.Context(), userID)
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
        "detectedSalaries": detections,
    })
}

func (h *AnalysisHandler) ConfirmSalaries(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    var req models.ConfirmSalariesRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid request body",
            },
        })
        return
    }
    
    analysis, err := h.analysisService.ConfirmSalaries(c.Request.Context(), userID, req.SalaryTransactionIDs)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "CONFIRMATION_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, analysis)
}
