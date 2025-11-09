package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "context"
    "time"
    
    "github.com/KotovBoris/AutoSave/backend/internal/banks"
    "github.com/KotovBoris/AutoSave/backend/internal/config"
    "github.com/KotovBoris/AutoSave/backend/internal/handlers"
    "github.com/KotovBoris/AutoSave/backend/internal/repository"
    "github.com/KotovBoris/AutoSave/backend/internal/router"
    "github.com/KotovBoris/AutoSave/backend/internal/services"
    "github.com/KotovBoris/AutoSave/backend/pkg/database"
    "github.com/KotovBoris/AutoSave/backend/pkg/jwt"
    "github.com/KotovBoris/AutoSave/backend/pkg/logger"
    
    "github.com/gin-gonic/gin"
)

func main() {
    // Load config
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
        os.Exit(1)
    }

    // Setup logger
    log := logger.New(cfg.GetLogLevel(), cfg.LogFormat)
    log.Info().Msg("Starting AutoSave Backend...")
    log.Info().Str("env", cfg.AppEnv).Msg("Environment loaded")

    // Setup Gin mode
    if cfg.IsProduction() {
        gin.SetMode(gin.ReleaseMode)
    }

    // Setup database connection
    dbConfig := database.Config{
        Host:     cfg.DBHost,
        Port:     cfg.DBPort,
        User:     cfg.DBUser,
        Password: cfg.DBPassword,
        Database: cfg.DBName,
        SSLMode:  cfg.DBSSLMode,
    }
    db, err := database.NewPostgresDB(dbConfig, log.Logger)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to database")
    }
    defer db.Close()

    // Initialize repositories
    repos := repository.NewRepositories(db.DB)
    log.Info().Msg("Repositories initialized")

    // Initialize utilities
    jwtUtil := jwt.NewJWTUtil(cfg.JWTSecret, cfg.JWTExpiry)
    log.Info().Msg("JWT utility initialized")

    // Initialize bank factory
    bankFactory := banks.NewFactory(cfg, log.Logger)
    log.Info().Msg("Bank factory initialized")

    // Initialize services
    authService := services.NewAuthService(repos.User, jwtUtil, log.Logger)
    bankService := services.NewBankService(repos.Bank, repos.Account, repos.Transaction, bankFactory, log.Logger)
    accountService := services.NewAccountService(repos.Account, repos.Transaction, log.Logger)
    analysisService := services.NewAnalysisService(repos.User, repos.Transaction, log.Logger)
    goalService := services.NewGoalService(repos.Goal, repos.Deposit, repos.User, repos.Bank, log.Logger)
    log.Info().Msg("Services initialized")

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService)
    bankHandler := handlers.NewBankHandler(bankService)
    accountHandler := handlers.NewAccountHandler(accountService)
    analysisHandler := handlers.NewAnalysisHandler(analysisService)
    goalHandler := handlers.NewGoalHandler(goalService)
    log.Info().Msg("Handlers initialized")

    // Setup router
    appRouter := router.NewRouter(
        authHandler,
        bankHandler,
        accountHandler,
        analysisHandler,
        goalHandler,
        jwtUtil,
        log.Logger,
        cfg.CORSAllowedOrigins,
    )
    engine := appRouter.Setup()
    log.Info().Msg("Router initialized")

    // Setup server
    server := &http.Server{
        Addr:    ":" + cfg.AppPort,
        Handler: engine,
    }

    // Graceful shutdown
    go func() {
        log.Info().Str("port", cfg.AppPort).Msg("Server is listening")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal().Err(err).Msg("Server failed to start")
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Warn().Msg("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal().Err(err).Msg("Server forced to shutdown")
    }

    log.Info().Msg("Server exiting")
}

