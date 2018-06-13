package article

import (
	"errors"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/golangboy/model"
	"time"
	"strings"
)


func ArticleIndexGet(c *gin.Context) {

	recentArts, _ := model.GetRecentArticleQuerys()
	articleCates, _ := model.GetArticleCategoryQuerys()
	articleArchs, _ := model.GetArticleArchiveQuerys()

	c.HTML(http.StatusOK, "client/index.html", gin.H{
		"recentArts"	: recentArts,
		"articleCates"	: articleCates,
		"articleArchs"	: articleArchs,
	})
	return
}

func ArticleListGet(c *gin.Context) {

	pageStr 		:= c.Query("page")
	perPageStr 		:= c.Query("per_page")

	pageInt, _ 		:= strconv.Atoi(pageStr)
	perPageInt, _	:= strconv.Atoi(perPageStr)

	offset 			:= pageInt * perPageInt

	articleTotal 	:= model.GetArtileCount()
	data 			:= model.GetArticleJsonData(perPageInt, offset)

	c.JSON(http.StatusOK, gin.H{
		"status" : 0,
		"data"	 : data,
		"total"	 : articleTotal,
	})
	return
}

func ArticleCategoryGet(c *gin.Context) {
	id				:= c.Param("id")
	recentArts, _ 	:= model.GetRecentArticleQuerys()
	articleCates, _ := model.GetArticleCategoryQuerys()
	articleArchs, _ := model.GetArticleArchiveQuerys()

	cateId, _ := strconv.Atoi(id)
	newCateId := uint64(cateId)

	cateArticles, err := model.GetArticleQuerysByCateId(newCateId)
	if err == nil {

		c.HTML(http.StatusOK, "client/category.html", gin.H{
			"cateArticles"	: cateArticles,
			"recentArts"	: recentArts,
			"articleCates"	: articleCates,
			"articleArchs"	: articleArchs,
		})
		return
	}
}


func ArticleArchiveGet(c *gin.Context) {
	yearMonth := c.Param("yearMonth")

	year, err 	:= time.Parse("2006-01", yearMonth)
	year_month 	:= year.Format("2006-01")

	recentArts, _ 	:= model.GetRecentArticleQuerys()
	articleCates, _ := model.GetArticleCategoryQuerys()
	articleArchs, _ := model.GetArticleArchiveQuerys()

	archArticles, err := model.GetArticleArchiveQuerysByTime(year_month)
	if err == nil {
		c.HTML(http.StatusOK, "client/archive.html", gin.H{
			"archArticles"	: archArticles,
			"recentArts"	: recentArts,
			"articleCates"	: articleCates,
			"articleArchs"	: articleArchs,
		})
		return
	}
}


func ArticleDetailGet(c*gin.Context) {
	id := c.Param("id")

	article, _ 	:= model.GetArticleById(id)
	recentArts, _ 	:= model.GetRecentArticleQuerys()
	articleCates, _ := model.GetArticleCategoryQuerys()
	articleArchs, _ := model.GetArticleArchiveQuerys()

	article.ViewCount++
	err := article.UpdateView()
	if err == nil {
		c.HTML(http.StatusOK, "client/detail.html", gin.H{
			"recentArts"	: recentArts,
			"articleCates"	: articleCates,
			"articleArchs"	: articleArchs,
			"artID"			: id,
			"artTitle"		: article.Title,
			"article"		: article,
		})
		return
	}
}


func AjaxArticleDetailGet(c *gin.Context) {
	id := c.Param("id")

	article, err := model.GetArticleById(id)
	if err == nil {
		var markdown []string

		markdown = append(markdown, article.Content)

		c.String(http.StatusOK, article.Content)
		return
	}
}

func ArticleSearchPost(c *gin.Context) {
	keyword			:= c.PostForm("keyword")
	recentArts, _ 	:= model.GetRecentArticleQuerys()
	articleCates, _ := model.GetArticleCategoryQuerys()
	articleArchs, _ := model.GetArticleArchiveQuerys()

	searchArts, err := model.GetArticlesByKeyword(keyword)
	if err != nil {
		searchArts = nil
	}
	c.HTML(http.StatusOK, "client/search.html", gin.H{
		"searchArts"	: searchArts,
		"recentArts"	: recentArts,
		"articleCates"	: articleCates,
		"articleArchs"	: articleArchs,
	})
	return
}


