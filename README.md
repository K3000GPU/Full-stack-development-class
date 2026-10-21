# 校园留言板

使用 Vue 3、Go/Gin 和 MySQL 的留言板，支持查看和发布留言。前端采用人物背景、黑白面板与斜切导航，数据保存在真实 MySQL 中。

## 功能

- 发布昵称和留言内容，成功后重新读取列表。
- 查看按 ID 倒序排列的留言及发布时间。
- 昵称去除首尾空白后为 1～30 个字符；内容不能为纯空白，原文最多 200 个字符。
- 中英文和 emoji 按 Unicode code point 计数；内容保留原始空白和换行。
- 提供加载状态、错误提示和重新加载，失败时保留表单输入。

## 环境与技术栈

| 项目 | 验证版本 |
| --- | --- |
| 操作系统 | Windows 11，PowerShell 7.6.5 |
| Node.js / npm | 24.21.0 / 11.19.0 |
| Go | 1.27.1 |
| Gin | 1.12.0 |
| MySQL 驱动 | go-sql-driver/mysql 1.10.1 |
| MySQL | Community Server 8.4.9 |
| Vue / Vite / GSAP | 3.5.44 / 8.3.4 / 3.15.0 |

前端使用 `package-lock.json` 固定依赖解析结果。后端由 `go.mod` 声明模块版本，`go.sum` 校验下载的模块内容；`go.sum` 不是与 npm lockfile 等价的锁文件。首次安装依赖需要联网。

## 项目结构

```text
campus-message-board/
├── frontend/
│   ├── src/views/MessageBoard.vue       页面状态与发布后刷新
│   ├── src/components/MessageForm.vue   输入与提交事件
│   ├── src/components/MessageList.vue   留言与加载、错误状态
│   ├── src/components/portrait/         人物布局与 GSAP 动效
│   ├── src/api/messages.js              listMessages、createMessage
│   ├── src/api/http.js                  HTTP 请求与错误处理
│   └── vite.config.js                  5173 端口与 /api 代理
├── backend/
│   ├── main.go                         配置、连接模块、注册路由
│   ├── go.mod / go.sum                  Go 模块与校验信息
│   ├── schema.sql                      创建缺失的 messages 表
│   └── internal/
│       ├── config/config.go            环境变量与数据库连接参数
│       ├── handler/message_handler.go  HTTP 解析与响应
│       ├── service/message_service.go  输入规则与业务过程
│       ├── repository/message_repository.go  参数化 SQL
│       └── model/message.go            记录、请求与响应结构
├── docs/                               调用链、UI 来源与验证说明
└── legacy/spring-boot/                  原后端归档，不参与默认运行
```

默认请求链：浏览器 → Vite `/api` 代理 → Gin Handler → Service → Repository → MySQL。

## 本地运行

以下命令使用 Windows PowerShell。`npm.cmd` 可避免 `npm.ps1` 受到脚本执行策略限制。

### 1. 获取代码

```powershell
$repoUrl = Read-Host '请输入项目 Git URL'
git clone $repoUrl campus-message-board
Set-Location .\campus-message-board
node --version
npm.cmd --version
go version
```

### 2. 初始化数据库

先确认要使用的 MySQL 实例、端口和账号。只在明确指定的项目实例上创建 `campus_message_board`，不要在未知实例或其他数据库中执行初始化。

```powershell
mysql --protocol=TCP -h 127.0.0.1 -P 3306 -u root -p
```

在 MySQL 客户端中执行，将密码占位符替换为本机专用密码：

```sql
CREATE DATABASE IF NOT EXISTS campus_message_board
  CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

CREATE USER IF NOT EXISTS 'campus_app'@'127.0.0.1'
  IDENTIFIED BY '请替换为本机项目密码';

GRANT SELECT, INSERT, CREATE ON campus_message_board.*
  TO 'campus_app'@'127.0.0.1';
```

已有项目数据库和专用账号时直接使用现有配置；`IF NOT EXISTS` 不会更改已有账号密码。示例账号用于本机 TCP 连接，远程 MySQL 需要按实际来源地址授权。

数据库由上述步骤预先创建。后端启动会检查连接，并执行嵌入程序的 `backend/schema.sql`；它只包含 `CREATE TABLE IF NOT EXISTS`，不会清空、重建已有表或创建其他数据库。

### 3. 启动后端

从项目根目录打开 PowerShell：

