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


//	Admin SignIn Index Get 管理员登录
func AdminSignIndexGet(c *gin.Context) {

	var err error
	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		articleCount := model.GetArticleCountByUserId(userInter.ID)
		pageCount := model.GetPageCountByUserId(userInter.ID)
		cateCount := model.GetCateCountByUserId(userInter.ID)

		var comments []string
		//comments = append(comments, "nihao","haodehen","xiaodai")

		c.HTML(http.StatusOK, "admin/index.html", gin.H{
			"pageCount":    	pageCount,
			"articleCount":    	articleCount,
			"cateCount":     	cateCount,
			"commentCount": 	0,
			"user":         	user,
			"comments":     	comments,
		})
		return
	} else {
		err = errors.New("当前未登录,请先登录")
	}

	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"message": err.Error(),
	})
	return
}


// User SignIn Get  用户登录
func UserSignInGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signin.html", nil)
}


//User SignIn Post  用户登录
func UserSignInPost(c *gin.Context) {
	useremail := c.PostForm("useremail")
	password := c.PostForm("password")

	useremail = strings.TrimSpace(useremail)

	var err error
	if len(useremail) > 0 && len(password) > 0 {
		var user model.User
		err = model.DB.First(&user, "email = ?", useremail).Error
		if err == nil {
			if user.Password == user.EncryptPassword(password, user.Salt()) {
				if !user.LockState {
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
						"id": user.ID,
					})

					var tokenString string
					tokenString, err = token.SignedString([]byte(config.ServerConfig.TokenSecret))
					if err == nil {
						cookie := &http.Cookie{
							Name:     "token",
							Value:    tokenString,
							MaxAge:   config.ServerConfig.TokenMaxAge,
							Path:     "/",
							Domain:   "",
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
						err = errors.New("服务器内部错误")
					}
				} else {
					err = errors.New("该账户被锁定,请联系管理员解锁")
				}
			} else {
				err = errors.New("登录密码不正确")
			}
		} else {
			err = errors.New("登录邮箱不存在,请注册后登录")
		}
	} else {
		err = errors.New("登录邮箱或登录密码不能为空")
	}
	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"message": err.Error(),
	})
}


// User SignUp Get  用户注册
func UserSignUpGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", nil)
}

// User SignUp Post 用户注册
func UserSignUpPost(c *gin.Context) {

	var err error

	useremail := c.PostForm("useremail")
	password :=	c.PostForm("password")

	if len(useremail) > 0 && len(password) > 0 {
		if len(password) >= 6 {
			useremail = utils.AvoidXSS(useremail)
			useremail = strings.TrimSpace(useremail)

			user := &model.User{
				Email:		useremail,
				Password:	password,
				IsAdmin:	true,
			}
			user.Password = user.EncryptPassword(password, user.Salt())
			err = user.Insert()
			if err == nil {
				c.JSON(http.StatusOK, gin.H{
					"succeed": true,
				})
				return
			} else {
				err = errors.New("该账户已存在,请登录")
			}
		} else {
			err = errors.New("登录密码至少6位,请重新输入")
		}
	} else {
		err = errors.New("登录邮箱或登录密码不为空")
	}

	c.JSON(http.StatusOK, gin.H{
		"succeed":false,
		"message":err.Error(),
	})
	return
}


// Signout 退出登录
func UserSignout(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		c.Redirect(http.StatusSeeOther, "/signin")
		return
	}
}


// Admin User Index 获取用户列表
func AdminUserIndexGet(c *gin.Context) {

	var err error
	if user, exists  := c.Get("user"); exists {

		var comments []string

		rows, rowErr := model.DB.Raw("select * from user").Rows()
		defer rows.Close()
		if rowErr == nil {
			var users []*model.User

			for rows.Next() {
				var newUser model.User
				err = model.DB.ScanRows(rows, &newUser)
				if err == nil {
					users = append(users, &newUser)
				} else {
					err = errors.New("服务器内部错误")
				}
			}

			c.HTML(http.StatusOK, "admin/user.html", gin.H{
				"users": users,
				"user": user,
				"comments": comments,
			})
			return
		} else {
			err = errors.New("获取用户列表失败")
		}
	} else {
		err = errors.New("当前未登录,请先登录")
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": err.Error(),
		})
	}

	return
}

// Admin User Lock 锁定用户
func AdminUserLock(c *gin.Context) {

	var err error
	if _, exists := c.Get("user"); exists {
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
	} else {
		err = errors.New("当前未登录,请先登录")
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": err.Error(),
		})
	}
	return
}