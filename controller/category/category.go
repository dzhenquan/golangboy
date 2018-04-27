package category

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"gin-blog/model"
	"net/http"
	"github.com/pkg/errors"
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

	if user, exists := c.Get("user"); exists {
		userInter := user.(model.User)
		fmt.Println("userID:", userInter.ID)
		fmt.Println("userEmail:", userInter.Email)

		cateName := c.PostForm("value")
		fmt.Println("cateName: ", cateName)

		// Find Category
		var findCate model.Category
		if err = model.DB.Where("name = ?", cateName).First(&findCate).Error; err == nil {
			fmt.Println("id: ", findCate.ID)
			fmt.Println("name: ", findCate.Name)

			c.JSON(http.StatusOK, gin.H{
				"data": findCate,
			})
			return
		}

		fmt.Println("没有发现 ", cateName," 开始插入!")

		var newCate model.Category
		newCate.Name = cateName

		if err = model.DB.Create(&newCate).Error; err == nil {
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
}