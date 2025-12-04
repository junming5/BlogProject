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
	Title    string    `gorm:"not null;type:varchar(255)" json:"title"`
	Content  string    `gorm:"not null;type:text" json:"content"`
	UserID   uint      `json:"user_id"`                         // 外键关联 User
	User     User      `gorm:"foreignKey:UserID" json:"author"` // GORM 关联对象
	Comments []Comment `json:"comments"`
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null;type:text" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
	PostID  uint   `json:"post_id"`
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

// PostRequest 用于接收创建和更新文章请求的输入
type PostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// --- 认证 Handler (Auth Handlers) ---

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

// --- 认证中间件 (Middleware) ---

// AuthRequired 是一个 Gin 中间件，用于验证请求中的 JWT Token
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Header 中获取 Token: Authorization: Bearer <token>
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort() // 终止后续操作
			return
		}

		// 提取实际的 Token 字符串
		tokenString = tokenString[7:]

		// 2. 解析和验证 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 确保签名方法是 HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
			}
			return jwtSecret, nil // 使用全局密钥进行验证
		})

		// 3. 检查解析结果
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 4. 将用户信息（如 UserID）存储在 Context 中，供后续 Handler 使用
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["user_id"].(float64)) // JWT number claims are float64
			c.Set("user_id", userID)
			c.Set("username", claims["username"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims invalid"})
			c.Abort()
			return
		}

		// Token 验证通过，继续处理请求
		c.Next()
	}
}

// --- 文章 CRUD Handlers ---

// CreatePost 处理创建新文章的请求 (已更新，使用 PostRequest DTO)
func CreatePost(c *gin.Context) {
	var input PostRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Binding error: %v", err.Error())})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	post := Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID.(uint),
	}

	if err := DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post_id": post.ID,
		"title":   post.Title,
	})
}

// GetPosts 处理获取所有文章列表的请求
func GetPosts(c *gin.Context) {
	var posts []Post
	// Preload("User") 确保同时加载关联的 User 信息
	// 忽略软删除的文章 (DeletedAt is NULL)
	if err := DB.Preload("User").Order("created_at desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost 处理获取单个文章详情的请求
func GetPost(c *gin.Context) {
	// 从 URL 参数获取文章 ID
	id := c.Param("id")
	var post Post

	// Preload("User") 和 Preload("Comments")
	if err := DB.Preload("User").Preload("Comments").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost 处理更新文章的请求
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint) // 从中间件获取当前用户ID

	// 1. 查找文章并检查作者
	var post Post
	if err := DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 2. 授权检查：确保当前用户是文章作者
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied: You are not the author of this post"})
		return
	}

	// 3. 绑定更新数据
	var input PostRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Binding error: %v", err.Error())})
		return
	}

	// 4. 更新字段并保存
	DB.Model(&post).Updates(map[string]interface{}{
		"Title":   input.Title,
		"Content": input.Content,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// DeletePost 处理删除文章的请求
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(uint) // 从中间件获取当前用户ID

	// 1. 查找文章并检查作者
	var post Post
	if err := DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 2. 授权检查：确保当前用户是文章作者
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied: You are not the author of this post"})
		return
	}

	// 3. 删除文章 (GORM 的 gorm.Model 提供了软删除功能)
	if err := DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
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

	// --- 文章公开读取路由组 (Public Posts Group) ---
	// 获取列表和详情不需要认证
	postsPublic := r.Group("/api/v1/posts")
	{
		postsPublic.GET("", GetPosts)    // GET /api/v1/posts -> 获取所有文章列表
		postsPublic.GET("/:id", GetPost) // GET /api/v1/posts/:id -> 获取单个文章详情
	}

	// --- 受保护的文章操作路由组 (Protected Posts Group) ---
	// 创建、更新、删除需要 JWT 认证
	protected := r.Group("/api/v1/posts")
	protected.Use(AuthRequired()) // 应用认证中间件
	{
		protected.POST("", CreatePost)       // POST /api/v1/posts -> 创建文章
		protected.PUT("/:id", UpdatePost)    // PUT /api/v1/posts/:id -> 更新文章
		protected.DELETE("/:id", DeletePost) // DELETE /api/v1/posts/:id -> 删除文章
	}

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
