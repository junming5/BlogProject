package models

import "gorm.io/gorm"

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null;type:text" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `gorm:"foreignKey:UserID" json:"user"`
	PostID  uint   `json:"post_id"`
	Post    Post
}

// CommentRequest 用于接收创建评论请求的输入
type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}
