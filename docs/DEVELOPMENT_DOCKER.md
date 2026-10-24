# Docker Desktop 开发环境共用指南

Windows 先按 [Windows 指南](PREPARATION.md) 安装 Docker Desktop、WSL 2 Ubuntu 并启用 WSL Integration；Mac 先按 [Mac 指南](PREPARATION_MACOS.md) 安装 Docker Desktop。之后使用本文的同一套项目命令。

| 平台 | 执行命令的位置 | 源码位置 |
| --- | --- | --- |
| Windows | WSL 2 的 Ubuntu 26.04.1 LTS 终端（发行版名称 `Ubuntu`） | Ubuntu 的 `~/projects/Full-stack-development-class` |
| Mac | macOS Terminal，Bash / zsh 均可 | Mac 的 `~/projects/Full-stack-development-class` |

Docker 命令统一写成 **`docker ...`，不加 sudo**。Docker Desktop 提供 Engine、CLI 和 Compose；Node.js、Go、MySQL 随项目容器提供。两种平台的浏览器都使用 `http://localhost:5173`。

本项目的三个服务为 `frontend`（前端）、`backend`（后端）和 `db`（数据库）。`db` 是 **database** 的缩写，也是 Compose 中 MySQL 的服务名；后端配置 `DB_HOST: db`，表示通过该服务名连接 MySQL。连接后使用的数据库名为 `campus_message_board`。

## 1. 检查 Docker Desktop 已启动

先打开宿主机上的 Docker Desktop，等待 Engine 运行，再在上述终端执行：

```bash
git --version

docker version

docker compose version

docker version --format '{{.Server.Os}}/{{.Server.Arch}}'

```

应看到 Git 版本、Docker 的 Client 与 Server 信息、Compose 版本以及 Engine 的 Linux 平台和架构。默认按下面的主流程构建启动；镜像拉取超时时使用本地镜像备用方案。

Docker Desktop 的 Containers 页面可查看容器状态和日志；终端与界面应能看到同一个项目。

## 2. 获取代码

```bash
mkdir -p ~/projects

cd ~/projects

git clone https://github.com/Yhimhas/Full-stack-development-class.git

cd Full-stack-development-class

ls

```

应看到 `README.md`、`compose.yaml`、`frontend`、`backend`、`docs`。这就是项目根目录；GitHub 仓库名为 `Full-stack-development-class`，本地目录也使用 `Full-stack-development-class`。私有仓库先完成对应终端中的 Git 认证。

已有克隆时进入已有目录；同步课程版本前先保存本机改动。Windows 源码放在 Ubuntu 的文件系统中，Mac 直接放在自己的用户目录；编辑器与 Agent 打开容器实际挂载的这一份源码。

## 3. 配置本机密码

> [!WARNING]
> **`.env` 是隐私配置文件，包含数据库密码。**
> 不要将 `.env` 提交到 Git 或推送到 GitHub，也不要把文件、内容或含密码的截图发给 Agent 或他人。
> `.gitignore` 已配置忽略 `.env`，请保留该忽略规则。

```bash
cp -n .env.example .env

chmod 600 .env

nano .env

```

`cp -n` 表示目标 `.env` 已存在时跳过复制，保留已有配置和密码；不存在时才从 `.env.example` 创建。

为两个变量分别设置不同的密码：

```dotenv
MYSQL_ROOT_PASSWORD='自行设置的管理密码'
DB_PASSWORD='自行设置的项目密码'
```

这两个密码对应 MySQL 中的两个账号，由你分别设置：

| 配置项 | 对应账号 | 用途 | 传给哪些容器 |
| --- | --- | --- | --- |
| `MYSQL_ROOT_PASSWORD`（管理密码） | MySQL 的 `root` 管理员 | 管理数据库、账号和权限等；首次空卷初始化时设置管理员密码 | 仅 `db` |
| `DB_PASSWORD`（项目密码） | `campus_app` 应用账号 | Go 后端连接 `campus_message_board` 库，执行建表、查询和新增留言；课堂手动查询也用这个账号 | `db`、`backend` |

Compose 把 `.env` 的 `DB_PASSWORD` 作为 `MYSQL_PASSWORD` 传给 `db`，用于首次创建 `campus_app`；同时以 `DB_PASSWORD` 传给后端，因此这两处必须使用同一个项目密码。后端日常通过应用账号访问项目库，管理密码留给数据库管理操作；前端和浏览器都不接收这两个密码。

