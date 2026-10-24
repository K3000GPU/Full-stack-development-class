# 镜像拉取超时：本地镜像备用方案

主流程按[共用指南](DEVELOPMENT_DOCKER.md)执行 `docker compose up --build -d --wait --wait-timeout 600`。仅当镜像拉取超时（包括构建器读取 Go / Node 基础镜像信息超时）时，使用本文件的备用方案。已成功构建启动时，继续页面与数据库验收即可。

## 镜像请求超时的处理顺序

先等待并重试原命令。确认在拉取镜像或读取基础镜像 metadata 阶段出现 `DeadlineExceeded: context deadline exceeded`、`net/http: TLS handshake timeout` 等超时后，可逐个拉取，已成功的镜像跳过：

```bash
docker pull golang:1.27.3-bookworm
docker pull node:24.22.1-bookworm-slim
docker pull mysql:8.4.11
```

拉取成功后回项目目录重试主流程。单独拉取是可选步骤，不能保证后续构建不访问网络；仍超时也可直接按下文导入完整镜像。`npm ci`、`go mod download` 报错或容器未 healthy 时，按对应日志排查，不能只凭 `DeadlineExceeded` 就判定为镜像请求超时。`--wait-timeout 600` 仅控制服务就绪等待时间。

## 1. 确认平台与准备文件

在 Ubuntu / Mac Terminal 执行：

```bash
docker version --format '{{.Server.Os}}/{{.Server.Arch}}'

```

当前课程镜像包为 `campus-message-board-images-linux-amd64.tar`，适用于 Engine 显示 `linux/amd64` 的电脑。`linux/arm64` 环境继续使用主流程，或取得对应架构的镜像包；本文件不作为 arm64 原生镜像使用。

将收到的 `.tar` 文件及配套 `SHA256SUMS` 文件放到 `~/course-images/`，与源码目录分开存放。镜像包包含以下三个完整镜像：

| 镜像 | 用途 |
| --- | --- |
| `campus-message-board-backend:latest` | Go 后端及已准备的模块依赖 |
| `campus-message-board-frontend:latest` | Node.js 前端开发服务及已安装的 npm 依赖 |
| `mysql:8.4.11` | MySQL 服务程序；留言保存在本机独立的数据卷 |

## 2. 校验并导入

已完成 [Windows](PREPARATION.md) 或 [Mac](PREPARATION_MACOS.md) 的安装、源码准备和 `.env` 配置后，执行：

```bash
cd ~/course-images

sha256sum -c SHA256SUMS

```

上面的校验命令用于 Ubuntu；Mac 使用 `shasum -a 256 -c SHA256SUMS`。应显示镜像文件名和 `OK`。校验失败时先检查文件是否完整，校验通过后再继续导入：

```bash
docker load --input ~/course-images/campus-message-board-images-linux-amd64.tar

```

导入后核对镜像及平台：

```bash
docker image inspect campus-message-board-backend:latest campus-message-board-frontend:latest mysql:8.4.11 --format '{{.RepoTags}} {{.Os}}/{{.Architecture}}'

```

应列出上表中的三个镜像，平台均为 `linux/amd64`。导入命令说明见 [Docker load 文档](https://docs.docker.com/reference/cli/docker/image/load/)。

## 3. 使用本机镜像启动

回到项目根目录，确认已按[共用指南](DEVELOPMENT_DOCKER.md)配置自己的 `.env`：

```bash
cd ~/projects/Full-stack-development-class

docker compose config --quiet

docker compose config --images

docker compose up --no-build --pull never -d --wait --wait-timeout 600

docker compose ps

```

`--no-build` 跳过镜像构建，`--pull never` 仅使用本机镜像，详见 [Compose up 文档](https://docs.docker.com/reference/cli/docker/compose/up/)。后续使用本备用分支启动时沿用这两个选项。

## 4. 继续验收

三个服务应为 `healthy`。打开 [http://localhost:5173](http://localhost:5173)，按[共用指南](DEVELOPMENT_DOCKER.md)完成页面发布、刷新、MySQL 查询和普通 down / up 后记录保留检查；重新 up 时使用本页的启动命令。

使用自己的 `.env` 和数据库密码，保留现有数据库及数据卷。依赖或源码版本变化时，使用匹配批次的镜像包；已有非空依赖卷不会因导入镜像自动更新。缺少依赖时先核对版本和卷内容，保留原卷。
