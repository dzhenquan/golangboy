package index

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func SigninGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"name":"xiaodai",
		"age":19,
	})
}

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", gin.H{
		"name":"xiaodai",
		"age":19,
	})
}

func SignUpPost(c *gin.Context) {

	email := c.PostForm("email")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	fmt.Println("email: ", email)
	fmt.Println("telephone: ", telephone)
	fmt.Println("password: ", password)

	c.HTML(http.StatusOK, "auth/signup.html", gin.H{
		"name":"xiaodai",
		"age":19,
	})
}