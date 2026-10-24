# 校园留言板

基于 Vue 3、Go/Gin 和 MySQL 的留言板，支持查看和发布留言。前端包含人物背景与 GSAP 动效，数据保存在真实 MySQL 中。

## 开发环境

统一使用 **Docker Desktop** 提供 Engine、CLI 与 Compose，项目的 Node.js、Go 和 MySQL 运行在 Linux 容器中，无需在宿主机另行安装这三个运行环境。

| 电脑系统 | 开发方式 | 准备文档 |
| --- | --- | --- |
| Windows | Docker Desktop + WSL 2 Ubuntu 26.04.1 LTS（发行版名称 `Ubuntu`），启用 WSL Integration | [Windows 环境准备](docs/PREPARATION.md) |
| macOS | Docker Desktop，直接使用 Mac Terminal 与本机源码 | [Mac 环境准备](docs/PREPARATION_MACOS.md) |

Windows 在 WSL Ubuntu 的 `~/projects` 保存源码并执行命令；Mac 直接使用本机的 `~/projects` 与 Terminal。编辑器与 Agent 打开容器实际挂载的同一份源码。两种平台的 Docker 命令均不加 sudo，浏览器统一访问 localhost。

| 服务 | 容器环境 | 职责 |
| --- | --- | --- |
| `frontend` | `node:24.22.1-bookworm-slim`；Vue 3.5.44 / Vite 8.3.4 / GSAP 3.15.0 | 前端开发服务，发布本机回环地址的 5173 端口 |
| `backend` | `golang:1.27.3-bookworm`；Gin 1.12.0 / MySQL 驱动 1.10.1 | HTTP API，8080 端口映射到本机回环地址 |
| `db` | `mysql:8.4.11` | 数据库存储；3306 仅供容器网络访问 |

`db` 是 **database（数据库）** 的缩写，在本项目中是 `compose.yaml` 为 MySQL 定义的服务名。后端通过这个服务名找到 MySQL，因此配置为 `DB_HOST: db`；实际存放留言的数据库名则是 `campus_message_board`。

上表版本来自当前 `compose.yaml`、Dockerfile 和依赖文件，表示配置目标，不代表已在 Windows、Mac 和各 CPU 架构上全部完成实机验收。Compose 不强制 CPU 架构；本地镜像包则必须与 Docker Engine 架构匹配。

前端依赖解析由 `package-lock.json` 固定；后端 `go.mod` 声明模块版本，`go.sum` 校验下载内容。镜像使用具体版本标签，不使用浮动的基础镜像版本。

## 快速启动

先按对应平台指南安装并启动 Docker Desktop。以下命令在 **Windows 的 WSL Ubuntu 终端**或 **Mac Terminal** 执行，Bash / zsh 均可；默认通过 Compose 下载所需镜像、构建并启动。镜像拉取超时时，可使用匹配项目版本与架构的本地镜像备用方案。

### 1. 获取代码

