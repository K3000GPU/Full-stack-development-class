# 全栈开发入门课｜Windows 课前准备

本次使用 **Docker Desktop + WSL 2 的 Ubuntu 26.04.1 LTS**，本机 WSL 发行版名称为 **`Ubuntu`**。Windows 上的 Docker Desktop 管理容器；Ubuntu 提供源码目录与终端，通过 WSL Integration 连接 Desktop。系统版本与发行版名称是两回事，命令中的 `-d` 使用 `wsl --list --verbose` 显示的名称。

准备顺序：**准备 WSL Ubuntu → 安装 Docker Desktop → 启用 WSL Integration → 克隆项目 → 启动与验收。**

Vue、Go/Gin、MySQL 都在容器中，不在 Ubuntu 里另装 Docker Engine。Mac 请阅读 [Mac 指南](PREPARATION_MACOS.md)。

## 1. 准备 WSL 2 Ubuntu

已有环境先在 PowerShell 检查：

```powershell
wsl --version

wsl --list --verbose

wsl -d Ubuntu -- cat /etc/os-release

```

本机已经确认 `Ubuntu` 使用 WSL 2、系统为 Ubuntu 26.04.1 LTS，直接复用，无需另装发行版。新电脑尚未安装 Ubuntu 时，先查看可安装列表：

```powershell
wsl --list --online

```

从列表选择所需 Ubuntu 版本，按其实际名称执行 `wsl --install -d <发行版名称>`（替换整个占位符），按提示完成安装、重启和 Linux 用户创建。输入密码时不显示字符是正常情况；安装后重新确认系统版本与注册名称，名称不同的电脑需相应调整本文的 `-d Ubuntu`。

在 PowerShell 确认发行版为 WSL 2：

```powershell
wsl --list --verbose

```

只有 `Ubuntu` 的 VERSION 为 1 时才需要转换；本机已是 2，无需执行：

```powershell
wsl --set-version Ubuntu 2

```

