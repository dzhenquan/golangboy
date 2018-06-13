package page

import (
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/golangboy/model"
)



func PageAboutMeGet(c *gin.Context) {

	page, err 	:= model.GetPageById("1")
	recentArts, _ 	:= model.GetRecentArticleQuerys()
	articleCates, _ := model.GetArticleCategoryQuerys()
	articleArchs, _ := model.GetArticleArchiveQuerys()

	if err == nil {
		c.HTML(http.StatusOK, "client/about.html", gin.H{
			"recentArts"	: recentArts,
			"articleCates"	: articleCates,
			"articleArchs"	: articleArchs,
			"artID"			: page.ID,
			"artTitle"		: page.Title,
			"article"		: page,
		})
		return
	}
}

func AjaxPageDetailGet(c *gin.Context) {
	id := c.Param("id")

	page, err := model.GetPageById(id)
	if err == nil {
		var markdown []string

		markdown = append(markdown, page.Content)

		c.String(http.StatusOK, page.Content)
		return
	}
}

func AdminPageGet(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		id := c.Param("id")

		page, err := model.GetPageById(id)
		if err == nil {
			var comments []string
			c.HTML(http.StatusOK, "page/display.html", gin.H{
				"user": user,
				"page": page,
				"comments": comments,
			})
			return
		} else {
			c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
				"message": "Sorry,I lost myself!",
			})
		}
	}
	return
}

func AdminCreatePageGet(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		var comments []string
		c.HTML(http.StatusOK, "page/new.html", gin.H{
			"user": user,
			"comments": comments,
		})
	}
	return
}

func AdminPageIndex(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		var comments []string

		pages, err := model.GetPageQuerysByUserId(userInter.ID)
		if err == nil {
			c.HTML(http.StatusOK, "admin/page.html", gin.H{
				"pages": pages,
				"user": user,
				"comments": comments,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	}
	return
}

func AdminCreatePagePost(c *gin.Context) {

	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

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
	return
}

func AdminEditPage(c *gin.Context) {
	if user, exists := c.Get("user"); exists {

		id := c.Param("id")
		page, err := model.GetPageById(id)
		if err == nil {

			var comments []string
			c.HTML(http.StatusOK, "page/modify.html", gin.H{
				"user": user,
				"page": page,
				"comments": comments,
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
	if _, exists := c.Get("user"); exists {

		id := c.Param("id")
		title := c.PostForm("title")
		content := c.PostForm("body")
		isPublished := c.PostForm("isPublished")
		published := "on" == isPublished

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
	if _, exists := c.Get("user"); exists {

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
	return
}

func AdminDeletePage(c *gin.Context) {
	if _, exists := c.Get("user"); exists {

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
	return
}