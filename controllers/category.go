package controllers

import (
	"blog/database"
	"blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func CreateCategory(c *gin.Context) {
	// 检查用户权限
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "category": category})
}

func GetCategory(c *gin.Context) {
	var category models.Category
	if err := database.DB.First(&category, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func ListCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func UpdateCategory(c *gin.Context) {
	// 检查用户权限
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	var category models.Category
	if err := database.DB.First(&category, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	category.Name = req.Name
	category.Description = req.Description

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "category": category})
}

func DeleteCategory(c *gin.Context) {
	// 检查用户权限
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
		return
	}

	var category models.Category
	if err := database.DB.First(&category, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	// 检查分类下是否有文章
	var count int64
	database.DB.Model(&models.Post{}).Where("category_id = ?", category.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下还有文章，无法删除"})
		return
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}