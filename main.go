package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5" // 引用 JWT 库
	"golang.org/x/crypto/bcrypt"       // 引用 bcrypt 库
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Global variable for the database connection and JWT Secret
var DB *gorm.DB

// 🚨🚨🚨 注意：在生产环境中，这个密钥必须通过环境变量安全地加载
var jwtSecret = []byte("your_super_secret_key_for_blog_system")

// --- 数据模型定义 (Models Definition) ---

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;type:varchar(50)" json:"username"`
	Password string `gorm:"not null;type:varchar(255)" json:"password"`
	Email    string `gorm:"unique;not null;type:varchar(100)" json:"email"`
	Posts    []Post
	Comments []Comment
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title    string `gorm:"not null;type:varchar(255)"`
	Content  string `gorm:"not null;type:text"`
	UserID   uint
	User     User
	Comments []Comment
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null;type:text"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}

// LoginRequest 专门用于接收登录请求的输入
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 专门用于接收注册请求的输入
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"` // 注册时 Email 必须
}

// --- 控制器函数 (Controller Handlers) ---

// Register 处理用户注册
func Register(c *gin.Context) {
	var input RegisterRequest
	// 使用 ShouldBindJSON 绑定输入数据，同时进行必要的验证
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名或邮箱是否已存在
	var existingUser User
	if DB.Where("username = ?", input.Username).Or("email = ?", input.Email).First(&existingUser).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建新用户对象
	user := User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword), // 存储加密后的密码
	}

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login 处理用户登录并返回 JWT
func Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		// 返回详细的错误信息，帮助我们定位是哪个字段的绑定出了问题
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Binding error: %v", err.Error())})
		return
	}

	var storedUser User
	// 根据用户名查找用户
	if err := DB.Where("username = ?", input.Username).First(&storedUser).Error; err != nil {
		// 统一返回 'Invalid username or password'，避免暴露是否存在该用户
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成 JWT Token
	claims := jwt.MapClaims{
		"user_id":  storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token 24小时后过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用全局密钥签名 Token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// --- 初始化与路由 (Initialization and Routing) ---

func main() {
	InitDB()

	r := gin.Default()

	// 公开路由 (无需认证)
	public := r.Group("/api/auth")
	{
		public.POST("/register", Register) // 用户注册接口
		public.POST("/login", Login)       // 用户登录接口
	}

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
	// 🚨 数据库连接字符串 (DSN) - 请确保已正确修改
	dsn := "root:gormpass@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到 MySQL 数据库: %v", err)
	}

	fmt.Println("成功连接到 MySQL 数据库!")

	// 自动迁移所有模型
	err = DB.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
}
