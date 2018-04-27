package router

import (
	"github.com/gin-gonic/gin"
	//"gin-blog/config"
	"path/filepath"
	"gin-blog/middleware"
	"gin-blog/controller/auth"
	"html/template"
	"gin-blog/utils"
	"gin-blog/controller/article"
	"gin-blog/controller/category"
	"gin-blog/controller/page"
	"gin-blog/controller/link"
)


//Route 路由
func Route(router *gin.Engine) {
	//apiPrefix := config.ServerConfig.APIPrefix

	funcMap := template.FuncMap{
		"dateFormat": utils.DateFormat,
		"substring":  utils.Substring,
		"isOdd":      utils.IsOdd,
		"isEven":     utils.IsEven,
		"truncate":   utils.Truncate,
		"add":        utils.Add,
		//"listtag":    helpers.ListTag,
	}

	router.SetFuncMap(funcMap)

	router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
	router.StaticFile("/favicon.ico", filepath.Join(getCurrentDirectory(), "./static/favicon.ico"))

	router.LoadHTMLGlob("views/**/*")

	api := router.Group("", middleware.RefreshTokenCookie)
	{
		api.GET("/signin", auth.SignInGet)
		api.POST("/signin", auth.SignInPost)

		api.GET("/signup", auth.SignUpGet)
		api.POST("/signup", auth.SignUpPost)
		api.GET("/logout", middleware.SigninRequired,
			auth.Signout)

		api.GET("/page/:id", page.AdminPageGet)
		api.GET("/post/:id", article.AdminArticleGet)
	}

	adminAPI := router.Group("/admin", middleware.RefreshTokenCookie)
	{
		adminAPI.POST("/logout", middleware.SigninRequired,
			auth.Signout)

		adminAPI.GET("/index", middleware.SigninRequired,
			auth.AdminSignInGet)

		// User
		adminAPI.GET("/user", middleware.SigninRequired,
			auth.AdminUserIndexGet)
		adminAPI.POST("/user/:id/lock", middleware.SigninRequired,
			auth.AdminUserLock)

		// Category
		adminAPI.POST("/new_tag", middleware.SigninRequired,
			category.AdminCreateCategory)
		/*adminAPI.POST("/new_tag", middleware.SigninRequired,
			category.AdminGetCategoryQuerys)*/

		// Link
		adminAPI.GET("/link", middleware.SigninRequired,
			link.AdminLinkIndex)
		adminAPI.POST("/new_link", middleware.SigninRequired,
			link.AdminLinkCreate)
		adminAPI.POST("/link/:id/edit", middleware.SigninRequired,
			link.AdminLinkUpdate)
		adminAPI.POST("/link/:id/delete", middleware.SigninRequired,
			link.AdminLinkDelete)

		// Page
		adminAPI.GET("/page", middleware.SigninRequired,
			page.AdminPageIndex)
		adminAPI.GET("/new_page", middleware.SigninRequired,
			page.AdminCreatePageGet)
		adminAPI.POST("/new_page", middleware.SigninRequired,
			page.AdminCreatePage)
		adminAPI.GET("/page/:id/edit", middleware.SigninRequired,
			page.AdminEditPage)
		adminAPI.POST("/page/:id/edit", middleware.SigninRequired,
			page.AdminUpdatePage)
		adminAPI.POST("/page/:id/publish", middleware.SigninRequired,
			page.AdminPublishPage)
		adminAPI.POST("/page/:id/delete", middleware.SigninRequired,
			page.AdminDeletePage)

		// Article
		adminAPI.GET("/post", middleware.SigninRequired,
			article.AdminArticleIndex)
		adminAPI.GET("/new_post", article.AdminNewPostGet)
		adminAPI.POST("/new_post", middleware.SigninRequired,
			article.AdminCreatePost)
		adminAPI.GET("/post/:id/edit", middleware.SigninRequired,
			article.AdminEditGET)
		adminAPI.POST("/post/:id/edit", middleware.SigninRequired,
			article.AdminUpdateArticle)
		adminAPI.POST("/post/:id/publish", middleware.SigninRequired,
			article.AdminArticlePublish)
		adminAPI.POST("/post/:id/delete", middleware.SigninRequired,
			article.AdminArticleDelete)
	}
}

func getCurrentDirectory() string {
	return ""
}


















