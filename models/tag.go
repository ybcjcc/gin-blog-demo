package models

import (
	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	Posts       []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

type PostTag struct {
	PostID uint `gorm:"primary_key"`
	TagID  uint `gorm:"primary_key"`
}