func AdminArticleGet(c *gin.Context) {
	if _, exists := c.Get("user"); exists {

		id := c.Param("id")

		article, _ 	:= model.GetArticleById(id)
		category, _ := model.GetCategoryById(article.CategoryID)

		user, err 	:= model.GetUserById(article.UserID)
		if err == nil {
			var comments []string

			c.HTML(http.StatusOK, "post/display.html", gin.H{
				"article"	: article,
				"category"	: category,
				"user"		: user,
				"comments"	: comments,
			})
			return
		}

		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"message": "Sorry,I lost myself!",
		})
		return
	}
}

func AdminArticleIndex(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		var comments []string

		// Find Article By UserId
		articles, err := model.GetArticleQuerysByUserId(userInter.ID)
		if err == nil {
			c.HTML(http.StatusOK, "admin/post.html", gin.H{
				"articles"	: articles,
				"Active"	: "posts",
				"user"		: user,
				"comments"	: comments,
			})
			return
		} else {
			err = errors.New("获取文章列表失败")
		}
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	}
	return
}

func AdminCreateArticleGet(c *gin.Context) {
	if user, exists := c.Get("user"); exists {

		var comments []string

		c.HTML(http.StatusOK, "post/new.html", gin.H{
			"user"		: user,
			"comments"	: comments,
		})
	}
	return
}

func AdminCreateArticlePost(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		tags 		:= c.PostForm("tags")
		title 		:= c.PostForm("title")
		body 		:= c.PostForm("body")
		isPublished := c.PostForm("isPublished")
		published 	:= "on" == isPublished

		tagList := strings.Split(tags, ",")

		var desc string

		if len(body) > 100 {
			desc = body[:100]
		} else {
			desc = body
		}

		cateID, _ := strconv.ParseUint(tagList[0], 10, 64)

		article := model.Article{
			Title		: title,
			Desc		: desc,
			Content		: body,
			UserID		: userInter.ID,
			ViewCount	: 0,
			CommentCount: 0,
			CategoryID	: cateID,
			IsPublished	: published,
		}

		err := article.Insert()
		if err == nil {
			c.Redirect(http.StatusMovedPermanently, "/admin/article")
		} else {
			c.HTML(http.StatusOK, "post/new.html", gin.H{
				"article": article,
				"message": err.Error(),
			})
		}
	}
	return
}


func AdminEditGET(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		id := c.Param("id")

		if article, err := model.GetArticleById(id); err == nil {
			var comments []string

			//Find Article Category
			if category, err := model.GetCategoryById(article.CategoryID); err == nil {

				article.Category = *category

				c.HTML(http.StatusOK, "post/modify.html", gin.H{
					"user"		: user,
					"article"	: article,
					"comments"	: comments,
				})
				return
			}
		}
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"message": "Sorry,I lost myself!",
		})
		return
	}
}

func AdminUpdateArticle(c *gin.Context) {
	if _, exists := c.Get("user");exists {

		id 			:= c.Param("id")
		cateId 		:= c.PostForm("tags")
		title 		:= c.PostForm("title")
		content 	:= c.PostForm("body")
		isPublished := c.PostForm("isPublished")

		published 	:= "on" == isPublished

		var desc string

		if len(content) > 100 {
			desc = content[:100]
		} else {
			desc = content
		}

		cateID, _ := strconv.ParseUint(cateId, 10, 64)

		pid, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			article := &model.Article{
				Title		: title,
				Desc		:desc,
				Content		:content,
				CategoryID	: cateID,
				IsPublished	: published,
			}
			article.ID = pid

			err = article.Update()
			if err == nil {
				c.Redirect(http.StatusMovedPermanently, "/admin/article")
				return
			} else {
				c.HTML(http.StatusOK, "post/modify.html", gin.H{
					"article": article,
					"message": err.Error(),
				})
			}
		} else {
			c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
				"message": "Sorry,I lost myself!",
			})
		}
	}
	return
}

func AdminDeleteArticle(c *gin.Context) {

	var err error

	id := c.Param("id")

	if article, err := model.GetArticleById(id); err == nil {

		err := article.Delete()
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"succeed": true,
			})
			return
		} else {
			err = errors.New("删除文章失败")
		}
	} else {
		err = errors.New("获取文章失败")
	}

	c.JSON(http.StatusOK, gin.H{
		"succeed": false,
		"message": err.Error(),
	})
	return
}

func AdminPublishArticle(c *gin.Context) {
	if _, exists := c.Get("user"); exists {

		id := c.Param("id")
		article, err := model.GetArticleById(id)
		if err == nil {
			article.IsPublished = !article.IsPublished
			err = article.Update()
		}

		c.JSON(http.StatusOK, gin.H{
			"succeed": err == nil,
		})
	}
}


