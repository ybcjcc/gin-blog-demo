package main

import (
	"blog/config"
	"blog/controllers"
	"blog/database"
	"blog/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}
	defer database.DB.Close()

	// 创建Gin引擎
	r := gin.Default()

	// 加载静态文件
	r.Static("/static", "./static")
	// 加载HTML模板
	r.LoadHTMLGlob("views/*")

	// 前台路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.GET("/categories", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"activeTab": "categories"})
	})
	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", nil)
	})
	// 添加博客详情页面路由
	r.GET("/post/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post.html", nil)
	})

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		api.POST("/login", controllers.Login)
		api.POST("/register", controllers.Register)

		// 需要认证的路由组
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			// 文章相关路由
			auth.POST("/posts", controllers.CreatePost)
			auth.PUT("/posts/:id", controllers.UpdatePost)
			auth.DELETE("/posts/:id", controllers.DeletePost)

			// 分类管理路由
			auth.POST("/categories", controllers.CreateCategory)
			auth.PUT("/categories/:id", controllers.UpdateCategory)
			auth.DELETE("/categories/:id", controllers.DeleteCategory)

			// 标签管理路由
			auth.POST("/tags", controllers.CreateTag)
			auth.PUT("/tags/:id", controllers.UpdateTag)
			auth.DELETE("/tags/:id", controllers.DeleteTag)
			auth.GET("/posts/:id/tags", controllers.GetPostTags)
			auth.PUT("/posts/:id/tags", controllers.UpdatePostTags)
		}

		// 公开路由
		api.GET("/posts", controllers.ListPosts)
		api.GET("/posts/:id", controllers.GetPost)
		api.GET("/categories", controllers.ListCategories)
		api.GET("/categories/:id", controllers.GetCategory)
		api.GET("/tags", controllers.ListTags)
		api.GET("/tags/:id", controllers.GetTag)
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	if err := r.Run(addr); err != nil {
		panic(fmt.Sprintf("启动服务器失败: %v", err))
	}
}