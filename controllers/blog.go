package controllers

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type BlogController struct {
    // 可以在这里注入所需的服务
}

// GetBlogDetail 获取博客详情
func (bc *BlogController) GetBlogDetail(c *gin.Context) {
    // 从URL参数中获取博客ID
    blogID := c.Param("id")
    
    // TODO: 根据实际情况从数据库获取博客详情
    // 这里是示例数据
    blog := map[string]interface{}{
        "id":      blogID,
        "title":   "示例博客标题",
        "content": "这是博客的详细内容...",
        "author":  "作者名称",
        "created": "2024-03-20",
    }
    
    c.JSON(http.StatusOK, gin.H{
        "code": 0,
        "msg":  "success",
        "data": blog,
    })
} 