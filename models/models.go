package models

import (
	"time"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"unique" json:"email"`
	Role     string `gorm:"default:'user'" json:"role"`
	Posts    []Post `json:"posts,omitempty"`
}

type Category struct {
	gorm.Model
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	Posts       []Post `json:"posts,omitempty"`
}

type Post struct {
	gorm.Model
	Title      string    `gorm:"not null" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user"`
	CategoryID uint      `json:"category_id"`
	Category   Category  `json:"category"`
	Tags       string    `json:"tags"`
	Status     string    `gorm:"default:'draft'" json:"status"`
	PublishAt  time.Time `json:"publish_at"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"user"`
	PostID  uint   `json:"post_id"`
	Post    Post   `json:"post"`
}

func InitDB(db *gorm.DB) error {
	// 自动迁移数据库表结构
	return db.AutoMigrate(&User{}, &Category{}, &Post{}, &Comment{}).Error
}