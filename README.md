# DevHelper

一个面向开发者的全功能工具平台，前后端分离架构，前端 React + TypeScript，后端 Go。

## 功能模块

### 已实现

- **JSON 工具** — 格式化、压缩、验证、格式转换（YAML/XML/TOML）、Schema 生成与验证、Diff 比较、JSONPath 查询、操作历史回退
- **用户系统** — 注册/登录（JWT 认证）、个人信息管理、管理员模式（用户增删改查）

### TODO

#### 编码与加密

- [ ] **Base64 编解码** — 文本/文件的 Base64 编码与解码，支持 URL-safe 变体
- [ ] **JWT 分析器** — 解析 JWT Token 结构（Header/Payload/Signature），验证签名，查看过期时间
- [ ] **哈希计算** — MD5、SHA1、SHA256、SHA512 等常用哈希算法
- [ ] **URL 编解码** — URL encode/decode，支持批量处理
- [ ] **加解密工具** — AES、RSA 等对称/非对称加密，密钥生成

#### 文档与格式

- [ ] **XML/HTML 查看与压缩** — 语法高亮、格式化、压缩、XPath 查询，HTML 结构树形预览
- [ ] **Markdown 预览** — 实时渲染预览，支持 GFM（GitHub Flavored Markdown），导出 HTML/PDF
- [ ] **CSV/Excel 查看** — 表格化展示 CSV 数据，支持排序和过滤
- [ ] **正则表达式测试** — 实时匹配高亮，支持多语言正则语法说明

#### 时间与数据转换

- [ ] **时间戳转换** — Unix 时间戳与可读时间互转，支持多时区，相对时间计算
- [ ] **进制转换** — 二进制、八进制、十进制、十六进制互转
- [ ] **颜色转换** — HEX、RGB、HSL、HSV 互转，颜色选择器
- [ ] **单位换算** — 存储大小、网络速度、长度、重量等常用单位换算

#### 文件操作

- [ ] **文件内容比较（Diff）** — 上传两个文件进行逐行对比，支持文本和代码文件，语法高亮差异
- [ ] **文件哈希校验** — 计算文件 MD5/SHA256，用于完整性验证
- [ ] **图片处理** — 格式转换、压缩、Base64 互转、EXIF 信息查看

#### 远程与网络

- [ ] **SSH 远程连接** — 基于 Web 的 SSH 终端，支持多会话，文件上传/下载（SFTP）
- [ ] **FTP 客户端/服务端** — 内置 FTP 客户端连接远程服务器；可启动本地 FTP 服务端对外共享文件
- [ ] **HTTP 请求测试** — 类 Postman 的 API 调试工具，支持保存请求历史
- [ ] **WebSocket 测试** — 连接 WebSocket 服务，发送/接收消息，查看帧详情
- [ ] **DNS 查询** — 查询域名的 A/AAAA/MX/TXT 等记录
- [ ] **网速测试** — 获取当前ip，并可以对当前网速进行测试
- [ ] **抓包工具集成** — 类似wireshark工具，对网络包进行抓取分析

#### 开发笔记

- [ ] **笔记读写** — Markdown 编辑器，支持标签分类、全文搜索
- [ ] **多端同步** — 支持同步到 GitHub / Gitee（通过 Git）或云存储（S3/OSS），可配置自动同步间隔
- [ ] **代码片段管理** — 保存常用代码片段，支持语言分类和搜索

#### 系统与运维

- [ ] **Cron 表达式解析** — 解析并可视化展示 Cron 表达式含义，生成下次执行时间列表
- [ ] **JSON Web Token 生成** — 配置 Header/Payload/Secret，生成标准 JWT
- [ ] **证书查看** — 解析 X.509 证书（PEM/DER），查看有效期、颁发者、SAN 等信息
- [ ] **二维码生成/解析** — 文本/URL 生成二维码，上传图片解析二维码内容
- [ ] **当前机器配置** — 查看当前机器配置，包括CPU、显卡、内存、主板、显示器、磁盘等

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18 + TypeScript + Vite |
| UI | Tailwind CSS + shadcn/ui |
| 编辑器 | Monaco Editor |
| 状态管理 | Zustand + TanStack Query |
| 后端 | Go 1.21+ + Gin |
| 数据库 | SQLite3（GORM） |
| 认证 | JWT（golang-jwt/jwt v5） |
| API 文档 | Swagger（Swaggo） |

## 快速开始

### 环境要求

- Go 1.23+
- Node.js 18+
- npm 9+

### 启动后端

```bash
cd backend
cp .env.example .env   # 按需修改配置
go run cmd/server/main.go
# 服务启动在 http://localhost:8080
# API 文档: http://localhost:8080/swagger/index.html
```

### 启动前端

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
# 访问 http://localhost:5173
```

### 使用 Makefile（推荐）

```bash
make dev        # 同时启动前后端
make backend    # 仅启动后端
make frontend   # 仅启动前端
```

## 项目结构

```
devhelper/
├── backend/                  # Go 后端
│   ├── cmd/server/main.go   # 入口
│   ├── internal/
│   │   ├── api/             # Handler + 路由 + 中间件
│   │   ├── config/          # 配置管理
│   │   ├── database/        # 数据库连接
│   │   ├── models/          # 数据模型
│   │   ├── repository/      # 数据访问层
│   │   ├── service/         # 业务逻辑
│   │   └── utils/           # 工具函数
│   └── docs/                # Swagger 文档
│
└── frontend/                 # React 前端
    └── src/
        ├── api/             # API 客户端
        ├── components/      # 组件（ui/layout/auth/json）
        ├── hooks/           # 自定义 Hooks
        ├── lib/             # 工具函数
        ├── pages/           # 页面
        └── store/           # 状态管理
```

## API 概览

| 分组 | 路径前缀 | 说明 |
|------|---------|------|
| 认证 | `/api/v1/auth` | 注册、登录、刷新 Token |
| JSON | `/api/v1/json` | 格式化、转换、Schema、Diff、查询 |
| 历史 | `/api/v1/history` | 操作历史回退记录 |
| Schema | `/api/v1/schemas` | 用户 Schema 管理 |
| 管理员 | `/api/v1/admin` | 用户管理（需 admin 角色） |

## 环境变量

### backend/.env

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DB_PATH` | `./devhelper.db` | SQLite 数据库路径 |
| `JWT_SECRET` | — | JWT 签名密钥（**必须修改**） |
| `JWT_ACCESS_EXPIRY` | `15m` | Access Token 有效期 |
| `JWT_REFRESH_EXPIRY` | `168h` | Refresh Token 有效期 |
| `SERVER_PORT` | `8080` | 服务端口 |
| `CORS_ORIGINS` | `http://localhost:5173` | 允许的跨域来源 |
| `ADMIN_INIT_EMAIL` | — | 首次启动时自动设为管理员的邮箱 |

### frontend/.env

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `VITE_API_BASE_URL` | `http://localhost:8080/api/v1` | 后端 API 地址 |
