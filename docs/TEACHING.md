# 教学导读：跟踪一次发布留言

按 [README](../README.md) 启动 MySQL、Go/Gin 和 Vite。本课用 45 分钟认识工程分层与一次完整请求，不以继续学习某种后端框架为前提。

## 课前准备

统一 Docker Desktop：Windows 使用 WSL 2 Ubuntu 26.04.1 LTS（发行版名称 `Ubuntu`），通过 WSL Integration 连接 Desktop；Mac 直接使用本机 Terminal。按 [Windows 指南](PREPARATION.md) 或 [Mac 指南](PREPARATION_MACOS.md) 准备，再执行 [共用指南](DEVELOPMENT_DOCKER.md)。

Windows 在 WSL Ubuntu 的 ~/projects、Mac 在本机 ~/projects 克隆项目并配置自己的 .env，按[共用指南](DEVELOPMENT_DOCKER.md)执行 `docker compose up --build -d --wait --wait-timeout 600`。Node.js、Go、MySQL 随容器提供，数据库和账号在空卷首次启动时创建，messages 表由后端创建。安装、镜像下载、首次构建与验收在课前完成。镜像拉取超时时才切换到[本地镜像备用方案](IMAGE_DISTRIBUTION.md)，导入课程提供的完整镜像后继续验收。

准备完成的标准：三个服务 healthy，宿主机浏览器能发布并刷新读取，使用 `docker compose exec db mysql -u campus_app -p campus_message_board` 查到相同记录，down 后重新启动仍可读取。各自使用自己的数据库，不依赖共享数据库或同一局域网。

Windows 和 Mac 浏览器都使用 localhost:5173。编辑器与 Agent 打开同一份挂载源码，Windows 使用 WSL 工作区，Mac 使用本机工作区。可对照[验证记录](VERIFICATION.md)中的截图检查自己的运行结果。

## Docker 中的运行关系

`compose.yaml` 定义 frontend、backend、db 三个服务；Dockerfile.dev 准备工具、依赖与启动命令。前端发布本机回环地址的 5173，Vite 把 /api 转发到 backend:8080，后端访问 db:3306。后端 8080 仅映射本机回环地址，MySQL 不发布端口。

`db` 是 **database（数据库）** 的缩写，这里用作 MySQL 的 Compose 服务名。`DB_HOST: db` 表示后端通过这个服务名寻找 MySQL；`DB_NAME: campus_message_board` 则指定连接后使用的数据库。

容器内监听 0.0.0.0，服务间用 Compose 服务名访问，每个容器的 localhost 都代表自己。Handler、Service、Repository 是一个后端程序的三个代码职责，并不是三个容器。

.env 由 Compose 读取再传给相应进程，Go 程序本身只读环境变量。前端不接收数据库密码。mysql_data 命名卷保存数据库文件；常规 down 保留卷。已有原生 MySQL 数据不会自动复制到新卷，已有卷的密码也不会随 .env 修改而自动改变。

## 目录与职责

| 文件 | 职责 |
| --- | --- |
| `frontend/src/components/MessageForm.vue` | 管理输入、校验并发出 submit 事件 |
| `frontend/src/views/MessageBoard.vue` | 页面状态、发布过程、发布后重新读取 |
| `frontend/src/components/MessageList.vue` | 展示记录、加载与错误状态，发出 retry 事件 |
| `frontend/src/components/portrait/PortraitLayout.vue` | 人物背景与动效，slot 保持业务组件实例 |
| `frontend/src/api/messages.js` | 定义 listMessages、createMessage |
| `frontend/src/api/http.js` | 发送 HTTP 并统一处理错误 |
| `compose.yaml`、`.env.example` | 服务组合、网络、数据卷与配置模板 |
| `frontend/Dockerfile.dev`、`backend/Dockerfile.dev` | 容器中的工具、依赖与启动命令 |
| `frontend/vite.config.js` | 5173 端口及 /api 转发 |
| `backend/main.go` | 启动、连接具体模块、注册路由 |
| `backend/internal/config/config.go` | 读取环境变量并构造数据库连接参数 |
| `backend/internal/handler/message_handler.go` | HTTP 请求解析、业务调用和响应 |
| `backend/internal/service/message_service.go` | 规则检查与业务过程 |
| `backend/internal/repository/message_repository.go` | 参数化 SQL 与数据库结果读取 |
| `backend/internal/model/message.go` | Message、请求、响应与错误结构 |
| `backend/schema.sql` | 在指定数据库创建缺失的表 |
| `backend/go.mod`、`backend/go.sum` | 模块版本与下载内容校验 |