例如，执行 `docker compose exec db mysql -u campus_app -p campus_message_board` 后，应在密码提示符输入 **`DB_PASSWORD`**。`root` 在这里指 MySQL 管理员；Ubuntu 登录和 `sudo` 使用的是 Linux 用户密码，与这两项分别配置。

这两项只在 MySQL 空数据卷首次初始化时设置对应账号的密码。已有数据卷时，仅编辑 `.env` 或重建容器不会修改数据库中已保存的密码；需要变更时，应由管理员在数据库内修改账号密码并同步配置，保留已有数据库和数据卷。

建议至少 16 位字母数字组合。保留外层单引号，密码中不使用单引号或换行。`Ctrl+O`、回车保存，`Ctrl+X` 退出；也可以用编辑器编辑此文件。

`.env` 由 Compose 读取，再将相应参数传给容器。前端不接收数据库密码。该文件由 Git 忽略，保留在自己的电脑。

## 4. 首次构建并启动（主流程）

在项目根目录执行：

```bash
docker compose config --quiet

docker compose up --build -d --wait --wait-timeout 600

docker compose ps

```

`config --quiet` 检查配置，不输出展开后的密码。首次构建需要下载镜像、npm 与 Go 依赖，并编译 Go；请在课前完成。

三个服务应为 `healthy`。`db` 首次使用空数据卷时自动创建数据库与应用账号；后端连接成功后创建缺失的表。前端健康检查也经过 API 转发访问后端，确保这条链路可用。首次 Go 编译和数据库初始化需要时间，在课前完成。

`--wait-timeout 600` 控制服务健康等待，不控制前面镜像下载或构建的总耗时。

### 镜像拉取超时时：使用本地镜像备用方案

若拉取 MySQL 镜像或构建器读取 Go / Node 基础镜像信息时超时，按[备用方案](IMAGE_DISTRIBUTION.md)导入提供的对应架构完整镜像包，并使用其中的 `docker compose up --no-build --pull never -d --wait --wait-timeout 600` 启动。该分支启动成功后，继续下面的页面、SQL 和持久化验收；主流程构建成功时直接继续。

现在打开自己的 Windows 或 Mac 浏览器：

