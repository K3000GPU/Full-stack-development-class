# 验证记录

验证日期：2026-10-21。当前默认实现为 Vue 3 + Go/Gin + MySQL；环境版本见 [README](../README.md)。

## 已通过

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

`legacy/spring-boot` 保留原后端与 Maven Wrapper；本轮没有重新验证该归档实现。当前验收针对 Go/Gin 默认启动链路。

本机账号密码、数据库文件和原始日志在被 Git 忽略的 `.local` 下，不随 clone 分发。未验证其他操作系统、部署环境或远程 MySQL。
