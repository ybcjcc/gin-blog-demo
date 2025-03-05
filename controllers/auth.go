package controllers

import (
	"blog/database"
	"blog/middleware"
	"blog/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_type": "invalid_request", "message": "无效的请求参数"})
		return
	}

	// 密码复杂度验证
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error_type": "weak_password", "message": "密码长度至少8位"})
		return
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error_type": "weak_password", "message": "必须包含至少一个大写字母"})
		return
	}
	if !regexp.MustCompile(`\d`).MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error_type": "weak_password", "message": "必须包含至少一个数字"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被使用"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     "user", // 默认角色
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"redirect": "/login",
		"user": gin.H{
			"id": user.ID,
			"username": user.Username,
			"email": user.Email
		}
	})
}