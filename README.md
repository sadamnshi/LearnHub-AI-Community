# 🤖 LearnHub AI Community (AI增强型学习社区平台)

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.x-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

LearnHub AI Community 是一个基于 **Golang (Gin)** 和 **Vue.js** 构建的**AI增强型学习社区平台**。该项目旨在为开发者和学习者提供一个集技术交流、知识分享与 AI 智能辅助于一体的现代化社区。不仅有传统的帖子发布与互动功能，还深度集成了 AI 助手（支持多轮对话上下文），随时解答学习疑惑。

---

## ✨ 核心功能 (Core Features)

- 🔐 **完整的用户系统**：基于 JWT (JSON Web Token) 的高安全性用户注册、登录、鉴权与身份信息管理。
- 🤖 **智能 AI 助手 (AI Chat)**：
  - 支持**多轮上下文对话**，拥有强大的代码与知识辅导能力。
  - **高并发保护**：后端引入 WaitGroup 异步加载历史请求与库表更新。
  - **API 限流防护**：基于 Go Channel 实现 Semaphore（信号量）限流，保护外部大模型 API 不被超刷。
  - **可靠的异步落库**：采用**指数退避重试 (Exponential Backoff Retry)** 算法，保障 AI 聊天记录在数据库高负载下也能安全持久化。
- 📝 **社区交流论坛 (Posts)**：提供帖子发布、详情查看、分享等功能，打造高互动性的学习社区。
- 🚀 **高性能缓存架构**：深度集成 **Redis**，对社区热点数据和高频查阅内容进行缓存（如 `post_cache`），大幅提升系统吞吐量。

---

## 🛠️ 技术栈 (Tech Stack)

### 👨‍💻 前端 (Frontend : `vue-demo`)
- **核心框架**: Vue.js (Vue CLI / UI)
- **路由管理**: Vue Router (`vue-router`)
- **网络请求**: Axios (封装了统一的 `token` 拦截与错误处理)
- **本地存储**: localStorage / SessionStorage (管理 JWT Token 与用户状态)

### ⚙️ 后端 (Backend : `go-serve`)
- **核心框架**: Gin (高性能 Go Web 框架)
- **ORM与数据库**: GORM + PostgreSQL (结构化数据持久存储)
- **缓存与中间件**: Redis (加速读写与状态管理)
- **AI 接入**: `go-openai` (兼容 OpenAI / 阿里云百炼大模型 API)
- **跨域与鉴权**: 自定义 CORS 中间件 + 高度封装的 JWT Auth 拦截器

---

## 🏗️ 项目架构 (Project Architecture)

整个项目采用经典的 **前后端分离** 架构与 **三层 MVC 架构**：

```text
├── go-serve/ (Golang 后端服务)
│   ├── config/       # 配置文件 (.env 环境变量管理)
│   ├── controllers/  # 控制层 (接收请求，参数校验：处理 AI、Post、User)
│   ├── services/     # 业务逻辑层 (复杂的并发流水线、限流与指数重试)
│   ├── repositories/ # 数据访问层 (封装数据库 CRUD 操作)
│   ├── models/       # 数据模型 (pg 表映射)
│   ├── databases/    # 中间件连接池 (pg、Redis 初始化)
│   ├── middleware/   # 核心中间件 (JWT 拦截鉴权、CORS 跨域处理)
│   └── cache/        # 缓存层 (结合 Redis 处理高频查询)
│
└── vue-demo/ (Vue.js 前端工程)
    ├── src/api/      # Axios 接口封装层 (按业务模块分发)
    ├── src/views/    # 页面级视图 (AIChat, Home, Profile, Login 等)
    ├── src/router/   # 页面路由配置
    └── src/utils/    # 工具类函数 (Token 解析验证等)
```

---

## 🚀 快速开始 (Quick Start)

### 1. 环境准备 (Prerequisites)
请确保你的本地或服务器环境中已安装以下组件：
- **Golang** (>= v1.18)
- **Node.js** (>= v16) & npm/yarn/pnpm
- **PostgreSQL** (>= 12.0)
- **Redis** (>= 6.0)

### 2. 后端部署 (Backend Setup)

1. **进入后端目录**:
   ```bash
   cd go-serve
   ```
2. **下载依赖**:
   ```bash
   go mod tidy
   ```
3. **环境配置**:
   将 `config/.env.example` 复制为 `config/.env`，并修改里面的配置：
   ```env
   # PostgreSQL 数据库配置
   DB_DSN="host=127.0.0.1 user=postgres password=你的密码 dbname=gin_demo port=5432 sslmode=disable"
   
   # Redis 配置
   REDIS_ADDR="127.0.0.1:6379"

   # AI 模型配置 (例如阿里云)
   ALIYUN_API_KEY="your_api_key_here"
   ALIYUN_BASE_URL="https://dashscope.aliyuncs.com/compatible-mode/v1"
   OpenAIModel="qwen-plus" # 或 gpt-4, gpt-3.5-turbo 等
   
   # 鉴权配置
   JWT_SECRET="YOUR_SUPER_SECRET_KEY"
   ```
4. **启动服务**:
   ```bash
   go run main.go
   ```
   *服务启动后将自动执行数据库迁移 (AutoMigrate) 并在 `http://localhost:8080` 暴露接口。*

### 3. 前端部署 (Frontend Setup)

1. **进入前端目录**:
   ```bash
   cd vue-demo
   ```
2. **安装依赖**:
   ```bash
   npm install
   ```
3. **环境配置**:
   如有必要，根据后端地址修改 `src/api/` 中的 `baseURL`（默认通常指向 `http://localhost:8080/api`）。
4. **启动开发服务器**:
   ```bash
   npm run serve
   ```
   *浏览器将自动打开前端页面，默认为 `http://localhost:8080` 或 `http://localhost:8081`。*

---

## 🤝 参与贡献 (Contributing)

欢迎任何人为这个项目做出贡献！如果你发现了 Bug 或者有新的特性需求，请：
1. Fork 这个仓库
2. 创建一个新的分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 发起一个 Pull Request

---

## 📄 许可证 (License)

本项目基于 [MIT License](LICENSE) 许可开源。欢迎个人及企业自由使用和在此基础上进行二次开发。