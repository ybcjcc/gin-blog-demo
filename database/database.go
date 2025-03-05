package database

import (
	"fmt"
	"blog/config"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.Mysql.User,
		config.AppConfig.Mysql.Password,
		config.AppConfig.Mysql.Host,
		config.AppConfig.Mysql.Port,
		config.AppConfig.Mysql.DBName,
	)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	DB = db
	// 设置连接池
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	// 启用日志
	DB.LogMode(true)

	return nil
}