```powershell
Set-Location .\backend
$env:PORT = '8080'
$env:DB_HOST = '127.0.0.1'
$env:DB_PORT = '3306'
$env:DB_NAME = 'campus_message_board'
$env:DB_USERNAME = Read-Host '请输入项目数据库用户名'
$dbSecret = Read-Host '请输入项目数据库密码' -AsSecureString
$env:DB_PASSWORD = [System.Net.NetworkCredential]::new('', $dbSecret).Password
Remove-Variable dbSecret

go mod download
go run .
```

后端默认监听 `http://127.0.0.1:8080`。保持此终端运行。

| 环境变量 | 默认值 / 要求 |
| --- | --- |
| `PORT` | `8080` |
| `DB_HOST` | `127.0.0.1` |
| `DB_PORT` | `3306` |
| `DB_NAME` | 必填，项目数据库名 |
| `DB_USERNAME` | 必填，项目专用账号 |
| `DB_PASSWORD` | 必填，账号密码 |

配置由当前进程的环境变量读取，不自动加载 `.env`。环境变量只对当前终端及其子进程有效；本机 MySQL 使用 `3307` 等其他端口时修改 `DB_PORT`。

### 4. 启动前端

从项目根目录打开另一个 PowerShell：

```powershell
Set-Location .\frontend
npm.cmd ci
npm.cmd run dev
```

打开 [http://127.0.0.1:5173](http://127.0.0.1:5173)。Vite 固定监听 `5173`，将 `/api` 代理到 `http://127.0.0.1:8080`，不改写路径。修改后端端口时也需同步修改 `frontend/vite.config.js` 的代理目标。

前后端和 MySQL 都需要保持运行。在前后端终端按 `Ctrl+C` 结束程序，停止服务不会删除数据库中的记录。

## API

| 方法 | 路径 | 成功响应 |
| --- | --- | --- |
| GET | `/api/messages` | `200`，留言数组；没有记录时为 `[]` |
| POST | `/api/messages` | `201`，新建记录 |

POST 请求使用 `Content-Type: application/json`：

```json
{"nickname":"  小林  ","content":"周五有人一起去图书馆吗？"}
```

响应示例：

```json
{
  "id": 1,
  "nickname": "小林",
  "content": "周五有人一起去图书馆吗？",
  "createdAt": "2026-10-21T19:30:00"
}
```

`id` 和 `createdAt` 来自数据库。表字段为 `id`、`nickname`、`content`、`created_at`；JSON 使用 `createdAt`。

错误响应统一为 `{"message":"错误说明"}`。输入或 JSON 无效返回 `400`，非 JSON 请求返回 `415`，数据库连接不可用返回 `503`，其他服务错误返回 `500`。响应不包含数据库凭据、SQL 或异常堆栈。

## 构建与检查

```powershell
# 后端：从项目根目录开始。
Set-Location .\backend
go mod verify
go test ./...
go build -o bin/campus-message-board.exe .

# 前端。
Set-Location ..\frontend
npm.cmd ci
npm.cmd run build
```

Go 单元测试不连接数据库。可执行文件位于 `backend/bin`，前端构建产物位于 `frontend/dist`；运行后端二进制时同样需要配置数据库环境变量。

前后端已经启动时，从项目根目录执行真实接口检查：

```powershell
node .\scripts\verify-api.mjs
```

脚本默认请求 `http://127.0.0.1:5173`，覆盖开发代理；每次新增并保留 3 条检查记录，验证有效发布、Unicode 边界和无效输入，不删除已有数据。运行前确认服务连接的是预期项目数据库。设置 `$env:API_BASE = 'http://127.0.0.1:8080'` 可以直连后端。

也可在同一数据库中独立核对结果：

```sql
SELECT DATABASE(), VERSION();
SELECT id, nickname, content, created_at
FROM campus_message_board.messages
ORDER BY id DESC
LIMIT 10;
```

## 常见问题

| 现象 | 检查方式 |
| --- | --- |
| 找不到 `go` 或 `node` | 检查安装与 PATH，重新打开终端 |
| Go 依赖下载超时 | 检查网络；可在当前终端设置 `$env:GOPROXY = 'https://goproxy.cn'` 后重试 |
| 缺少数据库配置 | 在启动后端的同一个终端设置必填环境变量 |
| MySQL 拒绝连接或认证失败 | 核对实例、端口、库名、密码与账号来源地址 |
| 前端显示读取失败 | 确认后端已在 8080 启动，再检查后端日志与浏览器 Network |
| 端口被占用 | 查清占用进程，不要随意关闭未知服务 |

## 更多文档

- [工程分层与完整调用链](docs/TEACHING.md)
- [UI 复用与素材来源](docs/UI_REUSE.md)
- [验证记录](docs/VERIFICATION.md)

本机配置、密码、数据库数据、依赖目录及构建产物均由 `.gitignore` 排除，不随 Git 分发。
