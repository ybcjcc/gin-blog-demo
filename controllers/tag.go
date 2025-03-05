package controllers

import (
	"blog/database"
	"blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type TagRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func CreateTag(c *gin.Context) {
	// 检查用户权限
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	tag := models.Tag{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建标签失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "tag": tag})
}

func GetTag(c *gin.Context) {
	var tag models.Tag
	if err := database.DB.First(&tag, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

func ListTags(c *gin.Context) {
	var tags []models.Tag
	if err := database.DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取标签列表失败"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func UpdateTag(c *gin.Context) {
	// 检查用户权限
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	var tag models.Tag
	if err := database.DB.First(&tag, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	tag.Name = req.Name
	tag.Description = req.Description

	if err := database.DB.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "tag": tag})
}

func DeleteTag(c *gin.Context) {
	// 检查用户权限
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	var tag models.Tag
	if err := database.DB.First(&tag, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "标签不存在"})
		return
	}

	// 检查标签是否被文章使用
	var count int64
	database.DB.Model(&models.Post{}).Where("tags LIKE ?", "%"+tag.Name+"%").Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该标签已被文章使用，无法删除"})
		return
	}

	if err := database.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetPostTags(c *gin.Context) {
	var post models.Post
	if err := database.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	tags := strings.Split(post.Tags, ",")
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func UpdatePostTags(c *gin.Context) {
	userID := c.GetUint("userID")
	role := c.GetString("role")

	var post models.Post
	if err := database.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	if post.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此文章的标签"})
		return
	}

	var req struct {
		Tags []string `json:"tags" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 更新文章标签
	post.Tags = strings.Join(req.Tags, ",")

	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章标签失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "tags": req.Tags})
}