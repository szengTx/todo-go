# 个人任务管理系统（Go + Gin）web课程作业

一个基于 Go 语言与 Gin 框架的轻量级任务管理应用，支持用户注册/登录、任务创建与完成、截止日期、日历视图以及前端美化。

## 亮点功能（含前端交互）
- 账户体系：注册 / 登录 / 登出，登录后导航自动切换“任务 / 个人主页 / 登出”。
- 任务管理：创建、完成切换、删除；支持“全天”或“精确到分钟”的截止时间。
- 日历视图：FullCalendar 月/周/列表视图，非全天任务在周视图按时间段展示，事件数据来自 `/tasks/json`。
- 个人主页：可编辑显示名 / 邮箱 / 头像；支持拖拽或文件选择上传头像（存储到 `static/uploads` 并实时预览）。
- 前端动效：
	- 全局淡入/滑入（`.fade-in` / `.slide-in-up`）与滚动显隐（`sr-hidden/sr-show`）。
	- 按钮/导航水波纹点击、卡片与任务项轻微倾斜悬停。
	- 回到顶部按钮平滑滚动。
- 响应式布局：顶部导航、卡片和表单在窄屏自动换行，任务列表网格化显示。

## 前端实现要点
- 样式：`static/css/style.css`
	- 渐变按钮、卡片阴影、任务条目 hover 抬升。
	- 表单行内布局（日期+时间+全天勾选），日历容器样式，头像上传区域的拖拽高亮。
	- 动效辅助类：`.fade-in`, `.slide-in-up`, `.sr-hidden/.sr-show`, `.ripple`, `.back-to-top`。
- 脚本：`static/js/app.js`
	- 回到顶部按钮创建与显隐控制，平滑滚动。
	- Scroll reveal：进入视口时为元素添加 `enter/sr-show`。
	- Ripple：为按钮/提交/导航/视图切换添加水波纹点击动画。
	- Hover tilt：卡片与任务项跟随鼠标轻微倾斜。
- 日历：FullCalendar（CDN）在 `templates/tasks.html` 初始化，events 来源 `/tasks/json`，非全天事件附带 `start/end`（默认 1 小时）以便 timeGrid 定位。
- 模板：
	- `templates/base.html`：公共骨架，加载 `style.css`、FullCalendar、`app.js`。
	- `templates/tasks.html`：任务表单（全天/日期/时间）、列表视图与日历视图切换，截止时间按全天或具体时间展示。
	- `templates/profile.html`：头像拖拽/选择上传（POST `/profile/upload`），实时预览并同步表单；编辑显示名/邮箱/头像 URL。

## 后端要点（简述）
- 路由：
	- 公共：`/register`, `/login`, `/logout`
	- 任务组（鉴权）：`/tasks`, `/tasks/create`, `/tasks/toggle/:id`, `/tasks/delete/:id`, `/tasks/json`
	- 个人主页（鉴权）：`/profile` GET/POST，`/profile/upload` 头像上传
- 数据模型：
	- `Task`：`Title`、`Completed`、`Deadline *time.Time`、`AllDay bool`
	- `User`：`Username`、`Password`（bcrypt）、`DisplayName`、`Email`、`AvatarURL`
- 存储：SQLite（`tasks.db`），GORM 自动迁移。
- 静态资源：`/static` 托管 CSS/JS/上传文件。

## 快速开始
```sh
go run main.go
# 浏览器访问 http://localhost:8080
```

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