项目仓库：[Yhimhas/Full-stack-development-class](https://github.com/Yhimhas/Full-stack-development-class)。仓库根目录即留言板项目，克隆时保留默认目录名 `Full-stack-development-class`：

```bash
mkdir -p ~/projects

cd ~/projects

git clone https://github.com/Yhimhas/Full-stack-development-class.git

cd Full-stack-development-class

```

已有项目时进入已有目录，同步代码前先保存本机改动。确认 Docker Desktop 已就绪：

```bash
docker version
docker compose version
docker version --format '{{.Server.Os}}/{{.Server.Arch}}'
```

应能看到 Docker Client、Server、Compose 版本及 `linux/amd64` 或 `linux/arm64`。若没有 Server 信息，先检查 Docker Desktop、当前 context 和 Windows 的 WSL Integration。

### 2. 配置本机密码

> [!WARNING]
> **`.env` 是隐私配置文件，包含数据库密码。**
> 不要将 `.env` 提交到 Git 或推送到 GitHub，也不要把文件、内容或含密码的截图发给 Agent 或他人。
> `.gitignore` 已配置忽略 `.env`，请保留该忽略规则。

```bash
cp -n .env.example .env

chmod 600 .env

nano .env

```

`cp -n` 表示已有 `.env` 时跳过复制，保留其中的配置和密码；不存在时从模板创建。

为 `MYSQL_ROOT_PASSWORD`、`DB_PASSWORD` 自行设置密码，可以相同，但都不能为空，保留外层单引号。示例：`DB_PASSWORD='自行设置的项目密码'`。保存方式：`Ctrl+O`、回车、`Ctrl+X`。密码不含单引号和换行；可使用至少 16 位字母数字组合。

| 密码 | 对应 MySQL 账号 | 使用场景 |
| --- | --- | --- |
| `MYSQL_ROOT_PASSWORD`（管理密码） | `root` 管理员 | 数据库、账号和权限管理；仅传给 `db` 容器 |
| `DB_PASSWORD`（项目密码） | `campus_app` 应用账号 | 后端连接项目库，执行建表、查询和新增留言；手动用 `mysql -u campus_app -p` 查询时也输入这个密码 |

Compose 用 `DB_PASSWORD` 同时配置数据库容器的 `MYSQL_PASSWORD` 和后端的 `DB_PASSWORD`，使应用账号与后端使用同一个密码。Ubuntu 登录或 `sudo` 的 Linux 用户密码与这两项分别配置。已有数据卷时，修改 `.env` 不会自动修改数据库账号密码；详细说明见[密码配置指南](docs/DEVELOPMENT_DOCKER.md#3-配置本机密码)。

`.env` 由 Compose 读取，并通过 `environment` 将必要参数传给容器；Go 程序只读取进程环境变量。根密码仅传给 `db`，项目密码传给 `db` 与 `backend`，前端不接收数据库密码。`.env` 已被 Git 忽略。

### 3. 构建并启动

在项目根目录执行：

```bash
docker compose config --quiet

docker compose up --build -d --wait --wait-timeout 600

docker compose ps

```

三个服务应为 `healthy`；打开下面的接口地址会看到 JSON 数组，空库为 `[]`。首次构建需要下载镜像、npm 与 Go 依赖并编译 Go。`--wait-timeout 600` 是服务健康等待时间，不是镜像下载和构建的总超时。

若镜像拉取或构建器读取基础镜像信息时超时，按[本地镜像备用方案](docs/IMAGE_DISTRIBUTION.md)校验并导入完整镜像包，再在项目根目录执行：

```bash
docker compose config --images
docker compose up --no-build --pull never -d --wait --wait-timeout 600
docker compose ps
```

当前文档中的镜像包为 `campus-message-board-images-linux-amd64.tar`，配套 `SHA256SUMS`，不包含在 Git 仓库内；仅适用于 `linux/amd64`。`linux/arm64` 使用主流程或对应架构的镜像包。备用包需包含后端、前端的完整应用镜像及 MySQL；只导入 Go / Node 基础镜像仍可能需要联网构建。后续沿用备用方案时，启动也保留 `--no-build --pull never`。

`db` 首次使用空数据卷时自动创建 `campus_message_board` 数据库和 `campus_app` 账号；后端连接成功后执行 `schema.sql`，创建缺失的 `messages` 表。后端等待数据库查询就绪，前端等待后端接口就绪，避免只按进程启动顺序判断可用。

### 4. 打开页面

Windows 和 Mac 浏览器均访问 [http://localhost:5173](http://localhost:5173)。可直接打开 [API JSON](http://localhost:5173/api/messages) 检查读取结果。Docker Desktop 管理端口转发，无需查询虚拟机 IP。

发布一条留言，刷新页面，应仍能读到。

### 5. 查看真实数据库

在项目根目录执行，输入 `.env` 中的 `DB_PASSWORD`：

```bash
docker compose exec db mysql -u campus_app -p campus_message_board

```

在 MySQL 提示符执行：

```sql
SHOW TABLES;

SELECT id, nickname, content, created_at
FROM messages ORDER BY id DESC LIMIT 10;

```

输入 `exit` 退出客户端；数据库服务继续运行。数据库数据保存在命名卷 `mysql_data` 中，Git 不复制已有留言。

## 项目结构与请求链

```text
Full-stack-development-class/
├── compose.yaml                         三个容器、网络、挂载与启动依赖
├── .env.example                         本机密码配置模板
├── frontend/
│   ├── Dockerfile.dev                   前端开发镜像与 npm 依赖
│   ├── package.json / package-lock.json  前端命令与依赖
│   ├── vite.config.js                   开发监听地址及 /api 转发
│   └── src/
│       ├── views/MessageBoard.vue       页面状态、发布后刷新
│       ├── components/MessageForm.vue   输入与提交事件
│       ├── components/MessageList.vue   留言列表与加载、错误状态
│       ├── components/portrait/         人物布局与 GSAP 动效
│       └── api/                         HTTP 请求与接口函数
├── backend/
│   ├── Dockerfile.dev                   后端开发镜像与 Go 依赖
│   ├── main.go                         配置、连接模块、注册路由
│   ├── go.mod / go.sum                  Go 模块与下载校验
│   ├── schema.sql                      创建缺失的数据表
│   └── internal/
│       ├── config/config.go            进程环境变量与连接参数
│       ├── handler/message_handler.go  HTTP 解析与响应
│       ├── service/message_service.go  输入规则与业务过程
│       ├── repository/message_repository.go  参数化 SQL
│       └── model/message.go            数据与响应结构
├── scripts/verify-api.mjs               真实接口检查
└── docs/                               环境准备、开发、镜像备用方案与调用链
```

```text
宿主机浏览器
  → localhost:5173 → Docker Desktop 转发 → frontend 容器中的 Vite
  → HTTP /api 转发到 backend:8080
  → Handler → Service → Repository（同一个 Go 程序内部调用）
  → MySQL 协议访问 db:3306 → 数据卷
```

浏览器中的 Vue 页面与组件负责交互；Node.js 容器提供前端开发服务。`backend`、`db` 是容器网络中的服务名，分别用于转发与数据库连接。每个容器的 `127.0.0.1` 指向自身；前端容器不能用它访问后端，后端也不能用它访问 MySQL。

Handler、Service、Repository 是代码分层，共同运行在一个后端容器内。Docker 管理运行环境，分层组织程序职责；两者不是同一个概念。

## 日常开发

```bash
# 查看状态与日志；日志跟随后按 Ctrl+C 只结束查看。
docker compose ps

docker compose logs -f --tail=50 backend

# Go 源码修改后重新启动后端；Vue 修改由 Vite 热更新。
docker compose restart backend

# 停止并移除本项目容器与网络，保留命名卷和留言。
docker compose down

# 下次启动。
docker compose up -d --wait --wait-timeout 600

```

源码通过 bind mount 映射到容器，`node_modules` 与 Go 缓存分别放在命名卷中。新增或更新前端依赖后执行 `docker compose exec frontend npm ci`，再执行 `docker compose restart frontend`；更改 Dockerfile 后重新 `up --build`。若使用本地镜像备用方案，镜像与依赖升级需使用匹配批次的文件并核对已有依赖卷；导入新镜像不会自动替换非空依赖卷。

MySQL 的初始库、账号和密码只在空数据卷首次初始化时应用。已有数据卷时修改 `.env` 不会更改数据库内的账号密码；使用原配置启动，或由了解现有配置的人执行账号变更。常规停机使用 `down`，不加 `--volumes`，不通过删除数据卷解决密码问题。已有原生 MySQL 数据不会自动迁入新卷。

当前 Compose 用于本地开发，运行 Vite 与 `go run`；线上部署需要另外准备构建产物、入口与配置。

## API

| 方法 | 路径 | 成功响应 |
| --- | --- | --- |
| GET | `/api/messages` | `200`，留言数组；无记录为 `[]` |
| POST | `/api/messages` | `201`，新建记录 |

POST 使用 `Content-Type: application/json`：

```json
{"nickname":"  小林  ","content":"周五有人一起去图书馆吗？"}
```

响应示例：

```json
{"id":1,"nickname":"小林","content":"周五有人一起去图书馆吗？","createdAt":"2026-10-21T19:30:00"}
```

`id`、`createdAt` 来自数据库；SQL 字段 `created_at` 对应 JSON 的 `createdAt`。昵称去首尾空白后为 1～30 个 Unicode code point；内容不能为纯空白，原文最多 200 个 code point，保留空白与换行。

错误统一为 `{"message":"错误说明"}`。非法输入或 JSON 返回 `400`，非 JSON 请求返回 `415`，数据库不可连接返回 `503`，其他服务错误返回 `500`。发布失败时保留输入；发布成功后重新 GET 列表。

## 构建与检查

项目已启动时，在项目根目录的 WSL Ubuntu / Mac Terminal 执行：

```bash
docker compose exec backend go mod verify

docker compose exec backend go test ./...

docker compose exec backend go build -o /tmp/campus-message-board .

docker compose exec frontend npm run build

docker compose exec frontend node /workspace/scripts/verify-api.mjs

```

Go 单元测试不连接数据库。接口脚本请求前端容器的 5173，覆盖开发转发，每次新增并保留 3 条检查记录，不删除已有数据。脚本不会使用模拟数据库，运行前确认连接的是本次项目库。

验收时对照页面、POST 响应和 MySQL 查询中的同一条记录，记下 id、昵称、内容、时间与总数；普通 `down` 后重新启动，再查询并刷新页面核对。完整步骤见[页面、SQL 与持久化验收](docs/DEVELOPMENT_DOCKER.md#5-页面sql-与持久化验收)。这些是验收方法，不能仅凭配置文件或健康检查替代实际验收。

## 常见问题与更多文档

- [Docker Desktop 共用指南：配置、启动、排错](docs/DEVELOPMENT_DOCKER.md)
- [镜像拉取超时：本地镜像备用方案](docs/IMAGE_DISTRIBUTION.md)
- [Windows / WSL 环境准备](docs/PREPARATION.md)
- [Mac / Docker Desktop 环境准备](docs/PREPARATION_MACOS.md)
- [工程分层与完整调用链](docs/TEACHING.md)

本机密码、依赖目录、构建产物与日志由 `.gitignore` 排除，MySQL 文件由 Docker 卷管理。
