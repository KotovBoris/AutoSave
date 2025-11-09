package router

import (
    "github.com/KotovBoris/AutoSave/backend/internal/handlers"
    "github.com/KotovBoris/AutoSave/backend/internal/middleware"
    "github.com/KotovBoris/AutoSave/backend/pkg/jwt"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
)

type Router struct {
    authHandler     *handlers.AuthHandler
    bankHandler     *handlers.BankHandler
    accountHandler  *handlers.AccountHandler
    analysisHandler *handlers.AnalysisHandler
    goalHandler     *handlers.GoalHandler
    jwtUtil         *jwt.JWTUtil
    logger          *zerolog.Logger
    corsOrigins     []string
}

func NewRouter(
    authHandler *handlers.AuthHandler,
    bankHandler *handlers.BankHandler,
    accountHandler *handlers.AccountHandler,
    analysisHandler *handlers.AnalysisHandler,
    goalHandler *handlers.GoalHandler,
    jwtUtil *jwt.JWTUtil,
    logger *zerolog.Logger,
    corsOrigins []string,
) *Router {
    return &Router{
        authHandler:     authHandler,
        bankHandler:     bankHandler,
        accountHandler:  accountHandler,
        analysisHandler: analysisHandler,
        goalHandler:     goalHandler,
        jwtUtil:         jwtUtil,
        logger:          logger,
        corsOrigins:     corsOrigins,
    }
}

func (r *Router) Setup() *gin.Engine {
    router := gin.New()
    
    // Global middleware
    router.Use(gin.Recovery())
    router.Use(middleware.LoggerMiddleware(r.logger))
    router.Use(middleware.CORSMiddleware(r.corsOrigins))
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API routes
    api := router.Group("/api")
    {
        // Auth (public)
        auth := api.Group("/auth")
        {
            auth.POST("/register", r.authHandler.Register)
            auth.POST("/login", r.authHandler.Login)
            auth.GET("/me", middleware.AuthMiddleware(r.jwtUtil), r.authHandler.GetMe)
        }
        
        // Protected routes
        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(r.jwtUtil))
        {
            // Banks
            banks := protected.Group("/banks")
            {
                banks.GET("", r.bankHandler.GetBanks)
                banks.POST("/connect", r.bankHandler.ConnectBank)
                banks.GET("/connected", r.bankHandler.GetConnectedBanks)
                banks.POST("/sync", r.bankHandler.SyncBanks)
                banks.DELETE("/:bankId", r.bankHandler.DisconnectBank)
            }
            
            // Accounts
            accounts := protected.Group("/accounts")
            {
                accounts.GET("", r.accountHandler.GetAccounts)
                accounts.GET("/:accountId/transactions", r.accountHandler.GetAccountTransactions)
            }
            
            // Analysis
            analysis := protected.Group("/analysis")
            {
                analysis.POST("/detect-salaries", r.analysisHandler.DetectSalaries)
                analysis.POST("/confirm-salaries", r.analysisHandler.ConfirmSalaries)
            }
            
            // Goals
            goals := protected.Group("/goals")
            {
                goals.GET("", r.goalHandler.GetGoals)
                goals.POST("", r.goalHandler.CreateGoal)
                goals.PUT("/:goalId", r.goalHandler.UpdateGoal)
                goals.DELETE("/:goalId", r.goalHandler.DeleteGoal)
                goals.PUT("/reorder", r.goalHandler.ReorderGoals)
            }
        }
    }
    
    return router
}