每层只有一个具体实现，没有 ORM、通用接口、工厂或依赖注入框架。

## 从入口找到后端

打开 `backend/main.go`：`config.Load()` 读取配置，`sql.Open()` 创建连接池，`PingContext()` 检查指定 MySQL，随后执行嵌入的 schema.sql。没有连接成功就不会开放 HTTP 端口。

同一文件显式创建模块：

```go
messageRepository := repository.New(database)
messageService := service.New(messageRepository)
messageHandler := handler.New(messageService)
```

再注册请求入口：

```go
router.GET("/api/messages", messageHandler.List)
router.POST("/api/messages", messageHandler.Create)
```

这里只启动一个 Go 后端进程。模块名称说明职责，并不代表多个服务。

## 按一次发布的顺序定位

行号从 1 开始，后续改代码时应重新核对；也可直接搜索方法名称。

| 顺序 | 文件与方法位置 | 发生的事情 | 调用类型 |
| --- | --- | --- | --- |
| 1 | [MessageForm.vue](../frontend/src/components/MessageForm.vue) 第 16 行 `handleSubmit()` | 检查输入，发出 `submit` 事件 | Vue 组件事件，浏览器内部 |
| 2 | [MessageBoard.vue](../frontend/src/views/MessageBoard.vue) 第 30 行 `publishMessage(payload)` | 设置提交状态，调用 createMessage | JavaScript 内部调用 |
| 3 | [messages.js](../frontend/src/api/messages.js) 第 8 行 `createMessage(message)` → [http.js](../frontend/src/api/http.js) 第 2 行 `request()` | 准备 POST 路径、JSON 和请求头 | JavaScript 内部调用 |
| 4 | `http.js` 第 6 行 `fetch()` → [vite.config.js](../frontend/vite.config.js) 第 10 行转发 | 浏览器发往 5173，Vite 转发到 8080；路径仍为 /api/messages | **HTTP POST 与 HTTP 转发** |
| 5 | [main.go](../backend/main.go) 第 72 行 `router.POST()` | 根据方法和路径，调用 MessageHandler.Create | HTTP 路由分发 |
| 6 | [message_handler.go](../backend/internal/handler/message_handler.go) 第 40 行 `Create()` | ShouldBindJSON 解析为 CreateMessageRequest，调用 Service | Go 内部调用 |
| 7 | [message_service.go](../backend/internal/service/message_service.go) 第 31 行 `Create()`、第 39 行 `validateInput()` | 昵称去首尾空白；按 Unicode 字符数检查长度；拒绝空白内容，再调用 Repository | Go 内部调用 |
| 8 | [message_repository.go](../backend/internal/repository/message_repository.go) 第 38 行 `Insert()`、第 46 行 `ExecContext()` | 事务内执行 `INSERT ... VALUES (?, ?)`，输入作为参数传入 | **database/sql + MySQL 通信** |
| 9 | 同一文件第 50 行 `LastInsertId()`、第 55 行 `QueryRowContext()`、第 61 行 `Commit()` | 取得真实主键，按 id 回读并 Scan 到 Message；提交成功后才返回 | MySQL 通信，随后 Go 内部返回 |
| 10 | [message.go](../backend/internal/model/message.go) 第 25 行 `NewMessageResponse()` → Handler 的 `c.JSON()` | 转换 createdAt，返回刚保存的记录和 201 | Go 内部调用，随后 **HTTP 响应** |
| 11 | `MessageBoard.vue` 第 46 行 → `MessageForm.vue` 第 33 行 `clearContent()` | 收到 POST 成功后清空内容、保留昵称，调用 loadMessages | JavaScript 内部调用 |
| 12 | `MessageBoard.vue` 第 17 行 `loadMessages()` → `messages.js` 第 4 行 `listMessages()` → `http.js` 的 fetch | 重新请求完整列表 | **新的 HTTP GET /api/messages** |
| 13 | `main.go` 第 71 行 → `MessageHandler.List()` 第 27 行 → `MessageService.List()` 第 26 行 → `MessageRepository.FindAll()` 第 19 行 | QueryContext 执行 SELECT，逐行 Scan，Handler 返回响应数组和 200 | Go 内部调用 + MySQL 通信，随后 HTTP 响应 |
| 14 | `MessageBoard.vue` 第 21 行 → `MessageList.vue` 第 33 行 `v-for` | 使用 GET 结果更新状态并重新渲染 | Vue 响应式更新，浏览器内部 |

