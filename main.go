package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User 用户模型，用于认证
type User struct {
	gorm.Model
	Username string    `gorm:"unique;not null;type:varchar(50)"`
	Password string    `gorm:"not null;type:varchar(255)"` // 密码通常需要较长的字段来存储哈希值
	Email    string    `gorm:"unique;not null;type:varchar(100)"`
	Posts    []Post    // 关联用户发布的文章
	Comments []Comment // 关联用户发表的评论
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title    string    `gorm:"not null;type:varchar(255)"`
	Content  string    `gorm:"not null;type:text"` // 内容使用TEXT类型
	UserID   uint      // 外键关联 User
	User     User      // GORM 关联对象
	Comments []Comment // 关联文章下的评论
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null;type:text"`
	UserID  uint   // 评论者ID
	User    User
	PostID  uint // 所属文章ID
	Post    Post
}

func main() {
	// 数据库初始化函数
	InitDB()

	// Gin 框架初始化
	r := gin.Default()

	// 简单的测试路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Blog System Backend!",
			"status":  "Server is running (MySQL)",
		})
	})

	// 运行服务器
	log.Println("服务器正在运行在 :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}

// InitDB 初始化数据库连接
func InitDB() {
	// 🚨 数据库连接字符串 (DSN)
	// 格式：用户名:密码@tcp(主机地址:端口)/数据库名称?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:gormpass@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error

	// 连接到 MySQL 数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到 MySQL 数据库: %v", err)
	}

	fmt.Println("成功连接到 MySQL 数据库!")

	// 自动迁移/同步 所有结构体到数据库表
	// 这将在数据库中创建 users, posts, comments 三个表
	err = DB.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
}
