package controllers

import (
	"blog/database"
	"blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CreatePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CategoryID uint   `json:"category_id" binding:"required"`
	Tags       string `json:"tags"`
	Status     string `json:"status" binding:"required,oneof=draft published"`
}

func CreatePost(c *gin.Context) {
	userID := c.GetUint("userID")
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	post := models.Post{
		Title:      req.Title,
		Content:    req.Content,
		UserID:     userID,
		CategoryID: req.CategoryID,
		Tags:       req.Tags,
		Status:     req.Status,
		PublishAt:  time.Now(),
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "post": post})
}

func GetPost(c *gin.Context) {
	var post models.Post
	if err := database.DB.Preload("User").Preload("Category").Preload("Comments").First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func ListPosts(c *gin.Context) {
	var posts []models.Post
	query := database.DB.Preload("User").Preload("Category")

	// 分页
	page := 1
	pageSize := 10
	offset := (page - 1) * pageSize

	// 查询条件
	if category := c.Query("category"); category != "" {
		query = query.Where("category_id = ?", category)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 排序
	query = query.Order("created_at DESC")

	var total int64
	query.Model(&models.Post{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"posts": posts,
	})
}

func UpdatePost(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var post models.Post
	if err := database.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	if post.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此文章"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 更新文章
	post.Title = req.Title
	post.Content = req.Content
	post.CategoryID = req.CategoryID
	post.Tags = req.Tags
	post.Status = req.Status

	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "post": post})
}

func DeletePost(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var post models.Post
	if err := database.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	if post.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文章"})
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}