- [留言板页面](http://localhost:5173)。
- [查看接口 JSON](http://localhost:5173/api/messages)，空库返回 `[]`。

Docker Desktop 将已发布的容器端口转发到宿主机，参考 [Docker Desktop 网络说明](https://docs.docker.com/desktop/features/networking/)。本项目把 5173、8080 发布到本机回环地址，MySQL 3306 只在容器网络中使用。

## 5. 页面、SQL 与持久化验收

从页面发布一条自己能识别的留言，再刷新。然后在项目根目录登录数据库容器：

```bash
docker compose exec db mysql -u campus_app -p campus_message_board

```

输入 `.env` 中的 **DB_PASSWORD**。在 `mysql>` 提示符执行：

```sql
SELECT DATABASE(), VERSION();

SHOW TABLES;

SELECT id, nickname, content, created_at
FROM messages ORDER BY id DESC LIMIT 10;

SELECT COUNT(*) AS total_messages FROM messages;

```

应看到刚刚发布的记录。记下该条记录的 id、昵称、内容、创建时间与总记录数，作为重启前基线。输入 `exit` 退出客户端，再在项目终端逐行执行：

```bash
docker compose down

docker volume inspect campus-message-board_mysql_data --format '{{.Name}} {{.CreatedAt}}'

docker compose up -d --wait --wait-timeout 600

docker compose ps

```

确认数据卷仍存在且三个服务重新 healthy，再用相同的 MySQL 登录与 SELECT 命令查询，对照同一 id 的全部字段与总记录数，最后刷新浏览器确认该留言仍可读，完成持久化检查。`down` 保留命名卷；日常停机不加 `--volumes`。Docker Desktop 的删除卷与恢复出厂设置操作也不作为日常停机方式。

已有卷时修改 `.env` 不会更改数据库里的账号密码；切换 Engine / Docker context 也不会自动迁移原有数据。遇到认证问题先核对首次初始化的密码与所连接的 Engine，保留现有数据。

## 6. 开发与检查

```bash
docker compose logs -f --tail=50 backend

```

`Ctrl+C` 只结束日志跟随。Vue 修改后由 Vite 热更新；Go 修改后执行：

```bash
docker compose restart backend

```

编译与测试也在容器内执行：

```bash
docker compose exec backend go test ./...

docker compose exec frontend npm run build

docker compose exec frontend node /workspace/scripts/verify-api.mjs

```

接口脚本每次保留三条新记录，不删除已有留言。主流程开发时，前端依赖变更后执行 `docker compose exec frontend npm ci` 并重启前端；Dockerfile 变更后重新 `up --build`。本地镜像备用分支需要使用与新源码配套的镜像包，并核对依赖卷，详见备用方案。

Windows 的编辑器与 Agent 使用 WSL 工作区，在终端执行 `pwd`、`uname -s`，应为 Ubuntu 中的项目目录及 `Linux`。Mac 的编辑器与 Agent 直接打开本机项目，终端 `uname -s` 显示 `Darwin` 是正常情况；Node / Go 命令通过 `docker compose exec` 执行。

## 7. 常见问题

```bash
docker context show

docker context ls

docker compose ps -a

docker compose logs --tail=80 db backend frontend

```

| 现象 | 检查方式 |
| --- | --- |
| 找不到 docker | Windows 核对 `Ubuntu` 的 WSL Integration；Mac 按安装指南检查 CLI 路径并重开终端 |
| 无 Server 信息 / 无法连接 daemon | 确认 Docker Desktop 已启动并等待 Engine 就绪；核对当前 context 与是否设置过 DOCKER_HOST |
| Docker 权限错误 | 不通过加 sudo 绕过；先核对是否误连了 WSL 中原先的独立 Engine |
| Containers 界面与终端看到的项目不同 | 检查 context、DOCKER_HOST 和 WSL Integration，确保连接同一套 Docker Desktop；不要盲目切换或重置 |
| 提示密码变量未设置 | 确认在项目根目录，检查 `.env` 中两个密码非空 |
| 镜像拉取 / 基础镜像信息读取超时 | 使用[本地镜像备用方案](IMAGE_DISTRIBUTION.md)，导入课程提供的完整应用镜像与 MySQL 镜像后启动 |
| 备用方案提示缺少本地镜像 / No such image | 核对课程镜像包导入结果、`docker compose config --images` 与本机镜像名及标签，使用同一课程版本 |
| 镜像架构不匹配 / exec format error | 核对 Engine 的 `linux/amd64` 或 `linux/arm64`；备用方案使用课程提供的对应架构镜像包 |
| 备用方案提示缺少 npm / Go 依赖 | 核对源码、完整应用镜像和依赖卷是否为同一批次，核对镜像包版本及校验结果，保留已有卷 |
| unhealthy / 首次启动超时 | 查看相应容器日志；确认首次编译或数据库初始化是否仍在进行，再重试等待 |
| MySQL 认证失败 | 核对初始化时密码与当前 Engine，保留卷，不通过删卷重设密码 |
| localhost:5173 打不开 | 检查 Desktop、三服务健康状态和 Ports 映射；Windows / Mac 浏览器均使用 localhost，不使用虚拟机 IP |
| Mac 提示挂载被拒绝 | 项目放在 `~/projects`，检查 Desktop 的 File sharing 及系统目录访问权限；见 [Desktop 设置说明](https://docs.docker.com/desktop/settings-and-maintenance/settings/) |
| 改代码页面没变 | 确认编辑的是实际挂载源码，Windows 优先从 WSL 工作区编辑；检查 Vite 日志，必要时重启前端 |
| Go 编译进程被 Killed | 检查 Desktop / WSL 的可用内存，适当调整资源再重试 |

## 课前完成清单

- [ ] Docker Desktop 已启动，Client、Server 与 Compose 检查通过。
- [ ] Windows WSL Integration 正常；Mac 使用本机 Terminal 和项目目录。
- [ ] 编辑器与 Agent 打开容器实际挂载的同一份源码。
- [ ] 三个服务 healthy，宿主机浏览器可以打开 localhost:5173。
- [ ] 发布、刷新与 SQL 结果对应，down / up 后记录仍在。
- [ ] 主流程构建通过，或拉取超时后通过备用方案导入启动；知道在 Desktop 界面及终端查看状态和日志。

反馈完成情况或错误信息即可，不发送密码与 `.env`。
