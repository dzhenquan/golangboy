package router

import (
	"path/filepath"
	"html/template"
	"github.com/gin-gonic/gin"
	"github.com/dzhenquan/golangboy/utils"
	"github.com/dzhenquan/golangboy/middleware"
	"github.com/dzhenquan/golangboy/controller/auth"
	"github.com/dzhenquan/golangboy/controller/page"
	"github.com/dzhenquan/golangboy/controller/link"
	"github.com/dzhenquan/golangboy/controller/article"
	"github.com/dzhenquan/golangboy/controller/category"
)


//Route 路由
func Route(router *gin.Engine) {
	//apiPrefix := config.ServerConfig.APIPrefix

	funcMap := template.FuncMap{
		"add"			: utils.Add,
		"isOdd"			: utils.IsOdd,
		"isEven"		: utils.IsEven,
		"truncate"		: utils.Truncate,
		"substring"		: utils.Substring,
		"dateFormat"	: utils.DateFormat,
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
		api.GET("/logout",  auth.SignOutGet)

		api.GET("/", article.ArticleIndexGet)
		api.GET("/login", auth.SignInGet)

		api.GET("/about/me", page.PageAboutMeGet)
		api.GET("/about/detail/:id", page.AjaxPageDetailGet)

		api.GET("/article/list", article.ArticleListGet)
		api.POST("/article/search", article.ArticleSearchPost)
		api.GET("/article/category/:id", article.ArticleCategoryGet)
		api.GET("/article/archive/:yearMonth", article.ArticleArchiveGet)

		api.GET("/article/detail/:id", article.ArticleDetailGet)
		api.GET("/ajax/article/detail/:id", article.AjaxArticleDetailGet)
	}

	adminAPI := router.Group("/admin", middleware.RefreshTokenCookie)
	{
		adminAPI.GET("/index", middleware.SigninRequired,
			auth.AdminSignIndexGet)

		// User
		adminAPI.GET("/user", middleware.SigninRequired,
			auth.AdminUserIndexGet)
		adminAPI.POST("/user/:id/lock", middleware.SigninRequired,
			auth.AdminUserLockPost)

		// Profile
		adminAPI.GET("/profile", middleware.SigninRequired,
			auth.AdminProfileGet)
		adminAPI.POST("/profile/upload/image", middleware.SigninRequired,
			auth.AdminUploadImage)
		adminAPI.POST("/profile/update/userpwd", middleware.SigninRequired,
			auth.AdminUpdateUserPwd)
		adminAPI.POST("/profile/update/userinfo", middleware.SigninRequired,
			auth.AdminUpdateUserInfo)

		// Category
		adminAPI.POST("/new_category", middleware.SigninRequired,
			category.AdminCreateCategory)
		/*adminAPI.POST("/new_tag", middleware.SigninRequired,
			category.AdminGetCategoryQuerys)*/
		adminAPI.GET("/category", middleware.SigninRequired,
				category.AdminCategoryIndex)
		adminAPI.GET("/category/:id", middleware.SigninRequired,
			category.AdminArticleByCateId)

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
			page.AdminCreatePagePost)
		adminAPI.GET("/page/:id", middleware.SigninRequired,
			page.AdminPageGet)
		adminAPI.GET("/page/:id/edit", middleware.SigninRequired,
			page.AdminEditPage)
		adminAPI.POST("/page/:id/edit", middleware.SigninRequired,
			page.AdminUpdatePage)
		adminAPI.POST("/page/:id/publish", middleware.SigninRequired,
			page.AdminPublishPage)
		adminAPI.POST("/page/:id/delete", middleware.SigninRequired,
			page.AdminDeletePage)

		// Article
		adminAPI.GET("/article", middleware.SigninRequired,
			article.AdminArticleIndex)
		adminAPI.GET("/new_article", middleware.SigninRequired,
			article.AdminCreateArticleGet)
		adminAPI.POST("/new_article", middleware.SigninRequired,
			article.AdminCreateArticlePost)

		adminAPI.GET("/article/:id", middleware.SigninRequired,
			article.AdminArticleGet)

		adminAPI.GET("/article/:id/edit", middleware.SigninRequired,
			article.AdminEditGET)
		adminAPI.POST("/article/:id/edit", middleware.SigninRequired,
			article.AdminUpdateArticle)
		adminAPI.POST("/article/:id/publish", middleware.SigninRequired,
			article.AdminPublishArticle)
		adminAPI.POST("/article/:id/delete", middleware.SigninRequired,
			article.AdminDeleteArticle)
	}
}

func getCurrentDirectory() string {
	return ""
}


















