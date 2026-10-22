# 验证记录

最近更新：2026-10-22。当前方案统一 Docker Desktop：Windows 使用 WSL 2 Ubuntu 的 Integration，Mac 使用本机 Terminal。业务实现为 Vue 3 + Go/Gin + MySQL，验证范围按下面的日期区分。

## 2026-10-22：统一 Docker Desktop

本次同步教学方案、README、Windows / Mac 准备指南、教学导读与共用操作指南。Windows 使用 WSL Integration，Mac 直接使用本机工作区；Docker 命令不加 sudo，两种平台的浏览器均访问 localhost:5173。前后端发布端口统一绑定本机回环地址，MySQL 不发布端口。

已完成：使用项目已有 Go YAML 库解析 Compose，检查三个服务、127.0.0.1 端口映射、健康检查、就绪依赖、Dockerfile 与挂载路径；核对文档链接、代码块和 45 分钟课时。未改动业务代码，因此沿用下面已有的本机业务编译和单元测试记录。

当前执行环境没有可用的 docker 命令，也没有 Mac 实机，未执行 docker compose config、镜像构建或容器启动。正式发放前需补齐：

- Windows：启动 Docker Desktop，启用 Ubuntu-24.04 的 WSL Integration，确认终端 Client / Server 与 Desktop Containers 页面连接同一套环境。
- Mac：安装与芯片匹配的 Docker Desktop，确认 CLI 与本机源码挂载正常。
- 两种平台分别完成 config --quiet、镜像构建、三个服务 healthy、本机 localhost 页面发布、SQL 对照、down / up 后保留，以及代码修改可见检查。

旧的独立 Engine / UTM 环境与数据不自动迁移或删除；旧文档链接保留为新指南入口。

## 2026-10-21：容器配置与业务检查

以下检查完成于 2026-10-21：

| 检查 | 结果 |
| --- | --- |
| 后端依赖与编译 | 本机 Go 1.27.1 执行 go mod verify、go test ./...、go build 通过 |
| 监听与连接配置 | 配置测试覆盖默认 127.0.0.1 监听及容器 HOST=0.0.0.0、DB_HOST=db |
| 前端构建 | 本机 Node.js 24.21.0 执行 npm.cmd run build 通过 |
| Vite 配置 | 分别加载默认配置与容器环境变量，验证监听地址及代理目标覆盖生效 |
| Compose 静态检查 | 使用项目已有 Go YAML 库解析；三个服务、Dockerfile 路径、健康检查、服务依赖和挂载路径检查通过；数据库不发布端口，前端不接收数据库凭据 |
| 文档一致性 | 教学方案与五份项目文档的代码块、39 处本地链接检查通过；课时合计 45 分钟 |

容器配置使用 Node 24.22.1、Go 1.27.3、MySQL 8.4.11 官方镜像标签；标签与架构信息参照 [Node 官方镜像](https://hub.docker.com/_/node)、[Go 官方镜像](https://hub.docker.com/_/golang)、[MySQL 官方镜像](https://hub.docker.com/_/mysql)。本轮没有成功执行镜像拉取或构建，不把本机编译结果视为上述镜像中的验证。

下面保留原生环境的既有验收，作为业务行为参考；它们不代表本轮容器与两种平台已经运行通过。

## 原生环境的既有验收

原生环境为 Windows / PowerShell、Go 1.27.1、Node.js 24.21.0、MySQL 8.4.9。

| 检查 | 结果 |
| --- | --- |
| 后端依赖与构建 | `go mod verify`、`go test ./...`、`go build -o bin/campus-message-board.exe .` 通过 |
| 前端构建 | `npm.cmd run build` 通过 |
| 真实数据库 | 本机独立 MySQL Community Server 8.4.9，`127.0.0.1:3307/campus_message_board` |
| 已有数据兼容 | 更换后端前的 9 条记录逐字段保持一致 |
| GET 与开发代理 | 8080 直连和 5173 的 `/api/messages` 均返回同一真实列表；代理保留 `/api` 路径 |
| POST | `scripts/verify-api.mjs` 新增 3 条记录（ID 10～12），返回真实主键和数据库时间；随后 GET 可读到 |
| 输入规则 | 14 次无效请求被拒绝；另验证 30 个 emoji 昵称、200 个 emoji 内容成功保存，MySQL CHAR_LENGTH 分别为 30、200 |
| 参数化 SQL 与页面文本 | 单引号和 SQL 样式文字按普通内容保存，HTML 样式文字在列表中按文本展示 |
| 浏览器实际发布 | 使用原 Vue 页面发布“迁移验收”，写入 ID 13；显示发布成功，内容框清空、昵称保留，列表更新为 13 条 |
| 独立数据库查询 | MySQL 客户端 SELECT 确认 ID 10～13 和总数 13；应用账号权限仅为该库的 SELECT、INSERT、CREATE |
| 重启与重复初始化 | 停止原 Go 后端，从干净 clone 执行 `go run .`；再次执行 IF NOT EXISTS 建表后，13 条记录所有字段与重启前完全一致 |
| 干净 Git clone | 本地 `git clone --no-hardlinks` 后，`go mod download`、`go mod verify`、单元测试和构建通过；`npm.cmd ci` 重新安装后，前端构建与 `npm.cmd run dev` 通过 |
| clone 请求链 | clone 的 Gin 8080、Vite 5173 代理均返回 200 与相同的 13 条记录；浏览器刷新后显示 ID 13；clone 中无跟踪文件改动 |

接口检查新增的 4 条记录保留在上述专用本机库中，没有清空旧记录或修改其他数据库。已有记录中提到旧后端的文字属于历史留言，未改写。

## 范围说明

Go 单元测试验证输入规则、配置和 HTTP 解析，不连接数据库。真实 MySQL 结论来自运行中的 API、浏览器发布和独立 MySQL 客户端查询，没有用模拟数组、localStorage 或其他数据库替代。

clone 验证重新安装和构建了源码依赖，未复制 node_modules、dist 或后端可执行文件。Go 允许使用本机模块/编译缓存；npm 在线安装出现 ECONNRESET 与下载重试，随后使用既有 `.local/npm-cache` 包缓存，以当前终端的 `npm_config_cache` 和 `npm_config_offline=true` 执行 `npm.cmd ci` 成功。没有改 lockfile 或绕过包校验。此次未完成一次完全无缓存的在线 npm 安装。

本轮检查针对 Go/Gin 默认启动链路，未重新验证旧版后端。

本机账号密码、数据库文件和原始日志在被 Git 忽略的 `.local` 下，不随 clone 分发。未验证其他操作系统、部署环境或远程 MySQL。
