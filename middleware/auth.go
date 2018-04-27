package middleware

import (
	"github.com/gin-gonic/gin"
	"gin-blog/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"gin-blog/config"
)

func getUser(c *gin.Context) (model.User, error) {
	var user model.User
	tokenString, cookieErr := c.Cookie("token")

	if cookieErr != nil {
		return user, errors.New("未登录")
	}

	token, tokenErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ServerConfig.TokenSecret), nil
	})

	if tokenErr != nil {
		return user,errors.New("未登录")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["id"].(float64))
		err := model.DB.First(&user, "id = ?", userID).Error
		if err != nil {
			return user, errors.New("未登录")
		}
		return user,nil
	}
	return user, errors.New("未登录")
}


// SigninRequired 必须是登录用户
func SigninRequired(c *gin.Context) {
	var user model.User
	var err error

	if user, err = getUser(c); err != nil {
		err = errors.New("未登录")
	}
	c.Set("user", user)
	c.Next()
}








