package models

import "gorm.io/gorm"

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;type:varchar(50)" json:"username"`
	Password string `gorm:"not null;type:varchar(255)" json:"password"`
	Email    string `gorm:"unique;not null;type:varchar(100)" json:"email"`
	Posts    []Post
	Comments []Comment
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