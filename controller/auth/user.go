package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gin-blog/model"
	"github.com/pkg/errors"
	"fmt"
	"strings"
	"gin-blog/utils"
	"github.com/dgrijalva/jwt-go"
	"gin-blog/config"
	"strconv"
)

func SignInGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signin.html", nil)
}

func AdminSignInGet(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		fmt.Println("user:", user)

		var comments []string
		comments = append(comments, "nihao","haodehen","xiaodai")

		c.HTML(http.StatusOK, "admin/index.html", gin.H{
			"pageCount":    2,
			"postCount":    3,
			"tagCount":     1,
			"commentCount": 5,
			"user":         user,
			"comments":     comments,
		})
	}
}

func SignInPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	username = strings.TrimSpace(username)

	var err error
	if len(username) > 0 && len(password) > 0 {
		var user model.User
		err = model.DB.First(&user, "email = ?", username).Error
		if err == nil && user.Password == user.EncryptPassword(password, user.Salt()) {
			if !user.LockState {
				fmt.Println("id:", user.ID)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"id": user.ID,
				})

				var tokenString string
				tokenString, err = token.SignedString([]byte(config.ServerConfig.TokenSecret))
				if err == nil {
					cookie := &http.Cookie{
						Name: "token",
						Value: tokenString,
						MaxAge: config.ServerConfig.TokenMaxAge,
						Path: "/",
						Domain: "",
						HttpOnly: true,
					}
					http.SetCookie(c.Writer, cookie)
					//c.SetCookie("token", tokenString, config.ServerConfig.TokenMaxAge, "/", "", true, true)
					if user.IsAdmin {
						c.Redirect(http.StatusMovedPermanently, "/admin/index")
					} else {
						c.Redirect(http.StatusMovedPermanently, "/")
					}
					return
				} else {
					err = errors.New("内部错误.")
				}
			} else {
				err = errors.New("Your account have been locked.")
			}
		} else {
			err = errors.New("Invalid username or password.")
		}
	} else {
		err = errors.New("Username or password cannot be null.")
	}
	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"message": err.Error(),
	})
}


func SignUpGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", nil)
}

func SignUpPost(c *gin.Context) {
	email := c.PostForm("email")
	telephone := c.PostForm("telephone")
	password :=	c.PostForm("password")

	email = utils.AvoidXSS(email)
	email = strings.TrimSpace(email)
	telephone = utils.AvoidXSS(telephone)
	telephone = strings.TrimSpace(telephone)

	fmt.Println("email: ", email)
	fmt.Println("telephone: ", telephone)
	fmt.Println("password: ", password)

	user := model.User{
		Email:		email,
		Telephone:	telephone,
		Password:	password,
		IsAdmin:	true,
	}

	var err error
	if len(user.Email) == 0 || len(user.Password) == 0 {
		err = errors.New("email or password cannot be null.")
	} else {
		user.Password = user.EncryptPassword(password, user.Salt())
	}

	if err = model.DB.Create(&user).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"succeed": true,
		})
		return
	} else {
		err = errors.New("email already exists.")
	}

	c.JSON(http.StatusOK, gin.H{
		"succeed":false,
		"message":err.Error(),
	})
}


// Signout 退出登录
func Signout(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}
}


func AdminUserIndexGet(c *gin.Context) {

	if user, exists  := c.Get("user"); exists {

		var comments []string
		comments = append(comments, "nihao","haodehen","xiaodai")

		var users []*model.User
		rows, err := model.DB.Raw("select * from user").Rows()
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				var newUser model.User
				err = model.DB.ScanRows(rows, &newUser)
				if err == nil {
					users = append(users, &newUser)
				} else {
					err = errors.New("扫描用户失败")
				}
			}

			c.HTML(http.StatusOK, "admin/user.html", gin.H{
				"users": users,
				"user": user,
				"comments": comments,
			})
			return
		}
	}
}

func AdminUserLock(c *gin.Context) {
	id := c.Param("id")
	userID, _ := strconv.ParseUint(id, 10, 64)

	// Find User
	var err error
	var newUser model.User
	err = model.DB.Where("id = ?", userID).First(&newUser).Error
	if err == nil {
		newUser.LockState = !newUser.LockState
		err = newUser.Lock()
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"succeed": true,
			})
			return
		} else {
			err = errors.New("锁定用户失败")
		}
	} else {
		err = errors.New("查找用户失败")
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": false,
		"message": err.Error(),
	})
	return
}