需要更新 WSL 时执行 `wsl --update` 并按提示完成。系统与虚拟化要求以 [Docker Desktop Windows 安装说明](https://docs.docker.com/desktop/setup/install/windows-install/) 和 [Microsoft WSL 安装说明](https://learn.microsoft.com/en-us/windows/wsl/install) 为准；保留现有发行版及项目。

## 2. 安装并启动 Docker Desktop

从 [Docker 官方 Windows 安装页](https://docs.docker.com/desktop/setup/install/windows-install/) 下载与电脑架构匹配的安装包，安装后打开 Docker Desktop，等待 Engine 就绪。

本课使用 Linux containers，设置如下，菜单文字以安装版本为准：

1. **Settings → General**：使用 **WSL 2 based engine**；部分版本默认开启，可能不显示切换项。
2. **Settings → Resources → WSL Integration**：启用 **Ubuntu**。本机默认发行版是 `podman-machine-default`，仅启用默认发行版集成不能证明 `Ubuntu` 已集成，须单独核对 `Ubuntu` 开关。
3. 应用设置，重新打开 Ubuntu 终端。若没有 Integration 选项，核对是否处于 Linux containers 模式。

步骤依据 [Docker 官方 WSL Integration 说明](https://docs.docker.com/desktop/features/wsl/)。Ubuntu 中运行的 docker 命令连接 Desktop 管理的 Engine，无需 `sudo docker` 或 `systemctl start docker`。

**若此前按旧指南安装过 WSL 内的独立 Engine：** 官方要求在使用 Desktop 集成前处理该发行版内已安装的 Engine / CLI，避免冲突。先备份需要保留的数据库并确认现有容器用途，再按官方说明迁移；可另建干净的课程发行版处理。不要删除原发行版、数据库目录或数据卷。独立 Engine 的卷不会自动迁到 Desktop，本指南不执行自动卸载或清理。

## 3. 在 Ubuntu 检查集成

在 PowerShell 进入指定发行版：

```powershell
wsl -d Ubuntu

```

下面开始在 **Ubuntu 终端**执行：

```bash
cat /etc/os-release

uname -s

docker version

docker compose version

docker version --format '{{.Server.Os}}/{{.Server.Arch}}'

```

应看到 Ubuntu 26.04.1 LTS、Linux、Docker Client 与 Server、Compose 版本，以及 `linux/amd64` 等 Engine 平台。找不到 docker 时先回 Desktop 检查 Integration，不使用 apt 安装另一套 Docker。

安装 Git 和编辑工具即可：

```bash
sudo apt update

sudo apt install -y git nano

```

这里 sudo 用于 Ubuntu 的软件安装，项目 Docker 命令不加 sudo。

## 4. 准备并运行项目

课程仓库为 [Yhimhas/Full-stack-development-class](https://github.com/Yhimhas/Full-stack-development-class)，仓库根目录即留言板项目。尚未准备 WSL 项目时，在 Ubuntu 逐行执行，每行成功后再执行下一行：

```bash
mkdir -p ~/projects

cd ~/projects

git clone https://github.com/Yhimhas/Full-stack-development-class.git

cd Full-stack-development-class

pwd

ls -a

```

如果已按先前步骤将 Windows 源码复制到了这个目录，或已有 WSL 克隆，不重复 clone、不删除目录，先确认其中的 `compose.yaml`、Git 状态和 remote。当前 Windows 来源 `C:\Users\lin\Desktop\全栈课程培训课\campus-message-board` 及其本机数据库继续保留。clone 不复制这些数据库、密码、缓存或本机未提交改动；本次 Agent 修订的文档与截图尚未 push 到 GitHub，新 clone 的文档可能仍为旧版本，进入本机发行版继续使用 `wsl -d Ubuntu`。后续编辑和 Compose 使用同一份 WSL 源码。

用户已在 Ubuntu 的 `~/projects` 完成此次 GitHub clone，输出确认对象接收与 delta 解析均为 100%。本机后续使用 `/home/yhimhas/projects/Full-stack-development-class`。Compose 内仍固定 `name: campus-message-board`，所以 Desktop 中的项目名与源码目录名不同是正常情况，不修改该项目名或已有数据卷。

继续完整执行 [Docker Desktop 共用指南](DEVELOPMENT_DOCKER.md)，包括：

1. 在 Ubuntu 的 `~/projects` 准备源码；完成上述 clone 或已有复制后跳过共用指南的克隆步骤，直接使用该目录。
2. 复制 `.env.example`，分别设置 MySQL `root` 的管理密码 `MYSQL_ROOT_PASSWORD`、`campus_app` 的项目密码 `DB_PASSWORD`；后端连接及课堂查询使用项目密码，详见[密码配置说明](DEVELOPMENT_DOCKER.md#3-配置本机密码)。
3. `docker compose up --build -d --wait --wait-timeout 600`，构建并等待三个服务 healthy；镜像拉取超时时才使用[本地镜像备用方案](IMAGE_DISTRIBUTION.md)。
4. Windows 浏览器打开 [http://localhost:5173](http://localhost:5173)。
5. 发布、刷新、容器内 SQL 查询与 down / up 后记录保留检查。

Docker Desktop 的 Containers 页面应看到 `campus-message-board` 及三个服务。软件安装、源码获取、镜像下载、首次构建与验收在课前完成；镜像拉取超时时使用课程提供的本地镜像备用包。

## 5. 编辑器与 Agent

源码放在 Ubuntu 的 `~/projects/Full-stack-development-class`，编辑器与 Agent 打开这份 WSL 工作区。集成终端中 `pwd` 应指向此目录，`uname -s` 显示 Linux。

支持 WSL 的编辑器可连接该发行版；例如 [VS Code WSL 官方文档](https://code.visualstudio.com/docs/remote/wsl)。在项目目录运行 `explorer.exe .` 可查看文件。不要给 Agent 再建一个 Windows 项目副本，避免编辑的代码与容器挂载不同。

## 6. 下次上课怎样打开

先启动 Windows 的 Docker Desktop，等待 Engine 就绪，然后在 PowerShell：

```powershell
wsl -d Ubuntu

```

在 Ubuntu：

```bash
cd ~/projects/Full-stack-development-class

docker compose up -d --wait --wait-timeout 600

docker compose ps

```

Windows 浏览器访问 localhost:5173。课后在项目目录执行 `docker compose down`，保留数据卷。

## 课前确认

- [ ] Docker Desktop 已启动，`Ubuntu` 使用 WSL 2、系统为 Ubuntu 26.04.1 LTS，且 Integration 已启用。
- [ ] Ubuntu 内 docker version 显示 Client / Server，Compose 检查通过。
- [ ] 编辑器与 Agent 打开 WSL 中同一份项目。
- [ ] 三个服务 healthy，Windows 浏览器能发布并刷新读取。
- [ ] MySQL 查询到记录，down / up 后仍在。

完整操作与排错见 [Docker Desktop 共用指南](DEVELOPMENT_DOCKER.md)，工程说明见 [README](../README.md)。
