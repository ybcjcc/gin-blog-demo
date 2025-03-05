package routes

import (
    "your-project/controllers"
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
    
    // 创建博客控制器实例
    blogController := &controllers.BlogController{}
    
    // 博客相关路由组
    blog := r.Group("/api/blog")
    {
        // 获取博客详情
        blog.GET("/detail/:id", blogController.GetBlogDetail)
        // 这里可以添加其他博客相关的路由
    }
    
    return r
} 