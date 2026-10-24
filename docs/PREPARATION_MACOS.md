# 全栈开发入门课｜Mac 课前准备

本次使用 **OrbStack**。代码、编辑器与 Agent 直接使用 Mac 本机环境；Vue、Go/Gin、MySQL 运行在 OrbStack 管理的 Linux 容器中。

准备顺序：**安装 OrbStack → 选择 Docker context → 检查 Git 与 Docker → 在 Mac 克隆项目 → 启动容器 → 浏览器与 SQL 验收。**

留言板直接在 Mac 本机按本指南运行。后续运维课另用 UTM + Ubuntu Server 练习 Linux；该环境按最新的《运维开发课前准备 Mac 补充》准备，不将留言板迁入 UTM。已有虚拟机与数据保留。Windows 请阅读 [Windows 指南](PREPARATION.md)。

## 1. 安装 OrbStack

在“关于本机”确认芯片：

| Mac 芯片 | OrbStack 安装包 |
| --- | --- |
| Apple Silicon，M 系列 | Apple Silicon / arm64 |
| Intel | Intel / amd64 |

从 [OrbStack 官网](https://orbstack.dev/download) 下载对应芯片的版本并打开，按提示完成设置，等待 Docker 引擎启动。已有可用的 OrbStack 和 Git 直接使用，无需重装。

OrbStack 自带 Docker CLI 和 Compose，后续仍使用 `docker` 命令。在 Mac Terminal 选择本次使用的引擎：

```bash
docker context use orbstack
```

已有其他引擎的数据库或容器时保留原环境；不同引擎的数据卷不会自动迁移到 OrbStack。

## 2. 在 Mac Terminal 检查

打开 **macOS Terminal**，Bash / zsh 均可：

```bash
uname -s

git --version

docker version

docker compose version

docker version --format '{{.Server.Os}}/{{.Server.Arch}}'

```

uname 显示 **Darwin** 是正常情况。Docker 应显示 Client 与 Server，Compose 显示版本；Engine 为 Linux，架构通常为 Apple Silicon 的 `arm64` 或 Intel 的 `amd64`。Docker 命令不加 sudo。

若 Git 尚不可用，可按系统提示安装 Apple Command Line Tools：

```bash
xcode-select --install

```

完成后重开终端，再执行 `git --version`。无需安装完整 Xcode。安装步骤可参考 [Apple Command Line Tools 说明](https://developer.apple.com/documentation/xcode/installing-the-command-line-tools)。

若找不到 docker，确认 OrbStack 已完成设置并重开 Terminal；若没有 Server 信息，确认引擎已启动，执行 `docker context show` 应显示 `orbstack`。不通过单独安装 Docker CLI 代替完整环境。

## 3. 在 Mac 本机克隆并运行

继续完整执行 [Docker Compose 共用指南](DEVELOPMENT_DOCKER.md)，在 Mac 的 `~/projects` 克隆仓库、填写 .env。两项密码可以相同，但都不能为空；已有数据卷时修改 .env 不会自动修改账号密码。再在项目目录构建启动：

```bash
cd ~/projects/Full-stack-development-class

docker compose up --build -d --wait --wait-timeout 600

docker compose ps

```

三个服务 healthy 后，用 Mac 浏览器打开 [http://localhost:5173](http://localhost:5173)。无需查询虚拟机 IP、设置 Shared Network 或端口转发。

首次下载与构建在课前完成。镜像请求超时时先等待重试，仍失败再按[本地镜像备用方案](IMAGE_DISTRIBUTION.md)处理。当前提供的包仅适用于 `linux/amd64`；Apple Silicon 的 `linux/arm64` 环境继续主流程或取得匹配架构的包，不能把 amd64 包当作 arm64 原生镜像。Mac 不必另外安装 Homebrew、Node.js、Go、MySQL Server；这些应用运行工具由项目容器提供。

## 4. 验收页面和数据库

页面发布一条留言并刷新，再在项目终端执行：

```bash
docker compose exec db mysql --default-character-set=utf8mb4 -u campus_app -p campus_message_board

```

输入 .env 的 DB_PASSWORD，在 mysql 提示符执行：

```sql
SELECT id, nickname, content, created_at
FROM messages ORDER BY id DESC LIMIT 10;

```

记下同一条记录的 id、昵称、内容、时间与总记录数，输入 `exit`。再按[共用指南的持久化验收](DEVELOPMENT_DOCKER.md#5-页面sql-与持久化验收)执行普通 down / up，重新查询同一 id 并刷新页面，对比字段与总数；不删除卷。使用本地镜像备用方案时，重新 up 也保留 `--no-build --pull never`。

## 5. 编辑器、Agent 与下次启动

编辑器和 Agent 直接打开 Mac 本机 `~/projects/Full-stack-development-class`，不再连接 UTM。源码挂载到容器，Vue 修改后热更新，Go 修改后执行 `docker compose restart backend`。

在 Terminal 运行命令是 Darwin 环境；需要 Go / npm 时用 `docker compose exec backend ...` 或 `docker compose exec frontend ...`，不用切换到 Ubuntu。

下次上课先启动 OrbStack，再打开项目终端：

```bash
cd ~/projects/Full-stack-development-class

docker compose up -d --wait --wait-timeout 600

docker compose ps

```

浏览器使用 localhost:5173，课后执行 `docker compose down`。使用本地镜像备用方案时，上面的 up 改用 `docker compose up --no-build --pull never -d --wait --wait-timeout 600`。

## 课前确认

- [ ] 安装包与芯片匹配，OrbStack 已启动，当前 context 为 `orbstack`。
- [ ] Git、Docker Client / Server、Compose 检查通过。
- [ ] 编辑器与 Agent 打开 Mac 本机项目，容器挂载同一份源码。
- [ ] 三个服务 healthy，localhost 页面能发布并刷新读取。
- [ ] SQL 查到相同记录，down / up 后记录保留。

完整步骤与排错见 [Docker Compose 共用指南](DEVELOPMENT_DOCKER.md)，工程说明见 [README](../README.md)。
