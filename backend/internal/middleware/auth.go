package middleware

import (
    "net/http"
    "strings"
    
    "github.com/KotovBoris/AutoSave/backend/pkg/jwt"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtUtil *jwt.JWTUtil) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "Authorization header required",
                },
            })
            c.Abort()
            return
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "Invalid authorization header format",
                },
            })
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        claims, err := jwtUtil.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "UNAUTHORIZED",
                    "message": "Invalid or expired token",
                },
            })
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        
        c.Next()
    }
}

func GetUserID(c *gin.Context) (int, bool) {
    userID, exists := c.Get("user_id")
    if !exists {
        return 0, false
    }
    
    id, ok := userID.(int)
    return id, ok
}