首次进入页面，`MessageBoard.vue` 第 66 行 `onMounted(refreshMessages)` 经第 59 行的 refreshMessages 读取列表，走第 12～14 步。

POST 内部为取得新记录而执行的 SELECT，属于该 POST 的后端过程。发布后页面重新请求列表，是另一次独立 GET。

## 区分通信与内部调用

```text
浏览器
  MessageForm → MessageBoard → messages.js → http.js
                 JavaScript 调用 / 组件事件
                            │
               HTTP POST 或 GET /api/messages
                            ↓
Vite :5173 ── HTTP 转发，保留路径 ──→ Go/Gin :8080
                                        Handler → Service → Repository
                                             Go 内部方法调用
                                                               │
                                                   database/sql + 参数化 SQL
                                                               ↓
                                             MySQL / campus_message_board.messages
```

页面 import 文件、组件 emit 事件，以及 Handler 调用 Service，都不会产生额外 HTTP 请求。Repository 使用 MySQL 协议与数据库通信，也不是 HTTP。浏览器 Network 不显示 Go 内部调用或 SQL。

页面导航的 `#publish`、`#messages` 只是锚点滚动。列表错误时点击“重新加载”则发出 retry 事件，由页面再次 GET。

## 规则、错误与 SQL

Service 的 `validateInput()` 检查昵称和内容，`utf8.RuneCountInString()` 与前端 Array.from 的长度规则一致。`trimWhitespace()` 采用与 JavaScript String.trim 一致的空白集合，包含 NBSP 和 FEFF。昵称归一化后保存；内容保留原文。

Handler 把 ValidationError 转为 400，把非法 JSON 与非字符串字段作为 400，把非 JSON 请求作为 415。数据库不可连接时返回 503，其他内部错误返回 500；详情只写服务端日志。

前端 request 抛出的 Error 由页面捕获并展示。读取失败与空列表分别显示；重试保留草稿。若 POST 成功而 GET 失败，会明确提示已发布，避免重复提交。

Repository 中重点看两条 SQL：

```sql
SELECT id, nickname, content, created_at FROM messages ORDER BY id DESC;

INSERT INTO messages (nickname, content) VALUES (?, ?);

```

`?` 表示参数占位符。昵称中的单引号是普通数据，不拼接成 SQL。新增后按 `WHERE id = ?` 回读数据库生成的时间。

| 表字段 | 含义 | JSON 字段 |
| --- | --- | --- |
| id | 自增主键，昵称可以重复 | id |
| nickname | 去掉首尾空白后的昵称 | nickname |
| content | 留言原文 | content |
| created_at | 数据库生成的发布时间 | createdAt |

Compose 在自己的 MySQL 卷中使用 campus_message_board 数据库。用 SELECT 对照 POST 返回的 id，再 down 并重新 up，检查卷中的记录仍可读取；不迁移或清空原生 MySQL 的已有数据。

## 45 分钟安排

| 时间 | 内容 | 检查点 |
| --- | --- | --- |
| 0～5 分钟 | 检查已有 clone、启动、查看成品 | 三个容器分别运行什么 |
| 5～10 分钟 | 目录、依赖、入口与配置 | Compose、main.js、main.go、依赖与环境变量 |
| 10～18 分钟 | 前端分工 | 表单事件、页面状态、API 模块 |
| 18～28 分钟 | Go 后端分工 | 显式路由、Handler、Service、Repository |
| 28～35 分钟 | 表、字段、主键、SQL 与数据卷 | INSERT 新增，SELECT 读取，重启保留 |
| 35～42 分钟 | 追踪一次发布 | Network、响应、数据库记录对应 |
| 42～45 分钟 | 用自己的话复述 | HTTP 与内部调用有什么区别 |

只需认识函数接收数据并返回结果、struct 描述数据、err 表示失败。指针、接口、并发、事务细节和框架内部机制不作为本课目标。

可以用三个问题检查：修改按钮、修改输入规则、修改查询排序，分别应先查看哪层；为什么前端不直接发 SQL；后端重启后为什么还能读到留言。

这是一节独立的全栈入门课。入口、配置、接口、数据保存以及进程、端口、日志的认识，也能用于理解模型服务与机器人系统，不以继续某个后端框架为前提。
