package link

import (
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/golangboy/model"
	"errors"
)

func AdminLinkIndex(c *gin.Context) {

	var err error

	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		var comments []string

		links, err := model.GetLinkQuerysByUserID(userInter.ID)
		if err == nil {
			c.HTML(http.StatusOK, "admin/link.html", gin.H{
				"links"		: links,
				"user"		: user,
				"comments"	: comments,
			})
			return
		} else {
			err = errors.New("查看链接失败")
		}
	} else {
		err = errors.New("用户不存在")
	}

	c.JSON(http.StatusOK, gin.H{
		"message":err.Error(),
	})
}

func AdminLinkCreate(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		name := c.PostForm("name")
		url := c.PostForm("url")
		sort := c.PostForm("sort")

		if len(name) == 0 || len(url) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": "error parameter!",
			})
		} else {
			_sort, err := strconv.ParseUint(sort, 10, 64)
			if err == nil {
				link := &model.Link{
					Name	: name,
					Url		: url,
					UserID	: userInter.ID,
					Sort	: int(_sort),
				}
				err = link.Insert()
			}
			c.JSON(http.StatusOK, gin.H{
				"succeed": err == nil,
			})
		}
	}
}

func AdminLinkUpdate(c *gin.Context) {

	if _, exists := c.Get("user"); exists {

		id := c.Param("id")
		name := c.PostForm("name")
		url := c.PostForm("url")
		sort := c.PostForm("sort")

		if len(name) == 0 || len(url) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": "error parameter!",
			})
		} else {
			var err error
			_id, err := strconv.ParseUint(id, 10, 64)
			_sort, err := strconv.ParseUint(sort, 10, 64)
			if err == nil {
				link := &model.Link{
					Name: name,
					Url: url,
					Sort: int(_sort),
				}
				link.ID = _id
				err = link.Update()
			}

			c.JSON(http.StatusOK, gin.H{
				"succeed": err == nil,
			})
		}
	}
}

func AdminLinkDelete(c *gin.Context) {

	if _, exists := c.Get("user"); exists {

		id := c.Param("id")

		var err error
		_id, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			link := &model.Link{}
			link.ID = _id

			err = link.Delete()
			if err == nil {

				c.JSON(http.StatusOK, gin.H{
					"succeed": true,
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": err.Error(),
			})
		}
	}
}