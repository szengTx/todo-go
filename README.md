# 个人任务管理系统（Go + Gin）

一个基于 Go 语言与 Gin 框架的轻量级任务管理应用，支持用户注册/登录、任务创建与完成、截止日期、日历视图以及前端美化。

## 功能特性
- 用户注册 / 登录 / 登出
- 任务创建、完成状态切换、删除
- 支持可选截止日期，任务列表显示日期
- 日历视图（FullCalendar），可在月/周/列表视图中查看任务
- SQLite 数据持久化

## 技术栈
- Go 1.21+
- Gin Web 框架
- GORM + SQLite
- FullCalendar（前端日历）
- 自定义 CSS（`static/css/style.css`）

## 快速开始
```sh
go run main.go
```
访问：`http://localhost:8080`

## 目录结构
- `main.go`：路由与中间件注册、静态资源托管
- `handlers/`：HTTP 处理器（认证、任务 CRUD、任务 JSON 接口）
- `models/`：数据模型（`Task` 支持可选 `Deadline`）
- `database/`：数据库初始化与迁移（自动迁移 `User`、`Task`）
- `templates/`：页面模板（任务列表/日历、登录、注册、基础布局）
- `static/css/`：自定义样式

## 常见问题与解决方案
### 1) 页面无法正确显示（注册页显示成任务页或模板错乱）
**问题原因**：
- 早期使用 `r.Group("/")` 绑定认证中间件，导致 `/register`、`/login` 等公共路由也被鉴权拦截并重定向到任务页。
- 全局 `SetHTMLTemplate(ParseGlob("templates/*.html"))` 在多模板同名定义时会出现覆盖/冲突，导致渲染到错误模板。

**解决办法**：
- 将受保护路由改为 `r.Group("/tasks")` 并仅在该组上使用 `AuthMiddleware()`，公共路由保持开放。
- 各处理器内部使用 `template.ParseFiles` 按需解析所需模板文件，避免全局模板覆盖。

### 2) 任务标题删除线在编辑器中提示语法错误
**问题原因**：
- 在 HTML `style` 属性中直接嵌入 Go 模板条件（`{{if .Completed}}...{{end}}`）触发前端 linter 误报。

**解决办法**：
- 使用 CSS 类代替内联样式：`<span class="task-title {{if .Completed}}completed{{end}}">`，并在 `static/css/style.css` 中定义 `.task-title.completed { text-decoration: line-through; }`。

## 主要接口
- `GET /register` / `POST /register`：注册
- `GET /login` / `POST /login`：登录
- `GET /tasks`：任务列表 / 日历页面
- `POST /tasks/create`：创建任务（可选字段 `deadline`，格式 `YYYY-MM-DD`）
- `GET /tasks/json`：返回当前用户的任务 JSON，供日历视图使用

## 开发与调试
- 模板文件在 `templates/`，样式在 `static/css/style.css`，日历数据来自 `/tasks/json`
- 数据库存储文件默认 `tasks.db`（SQLite），在 `database.InitDB()` 中配置

## 许可
MIT
