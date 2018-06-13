package category

import (
	"fmt"
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/golangboy/model"
	"strconv"
)



func AdminGetCategoryQuerys(c *gin.Context) {

	var err error

	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)

		fmt.Println("userID:", userInter.ID)
		fmt.Println("userEmail:", userInter.Email)


		//rows, err1 := model.DB.Where("user_id = ?", userInter.ID).Find(&model.Article{}).Rows();
		rows, err1 := model.DB.Raw("select DISTINCT category.name from category inner join article on category.id=article.category_id where article.user_id = ?", userInter.ID).Rows();
		if err1 == nil {
			for rows.Next() {
				//var article model.Article
				//var cate string
				var cate model.Category
				model.DB.ScanRows(rows, &cate)
				fmt.Println("article_cate_title: ", cate.Name)
			}
		}

		defer rows.Close()

		c.JSON(http.StatusOK, gin.H{
			"succeed":true,
		})

		return
	} else {
		err = errors.New("用户不存在")
	}

	c.JSON(http.StatusOK, gin.H{
		"message":err.Error(),
	})
}



func AdminCreateCategory(c *gin.Context) {

	var err error

	if _, exists := c.Get("user"); exists {

		cateName := c.PostForm("value")

		// Find Category
		if findCate, err := model.GetCategoryByName(cateName); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"data": findCate,
			})
			return
		}

		var newCate model.Category
		newCate.Name = cateName

		// Insert Category
		if err = newCate.Insert(); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"data": newCate,
			})
			return
		} else {
			err = errors.New("插入失败")
		}
	} else {
		err = errors.New("用户不存在")
	}

	c.JSON(http.StatusOK, gin.H{
		"message":err.Error(),
	})
	return
}


func AdminCategoryIndex(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		userInter := user.(model.User)

		var comments []string

		categorys, err := model.GetCategoryQuerysByUserId(userInter.ID)
		if err != nil {
			fmt.Println("err: ", err)
			categorys = nil
		}
		c.HTML(http.StatusOK, "admin/category.html", gin.H{
			"categorys"	: categorys,
			"user"		: user,
			"comments"	: comments,
		})
		return
	}
}


func AdminArticleByCateId(c *gin.Context) {

	if user, exists := c.Get("user"); exists {

		id := c.Param("id")
		userInter := user.(model.User)

		var comments []string

		cateId, _ := strconv.Atoi(id)

		articles, err := model.AdminGetArticleByCategory(userInter.ID, uint64(cateId))
		if err == nil {
			c.HTML(http.StatusOK, "admin/cate_article.html", gin.H{
				"articles"	: articles,
				"user"		: user,
				"comments"	: comments,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	}
}
