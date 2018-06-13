package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/dzhenquan/golangboy/config"
)

// RefreshTokenCookie 刷新过期时间
func RefreshTokenCookie(c *gin.Context) {
    tokenString, err := c.Cookie("token")

    if tokenString != "" && err == nil {
        cookie := &http.Cookie{
            Name: "token",
            Value: tokenString,
            MaxAge: config.ServerConfig.TokenMaxAge,
            Path: "/",
            Domain: "",
            HttpOnly: true,
        }
        http.SetCookie(c.Writer, cookie)
    }
    c.Next()
}
