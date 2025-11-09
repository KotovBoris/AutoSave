package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
)

func LoggerMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery
        
        c.Next()
        
        latency := time.Since(start)
        statusCode := c.Writer.Status()
        clientIP := c.ClientIP()
        method := c.Request.Method
        
        if raw != "" {
            path = path + "?" + raw
        }
        
        logEvent := logger.Info()
        
        if statusCode >= 500 {
            logEvent = logger.Error()
        } else if statusCode >= 400 {
            logEvent = logger.Warn()
        }
        
        logEvent.
            Str("method", method).
            Str("path", path).
            Int("status", statusCode).
            Dur("latency", latency).
            Str("ip", clientIP).
            Msg("HTTP request")
    }
}
