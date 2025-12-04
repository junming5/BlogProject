# Go Gin + GORM 博客系统后端

## 📜 1. 项目概述

本项目是一个基于 Go 语言构建的博客系统后端服务。使用 Gin 框架处理路由和 HTTP 请求，GORM 作为 ORM 负责数据库操作。项目实现了用户身份认证、文章管理（CRUD）和评论功能。

## 📦 2. 提交内容说明

本次提交包含以下完整项目代码及必要配置：

1. **完整的项目代码文件**：包含所有业务逻辑和路由处理代码。

2. **依赖管理文件 (`go.mod` 和 `go.sum`)**：确保依赖可重现安装。

3. **项目说明文档 (`README.md`)**：即本文档，提供所有运行、配置和测试指南。

## ⚙️ 3. 运行环境与依赖

运行本项目需要以下环境和工具：

1. **Go 语言环境**：Go 1.21 或更高版本。

2. **数据库**：MySQL 5.7 或 8.0 版本。

3. **接口测试工具**：Postman、cURL 或其他 HTTP 客户端工具。

**项目主要依赖库：**

* `github.com/gin-gonic/gin`: Web 框架

* `gorm.io/gorm`: ORM 库

* `gorm.io/driver/mysql`: MySQL 驱动

* `github.com/golang-jwt/jwt/v5`: JWT 认证

## 🚀 4. 安装、配置与启动步骤

### 4.1 数据库配置

1. **创建数据库**：在 MySQL 中创建目标数据库，例如 `blog_db`。

   ```
   CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   
   ```

2. **修改连接配置**：请在项目主文件（通常是 `main.go`）中，修改数据库连接字符串 `dsn` 为您实际的 MySQL 用户名和密码。

   > **示例（请替换 `[username]` 和 `[password]`）：**

   > ```
   > dsn := "[username]:[password]@tcp(127.0.0.1:3306)/blog_db?charset=utf8mb4&parseTime=True&loc=Local"
   > 
   > ```

### 4.2 依赖安装

在项目根目录下，执行 Go 依赖安装命令：

```
go mod tidy

```

### 4.3 启动服务

执行以下命令启动后端服务：

```
go run main.go

```

服务成功启动后，控制台将显示连接成功的日志信息，且服务默认监听 `http://localhost:8080` 端口。GORM 将自动执行数据库迁移并创建所需的数据表（`users`, `posts`, `comments`）。

## 🧪 5. API 接口测试用例与结果

所有 API 接口的基础 URL 为：`http://localhost:8080`。

### 5.1 用户认证模块 (`/api/auth`)

| 序号 | 接口名称 | 路径 | 方法 | Body 示例 | 预期状态码 | 预期结果/说明 | 
 | ----- | ----- | ----- | ----- | ----- | ----- | ----- | 
| **1.1** | 用户注册 | `/api/auth/register` | `POST` | `{"username": "testuser", "password": "password123", "email": "test@example.com"}` | `201 Created` | 成功创建新用户，返回用户ID。 | 
| **1.2** | 用户登录 | `/api/auth/login` | `POST` | `{"username": "testuser", "password": "password123"}` | `200 OK` | 登录成功，返回包含 JWT **Token** 的响应体。**（此 Token 用于后续认证）** | 
| **1.3** | 登录失败 | `/api/auth/login` | `POST` | `{"username": "testuser", "password": "wrongpassword"}` | `401 Unauthorized` | 密码错误，认证失败。 | 

### 5.2 文章管理模块 (`/api/v1/posts`)

所有需要认证的接口，请求头必须携带：`Authorization: Bearer <TOKEN>`。

| 序号 | 接口名称 | 路径 | 方法 | 权限 | Body 示例 | 预期状态码 | 预期结果/说明 | 
 | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | 
| **2.1** | 创建文章 | `/api/v1/posts` | `POST` | 认证 | `{"title": "项目测试文章", "content": "这是用于接口测试的示例内容。"}` | `201 Created` | 返回文章 ID 和创建成功信息。 | 
| **2.2** | 获取所有文章 | `/api/v1/posts` | `GET` | 公开 | 无 | `200 OK` | 返回文章列表（包含作者信息）。 | 
| **2.3** | 获取单篇文章 | `/api/v1/posts/1` | `GET` | 公开 | 无 | `200 OK` | 返回 ID 为 1 的文章详情（包含评论列表）。 | 
| **2.4** | 更新文章 | `/api/v1/posts/1` | `PUT` | 作者认证 | `{"title": "更新后的标题", "content": "内容已修改！"}` | `200 OK` | 文章内容更新成功。 | 
| **2.5** | 删除文章 | `/api/v1/posts/1` | `DELETE` | 作者认证 | 无 | `204 No Content` | 文章删除成功。 | 
| **2.6** | 权限拒绝测试 | `/api/v1/posts/1` | `PUT` | 认证 | `{"title": "..."}` | `403 Forbidden` | 非文章作者尝试更新，权限不足。 | 

### 5.3 评论模块 (`/api/v1/posts/:post_id/comments`)

| 序号 | 接口名称 | 路径 | 方法 | 权限 | Body 示例 | 预期状态码 | 预期结果/说明 | 
 | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | 
| **3.1** | 创建评论 | `/api/v1/posts/1/comments` | `POST` | 认证 | `{"content": "对文章 ID 1 的评论内容。"}` | `201 Created` | 成功发表评论。 | 
| **3.2** | 获取文章评论 | `/api/v1/posts/1/comments` | `GET` | 公开 | 无 | `200 OK` | 返回文章 ID 1 下的评论列表。 | 
| **3.3** | 删除评论 | `/api/v1/comments/1` | `DELETE` | 评论作者认证 | 无 | `204 No Content` | 成功删除 ID 为 1 的评论。 | 

### 5.4 错误处理测试

| 序号 | 测试场景 | 路径 | 方法 | 预期状态码 | 预期结果/说明 | 
 | ----- | ----- | ----- | ----- | ----- | ----- | ----- | 
| **4.1** | 资源未找到 | `/api/v1/posts/9999` | `GET` | `404 Not Found` | 数据库中不存在的资源。 | 
| **4.2** | 认证信息缺失 | `/api/v1/posts` | `POST` | `401 Unauthorized` | 缺少 Token 或 Token 无效。 | 
| **4.3** | 字段验证失败 | `/api/auth/register` | `POST` | `400 Bad Request` | Body 中缺少必填字段或格式错误。 | 

**特别说明：** 测试时请严格按照 **步骤 5.1 -> 5.2 -> 5.3** 的顺序进行，以确保获取到有效的 JWT Token 和文章 ID。