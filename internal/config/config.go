package config

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/joho/godotenv"
    "github.com/rs/zerolog"
)

type Config struct {
    // App
    AppEnv     string
    AppPort    string
    AppHost    string
    JWTSecret  string
    JWTExpiry  time.Duration

    // Database
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string

    // Redis
    RedisURL string

    // Team credentials
    TeamID     string
    TeamSecret string

    // VBank
    VBankClientID     string
    VBankClientSecret string
    VBankAPIURL       string

    // ABank
    ABankClientID     string
    ABankClientSecret string
    ABankAPIURL       string

    // SBank
    SBankClientID     string
    SBankClientSecret string
    SBankAPIURL       string

    // Logging
    LogLevel  string
    LogFormat string

    // CORS
    CORSAllowedOrigins []string
}

func Load() (*Config, error) {
    // Load .env file if exists
    if err := godotenv.Load(); err != nil {
        // Not an error in production where env vars are set differently
        if !strings.Contains(err.Error(), "no such file") {
            return nil, fmt.Errorf("error loading .env file: %w", err)
        }
    }

    cfg := &Config{
        // App
        AppEnv:     getEnv("APP_ENV", "development"),
        AppPort:    getEnv("APP_PORT", "8080"),
        AppHost:    getEnv("APP_HOST", "localhost"),
        JWTSecret:  getEnv("JWT_SECRET", "your-super-secret-jwt-key"),

        // Database
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "autosave"),
        DBPassword: getEnv("DB_PASSWORD", "autosave_password"),
        DBName:     getEnv("DB_NAME", "autosave_db"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

        // Redis
        RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),

        // Team
        TeamID:     getEnv("TEAM_ID", "team242"),
        TeamSecret: getEnv("TEAM_SECRET", ""),

        // VBank
        VBankClientID:     getEnv("VBANK_CLIENT_ID", ""),
        VBankClientSecret: getEnv("VBANK_CLIENT_SECRET", ""),
        VBankAPIURL:       getEnv("VBANK_API_URL", "https://vbank.open.bankingapi.ru"),

        // ABank
        ABankClientID:     getEnv("ABANK_CLIENT_ID", ""),
        ABankClientSecret: getEnv("ABANK_CLIENT_SECRET", ""),
        ABankAPIURL:       getEnv("ABANK_API_URL", "https://abank.open.bankingapi.ru"),

        // SBank
        SBankClientID:     getEnv("SBANK_CLIENT_ID", ""),
        SBankClientSecret: getEnv("SBANK_CLIENT_SECRET", ""),
        SBankAPIURL:       getEnv("SBANK_API_URL", "https://sbank.open.bankingapi.ru"),

        // Logging
        LogLevel:  getEnv("LOG_LEVEL", "debug"),
        LogFormat: getEnv("LOG_FORMAT", "console"),
    }

    // Parse JWT expiry
    expiryStr := getEnv("JWT_EXPIRY", "168h")
    expiry, err := time.ParseDuration(expiryStr)
    if err != nil {
        return nil, fmt.Errorf("invalid JWT_EXPIRY format: %w", err)
    }
    cfg.JWTExpiry = expiry

    // Parse CORS origins
    origins := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
    cfg.CORSAllowedOrigins = strings.Split(origins, ",")

    return cfg, nil
}

func (c *Config) GetDatabaseURL() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
    )
}

func (c *Config) GetLogLevel() zerolog.Level {
    switch c.LogLevel {
    case "trace":
        return zerolog.TraceLevel
    case "debug":
        return zerolog.DebugLevel
    case "info":
        return zerolog.InfoLevel
    case "warn":
        return zerolog.WarnLevel
    case "error":
        return zerolog.ErrorLevel
    case "fatal":
        return zerolog.FatalLevel
    case "panic":
        return zerolog.PanicLevel
    default:
        return zerolog.InfoLevel
    }
}

func (c *Config) IsDevelopment() bool {
    return c.AppEnv == "development"
}

func (c *Config) IsProduction() bool {
    return c.AppEnv == "production"
}

func (c *Config) IsTest() bool {
    return c.AppEnv == "test"
}

// GetBankConfig returns config for specific bank
func (c *Config) GetBankConfig(bankID string) (clientID, clientSecret, apiURL string) {
    switch bankID {
    case "vbank":
        return c.VBankClientID, c.VBankClientSecret, c.VBankAPIURL
    case "abank":
        return c.ABankClientID, c.ABankClientSecret, c.ABankAPIURL
    case "sbank":
        return c.SBankClientID, c.SBankClientSecret, c.SBankAPIURL
    default:
        return "", "", ""
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultValue
}
