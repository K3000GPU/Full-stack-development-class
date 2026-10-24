# 全栈开发入门课｜Mac 课前准备

本次使用 **Docker Desktop for Mac**。代码、编辑器与 Agent 直接使用 Mac 本机环境；Vue、Go/Gin、MySQL 运行在 Desktop 管理的 Linux 容器中。

准备顺序：**安装 Docker Desktop → 检查 Git 与 Docker → 在 Mac 克隆项目 → 启动容器 → 浏览器与 SQL 验收。**

本课无需准备 UTM、Ubuntu 虚拟机或 SSH 工作区。已有虚拟机和数据可以保留。Windows 请阅读 [Windows 指南](PREPARATION.md)。

## 1. 下载与芯片匹配的 Docker Desktop

在“关于本机”确认芯片：

| Mac 芯片 | Docker Desktop 安装包 |
| --- | --- |
| Apple Silicon，M 系列 | Mac with Apple silicon |
| Intel | Mac with Intel chip |

从 [Docker 官方 Mac 安装页](https://docs.docker.com/desktop/setup/install/mac-install/) 下载；系统版本与内存要求以该页为准。项目使用多架构镜像，不需要为课程强制启用 amd64 模拟。

打开安装包，把 Docker 拖入 Applications，启动 Docker.app，按安装版本提示完成设置，等待 Engine 运行。Docker Desktop 管理容器所需的 Linux 环境。

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

若找不到 docker，先重开 Terminal，并检查 Docker Desktop 中的 CLI 安装设置。官方安装说明列出了 CLI 路径；如果确实安装在用户目录，可以在当前终端临时添加：

```bash
export PATH="$HOME/.docker/bin:$PATH"

docker version

```

还未安装 Docker Desktop 时，不通过单独安装 Docker CLI 代替完整环境。

## 3. 在 Mac 本机克隆并运行

继续完整执行 [Docker Desktop 共用指南](DEVELOPMENT_DOCKER.md)，在 Mac 的 `~/projects` 克隆仓库、填写 .env，再在项目目录构建启动：

```bash
cd ~/projects/Full-stack-development-class

docker compose up --build -d --wait --wait-timeout 600

docker compose ps

```

三个服务 healthy 后，用 Mac 浏览器打开 [http://localhost:5173](http://localhost:5173)。无需查询虚拟机 IP、设置 Shared Network 或端口转发。

首次下载与构建在课前完成。镜像拉取超时时，才按[本地镜像备用方案](IMAGE_DISTRIBUTION.md)导入课程提供的对应架构镜像包：Apple Silicon 使用 arm64 包，Intel 使用 amd64 包。Mac 不必另外安装 Homebrew、Node.js、Go、MySQL Server；这些应用运行工具由项目容器提供。

## 4. 验收页面和数据库

页面发布一条留言并刷新，再在项目终端执行：

```bash
docker compose exec db mysql -u campus_app -p campus_message_board

```

输入 .env 的 DB_PASSWORD，在 mysql 提示符执行：

```sql
SELECT id, nickname, content, created_at
FROM messages ORDER BY id DESC LIMIT 10;

```

确认记录后输入 `exit`。再执行 `docker compose down` 与 `docker compose up -d --wait --wait-timeout 600`，页面应仍能读取该记录；不删除卷。

## 5. 编辑器、Agent 与下次启动

编辑器和 Agent 直接打开 Mac 本机 `~/projects/Full-stack-development-class`，不再连接 UTM。源码挂载到容器，Vue 修改后热更新，Go 修改后执行 `docker compose restart backend`。

在 Terminal 运行命令是 Darwin 环境；需要 Go / npm 时用 `docker compose exec backend ...` 或 `docker compose exec frontend ...`，不用切换到 Ubuntu。

下次上课先启动 Docker Desktop，再打开项目终端：

```bash
cd ~/projects/Full-stack-development-class

docker compose up -d --wait --wait-timeout 600

docker compose ps

```

浏览器使用 localhost:5173，课后执行 `docker compose down`。

## 课前确认

- [ ] 安装包与芯片匹配，Docker Desktop 已启动。
- [ ] Git、Docker Client / Server、Compose 检查通过。
- [ ] 编辑器与 Agent 打开 Mac 本机项目，容器挂载同一份源码。
- [ ] 三个服务 healthy，localhost 页面能发布并刷新读取。
- [ ] SQL 查到相同记录，down / up 后记录保留。

完整步骤与排错见 [Docker Desktop 共用指南](DEVELOPMENT_DOCKER.md)，工程说明见 [README](../README.md)。
