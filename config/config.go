package config

type Config struct {
	Mysql struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
	Server struct {
		Port int
	}
	JWT struct {
		Secret string
	}
}

var AppConfig Config

func init() {
	// 默认配置
	AppConfig.Mysql.Host = "localhost"
	AppConfig.Mysql.Port = 3306
	AppConfig.Mysql.User = "root"
	AppConfig.Mysql.Password = "root"
	AppConfig.Mysql.DBName = "blog"

	AppConfig.Server.Port = 8080

	AppConfig.JWT.Secret = "your-secret-key"
}