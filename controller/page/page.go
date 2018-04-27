package page

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gin-blog/model"
	"fmt"
	"strconv"
)


func AdminPageGet(c *gin.Context) {
	id := c.Param("id")

	page, err := model.GetPageById(id)
	if err == nil && page.IsPublished {
		page.View++
		c.HTML(http.StatusOK, "page/display.html", gin.H{
			"page": page,
		})
		return
	} else {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"message": "Sorry,I lost myself!",
		})
	}
}

func AdminCreatePageGet(c *gin.Context) {
	c.HTML(http.StatusOK, "page/new.html", nil)
}

func AdminPageIndex(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)
		var comments []string
		comments = append(comments, "nihao","haodehen","xiaodai")

		rows, err := model.DB.Raw("select * from page where user_id = ?", userInter.ID).Rows()
		defer rows.Close()

		if err == nil {
			var page model.Page
			var pages []*model.Page

			for rows.Next() {
				model.DB.ScanRows(rows, &page)
				pages = append(pages, &page)
			}
			c.HTML(http.StatusOK, "admin/page.html", gin.H{
				"pages": pages,
				"user": user,
				"comments": comments,
			})
			return
		}
	}
}

func AdminCreatePage(c *gin.Context) {

	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		title := c.PostForm("title")
		content := c.PostForm("body")
		isPublished := c.PostForm("isPublished")

		published := "on" == isPublished

		page := &model.Page{
			Title: title,
			Content: content,
			UserID: userInter.ID,
			IsPublished: published,
		}
		err := page.Insert()
		if err == nil {
			c.Redirect(http.StatusMovedPermanently, "/admin/page")
			return
		} else {
			c.HTML(http.StatusOK, "page/new.html", gin.H{
				"message": err.Error(),
				"page": page,
			})
		}
	}


}

func AdminEditPage(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		id := c.Param("id")
		page, err := model.GetPageById(id)
		if err == nil {
			c.HTML(http.StatusOK, "page/modify.html", gin.H{
				"page": page,
			})
			return
		} else {
			c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
				"message": "Sorry,I lost myself!",
			})
		}
	}
}

func AdminUpdatePage(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		id := c.Param("id")
		title := c.PostForm("title")
		content := c.PostForm("body")
		isPublished := c.PostForm("isPublished")
		published := "on" == isPublished

		fmt.Println("content:", content)

		pid, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			page := &model.Page{
				Title: title,
				Content: content,
				IsPublished: published,
			}
			page.ID = pid

			err = page.Update()
			if err == nil {
				c.Redirect(http.StatusMovedPermanently, "/admin/page")
				return
			}
		}
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func AdminPublishPage(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		id := c.Param("id")
		page, err := model.GetPageById(id)
		if err == nil {
			page.IsPublished = !page.IsPublished
			err = page.Update()
		}

		c.JSON(http.StatusOK, gin.H{
			"succeed": err == nil,
		})
	}
}

func AdminDeletePage(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		id := c.Param("id")
		pageID, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			page := &model.Page{}
			page.ID = pageID

			err = page.Delete()
			if err == nil {
				c.JSON(http.StatusOK, gin.H{
					"succeed": true,
				})
			}
		}
	}
}