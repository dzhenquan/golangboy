package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"gin-blog/model"
	"strconv"
	"github.com/pkg/errors"
)


func UserArticleIndexGet(c *gin.Context) {

	rows, err := model.DB.Raw("select * from article ORDER BY updated_at desc LIMIT 10").Rows()
	if err == nil {

		var articles []*model.Article
		for rows.Next() {
			var article model.Article
			model.DB.ScanRows(rows, &article)

			if article.IsPublished {

				user, _ := model.GetUserById(article.UserID)
				cate, _ := model.GetCategoryById(article.CategoryID)

				article.User = *user
				article.Category = *cate
				articles = append(articles, &article)
			}
		}
		c.HTML(http.StatusOK, "auth/index.html", gin.H{
			"Articles": articles,
		})
		return
	}
	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"succeed": false,
		"message": err.Error(),
	})
	return
}


func AdminArticleGet(c *gin.Context) {
	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID:", userInter.ID)
		fmt.Println("userEmail:", userInter.Email)

		id := c.Param("id")
		article, _ := model.GetArticleById(id)

		category, _ := model.GetCategoryById(article.CategoryID)

		user, err := model.GetUserById(article.UserID)
		if err == nil && article.IsPublished {
			//article.IsPublished = !article.IsPublished
			//err = article.Update()
			var comments []string

			c.HTML(http.StatusOK, "post/display.html", gin.H{
				"article": article,
				"category": category,
				"user": user,
				"comments": comments,
			})
			return
		}
		fmt.Println(err)
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"message": "Sorry,I lost myself!",
		})
		return
	}
}

func AdminArticleIndex(c *gin.Context) {

	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)
		fmt.Println("userID:", userInter.ID)
		fmt.Println("userEmail:", userInter.Email)

		var comments []string

		// Find Article By UserId
		var articles []*model.Article

		//model.DB.Where("id = ?", userInter.ID).Find(&article)

		rows, err := model.DB.Raw("select * from article where user_id = ?", userInter.ID).Rows()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"succeed": false,
			})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var article model.Article
			model.DB.ScanRows(rows, &article)
			articles = append(articles, &article)
		}

		c.HTML(http.StatusOK, "admin/post.html", gin.H{
			"articles":    articles,
			"Active":   "posts",
			"user":     user,
			"comments": comments,
		})
	}
	/*
	fmt.Println("username",c)
	var comments []string
	comments = append(comments, "nihao","haodehen")

	var user model.User

	err := model.DB.First(&user, "email = ?", "123@qq.com").Error
	if err == nil {
		c.HTML(http.StatusOK, "admin/post.html", gin.H{
			"pageCount":    2,
			"postCount":    3,
			"tagCount":     1,
			"commentCount": 5,
			"user":         user,
			"comments":     comments,
		})
	}*/
}

func AdminNewPostGet(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)
		fmt.Println("userID:", userInter.ID)
		fmt.Println("userEmail:", userInter.Email)

		var comments []string
		c.HTML(http.StatusOK, "post/new.html", gin.H{
			"user": user,
			"comments": comments,
		})
	}
}

func AdminCreatePost(c *gin.Context) {
	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)
		fmt.Println("userID:", userInter.ID)
		fmt.Println("userEmail:", userInter.Email)


		tags := c.PostForm("tags")
		title := c.PostForm("title")
		body := c.PostForm("body")
		isPublished := c.PostForm("isPublished")
		published := "on" == isPublished

		fmt.Println("title: ", title)
		fmt.Println("tags: ", tags)
		fmt.Println("pubisted: ", published)
		fmt.Println("body: ", body)

		var desc string
		if len(body) > 100 {
			desc = body[:100]
		} else {
			desc = body
		}

		cateID, _ := strconv.ParseUint(tags, 10, 64)

		article := model.Article{
			Title: title,
			Desc: desc,
			Content: body,
			UserID: userInter.ID,
			CategoryID: cateID,
			IsPublished: published,
		}

		err := article.Insert()
		if err == nil {
			c.Redirect(http.StatusMovedPermanently, "/admin/post")
		} else {
			c.HTML(http.StatusOK, "post/new.html", gin.H{
				"article":    article,
				"message": err.Error(),
			})
		}
	}
}


func AdminEditGET(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)
		fmt.Println("userID:", userInter.ID)
		fmt.Println("userEmail:", userInter.Email)

		id := c.Param("id")

		fmt.Println("id:", id)

		var article model.Article
		err := model.DB.Where("id = ?", id).First(&article).Error
		if err == nil {

			var comments []string
			fmt.Println("cateID:", article.CategoryID)
			//Find Article Category
			err := model.DB.Where("id = ?", article.CategoryID).First(&article.Category).Error
			if err == nil {
				fmt.Println("cateName:", article.Category.Name)
				c.HTML(http.StatusOK, "post/modify.html", gin.H{
					"user": user,
					"article": article,
					"comments": comments,
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
	if user, exists := c.Get("user");exists {
		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

		id := c.Param("id")
		cateId := c.PostForm("tags")
		title := c.PostForm("title")
		content := c.PostForm("body")
		isPublished := c.PostForm("isPublished")

		published := "on" == isPublished

		cateID, _ := strconv.ParseUint(cateId, 10, 64)
		pid, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			article := &model.Article{
				Title: title,
				Content:content,
				CategoryID: cateID,
				IsPublished: published,
			}
			article.ID = pid

			err = article.Update()
			if err == nil {
				c.Redirect(http.StatusMovedPermanently, "/admin/post")
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

func AdminArticleDelete(c *gin.Context) {
	id := c.Param("id")
	var err error
	articleID, err := strconv.ParseUint(id, 10, 64)
	if err == nil {
		var article model.Article
		err := model.DB.Where("id = ?", articleID).First(&article).Error
		if err == nil {
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
	} else {
		err = errors.New("获取文章ID失败")
	}

	c.JSON(http.StatusOK, gin.H{
		"succeed": false,
		"message": err.Error(),
	})
	return
}

func AdminArticlePublish(c *gin.Context) {
	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		fmt.Println("userID: ", userInter.ID)
		fmt.Println("userEmail: ", userInter.Email)

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