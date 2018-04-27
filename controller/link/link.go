package link

import (
	"github.com/gin-gonic/gin"
	"gin-blog/model"
	"fmt"
	"net/http"
	"strconv"
)

func AdminLinkIndex(c *gin.Context) {

	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		var comments []string
		comments = append(comments, "nihao", "wohao", "henhao")

		rows, err := model.DB.Raw("select * from link where user_id = ?", userInter.ID).Rows()
		defer rows.Close()

		if err == nil {
			var link model.Link
			var links []*model.Link

			for rows.Next() {
				model.DB.ScanRows(rows, &link)
				links = append(links, &link)
			}
			c.HTML(http.StatusOK, "admin/link.html", gin.H{
				"links": links,
				"user": user,
				"comments": comments,
			})
			return
		}
	}
}

func AdminLinkCreate(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

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
					Name: name,
					Url: url,
					UserID: userInter.ID,
					Sort: int(_sort),
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

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

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

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		id := c.Param("id")

		var err error
		_id, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			link := &model.Link{}
			link.ID = _id

			err = link.Delete()
		}
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"succeed": true,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
				"message": err.Error(),
			})
		}
	}
}