package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/autosave/backend/internal/middleware"
    "github.com/autosave/backend/internal/models"
    "github.com/autosave/backend/internal/services"
    "github.com/autosave/backend/pkg/validator"
    "github.com/gin-gonic/gin"
)

type GoalHandler struct {
    goalService *services.GoalService
}

func NewGoalHandler(goalService *services.GoalService) *GoalHandler {
    return &GoalHandler{
        goalService: goalService,
    }
}

func (h *GoalHandler) GetGoals(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    goals, err := h.goalService.GetUserGoals(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, goals)
}

func (h *GoalHandler) CreateGoal(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    var req models.CreateGoalRequest
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
    
    goal, err := h.goalService.CreateGoal(c.Request.Context(), userID, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "CREATE_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusCreated, goal)
}

func (h *GoalHandler) UpdateGoal(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    goalID, err := strconv.Atoi(c.Param("goalId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid goal ID",
            },
        })
        return
    }
    
    var req models.UpdateGoalRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid request body",
            },
        })
        return
    }
    
    if err := h.goalService.UpdateGoal(c.Request.Context(), userID, goalID, req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "UPDATE_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Goal updated",
    })
}

func (h *GoalHandler) DeleteGoal(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    goalID, err := strconv.Atoi(c.Param("goalId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid goal ID",
            },
        })
        return
    }
    
    if err := h.goalService.DeleteGoal(c.Request.Context(), userID, goalID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "DELETE_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusNoContent, nil)
}

func (h *GoalHandler) ReorderGoals(c *gin.Context) {
    userID, _ := middleware.GetUserID(c)
    
    var req models.ReorderGoalsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Invalid request body",
            },
        })
        return
    }
    
    if err := h.goalService.ReorderGoals(c.Request.Context(), userID, req.GoalIDs); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "REORDER_FAILED",
                "message": err.Error(),
            },
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Goals reordered",
    })
}
