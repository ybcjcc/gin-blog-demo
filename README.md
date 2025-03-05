# 博客系统

基于Golang和Gin框架构建的现代化博客平台，提供完整的用户认证、文章管理和内容展示功能。

## 功能特性

✅ 用户注册与登录认证  
✅ 文章发布与管理  
✅ 分类与标签系统  
✅ 响应式网页设计  
✅ RESTful API接口

## 技术栈

- **后端框架**: Gin 1.9.1
- **数据库**: MySQL 8.0
- **ORM**: GORM 1.25.7
- **前端**: Bootstrap 5.1.3
- **会话管理**: Gorilla Sessions

## 快速开始

### 环境要求

- Go 1.20+ 
- MySQL 8.0+
- Node.js 16+ (前端依赖)

### 安装步骤

1. 克隆仓库：
```bash
git clone https://github.com/your-repo/blog-system.git
```

2. 初始化配置：
```bash
cp config/config.example.yaml config/config.yaml
```

3. 安装依赖：
```bash
go mod tidy
```

4. 数据库初始化：
```bash
mysql -u root -p < database/init.sql
```

5. 启动服务：
```bash
go run main.go
```

## 项目结构

```
blog-system/
├── config/            # 配置文件
├── controllers/       # 控制器层
├── database/          # 数据库脚本
├── middleware/        # 中间件
├── models/            # 数据模型
├── static/            # 静态资源
├── views/             # 前端模板
└── main.go            # 入口文件
```

## 配置说明

编辑 `config/config.yaml`：

```yaml
database:
  host: localhost
  port: 3306
  user: blog
  password: secure_password
  name: blog_db

server:
  port: 8080
  session_secret: your_session_secret
```

## API文档

### 用户认证

**注册接口**  
`POST /api/auth/register`

```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "SecurePass123!"
}
```

**登录接口**  
`POST /api/auth/login`

### 文章管理

**创建文章**  
`POST /api/posts` (需认证)

**获取文章列表**  
`GET /api/posts?page=1&size=10`

## 许可证

MIT License