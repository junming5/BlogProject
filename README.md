⸻

Go Gin + GORM 博客系统后端

📜 1. 项目概述

本项目是一个基于 Go 语言构建的博客系统后端服务：
	•	使用 Gin 处理路由和 HTTP 请求
	•	使用 GORM 操作数据库
	•	实现了 用户认证、文章管理（CRUD）、评论功能

⸻

📦 2. 提交内容说明

本次提交包含以下内容：
	•	完整项目代码（业务逻辑、控制器、模型、路由等）
	•	依赖管理文件：go.mod、go.sum
	•	项目说明文档：README.md
	•	数据库自动迁移（使用 GORM AutoMigrate）

⸻

⚙️ 3. 运行环境与依赖

运行本项目需要：
	•	Go 1.21+
	•	MySQL 5.7 或 8.0
	•	Postman / Curl / Apifox
	•	项目依赖：
	•	github.com/gin-gonic/gin
	•	gorm.io/gorm
	•	gorm.io/driver/mysql
	•	github.com/golang-jwt/jwt/v5

⸻

🚀 4. 安装、配置与启动步骤

4.1 数据库配置
	1.	创建数据库：

CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

	2.	在 main.go 修改数据库连接：

dsn := "[username]:[password]@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"


⸻

4.2 安装依赖

go mod tidy


⸻

4.3 启动服务

go run main.go

服务默认运行在：

http://localhost:8080

启动后 GORM 会自动创建表：
	•	users
	•	posts
	•	comments

⸻

🧪 5. API 接口测试

基础 URL：

http://localhost:8080

⸻

5.1 用户认证模块 (/api/auth)

序号	接口名称	路径	方法	Body 示例	状态码	说明
1.1	用户注册	/api/auth/register	POST	{"username":"testuser","password":"password123","email":"test@example.com"}	201	创建成功
1.2	用户登录	/api/auth/login	POST	{"username":"testuser","password":"password123"}	200	返回 JWT Token
1.3	登录失败	/api/auth/login	POST	{"username":"testuser","password":"wrongpassword"}	401	密码错误


⸻

5.2 文章管理模块 (/api/v1/posts)

认证接口需要 Header：

Authorization: Bearer <TOKEN>

序号	接口名称	路径	方法	权限	Body 示例	状态码	说明
2.1	创建文章	/api/v1/posts	POST	认证	{"title":"项目测试文章","content":"示例内容"}	201	返回文章 ID
2.2	获取所有文章	/api/v1/posts	GET	公开	—	200	返回文章列表
2.3	获取文章详情	/api/v1/posts/1	GET	公开	—	200	包含评论
2.4	更新文章	/api/v1/posts/1	PUT	作者认证	{"title":"更新后的标题","content":"内容已修改！"}	200	更新成功
2.5	删除文章	/api/v1/posts/1	DELETE	作者认证	—	204	删除成功
2.6	权限不足	/api/v1/posts/1	PUT	认证	{"title":"..."}	403	非作者无法操作


⸻

5.3 评论模块 (/api/v1/posts/:post_id/comments)

序号	接口名称	路径	方法	权限	Body 示例	状态码	说明
3.1	创建评论	/api/v1/posts/1/comments	POST	认证	{"content":"对文章的评论"}	201	创建成功
3.2	获取评论列表	/api/v1/posts/1/comments	GET	公开	—	200	返回评论列表
3.3	删除评论	/api/v1/comments/1	DELETE	评论作者认证	—	204	删除成功


⸻

5.4 错误处理测试

场景	路径	方法	状态码	说明
资源不存在	/api/v1/posts/9999	GET	404	无此资源
缺少 Token	/api/v1/posts	POST	401	未认证
字段验证失败	/api/auth/register	POST	400	缺少字段或格式错误


⸻

🔖 特别说明

测试顺序建议：
	1.	先进行 用户注册 / 登录 获取 Token
	2.	再进行 文章相关接口 测试
	3.	最后测试 评论模块

保证 Token、文章 ID、评论 ID 顺序正确。

⸻

如需，我也可以帮你：

✅ 生成项目目录结构
✅ 帮你补充 Swagger 文档
✅ 生成 Postman / Apifox 导入文件
✅ 提供完整后端代码模板

只需告诉我即可！