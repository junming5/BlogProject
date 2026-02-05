package models

import "gorm.io/gorm"

// Post 文章模型
type Post struct {
	gorm.Model
	Title    string    `gorm:"not null;type:varchar(255)" json:"title"`
	Content  string    `gorm:"not null;type:text" json:"content"`
	UserID   uint      `json:"user_id"`                         // 外键关联 User
	User     User      `gorm:"foreignKey:UserID" json:"author"` // GORM 关联对象
	Comments []Comment `json:"comments"`
}

// PostRequest 用于接收创建和更新文章请求的输入
type PostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
