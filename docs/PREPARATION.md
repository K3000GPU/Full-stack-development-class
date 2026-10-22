# 全栈开发入门课｜Windows 课前准备

本次使用 **Docker Desktop + WSL 2 的 Ubuntu-24.04**。Windows 上的 Docker Desktop 管理容器；Ubuntu 提供源码目录与终端，通过 WSL Integration 连接 Desktop。

准备顺序：**准备 WSL Ubuntu → 安装 Docker Desktop → 启用 WSL Integration → 克隆项目 → 启动与验收。**

Vue、Go/Gin、MySQL 都在容器中，不在 Ubuntu 里另装 Docker Engine。Mac 请阅读 [Mac 指南](PREPARATION_MACOS.md)。

## 1. 准备 WSL 2 Ubuntu

已有环境先在 PowerShell 检查：

```powershell
wsl --version
wsl --list --verbose
```

尚未安装课程发行版时，在管理员 PowerShell 执行：

```powershell
wsl --install -d Ubuntu-24.04
```

按提示重启 Windows，打开 Ubuntu，创建 Linux 用户名和密码。输入密码时不显示字符是正常情况。

在 PowerShell 确认发行版为 WSL 2：

```powershell
wsl --list --verbose
```

如果 Ubuntu-24.04 的 VERSION 为 1，转换：

```powershell
wsl --set-version Ubuntu-24.04 2
```

需要更新 WSL 时执行 `wsl --update` 并按提示完成。系统与虚拟化要求以 [Docker Desktop Windows 安装说明](https://docs.docker.com/desktop/setup/install/windows-install/) 和 [Microsoft WSL 安装说明](https://learn.microsoft.com/en-us/windows/wsl/install) 为准；保留现有发行版及项目。

## 2. 安装并启动 Docker Desktop

从 [Docker 官方 Windows 安装页](https://docs.docker.com/desktop/setup/install/windows-install/) 下载与电脑架构匹配的安装包，安装后打开 Docker Desktop，等待 Engine 就绪。

本课使用 Linux containers，设置如下，菜单文字以安装版本为准：

1. **Settings → General**：使用 **WSL 2 based engine**；部分版本默认开启，可能不显示切换项。
2. **Settings → Resources → WSL Integration**：启用 **Ubuntu-24.04**。
3. 应用设置，重新打开 Ubuntu 终端。若没有 Integration 选项，核对是否处于 Linux containers 模式。

步骤依据 [Docker 官方 WSL Integration 说明](https://docs.docker.com/desktop/features/wsl/)。Ubuntu 中运行的 docker 命令连接 Desktop 管理的 Engine，无需 `sudo docker` 或 `systemctl start docker`。

**若此前按旧指南安装过 WSL 内的独立 Engine：** 官方要求在使用 Desktop 集成前处理该发行版内已安装的 Engine / CLI，避免冲突。先备份需要保留的数据库并确认现有容器用途，再按官方说明迁移；可另建干净的课程发行版处理。不要删除原发行版、数据库目录或数据卷。独立 Engine 的卷不会自动迁到 Desktop，本指南不执行自动卸载或清理。

## 3. 在 Ubuntu 检查集成

在 PowerShell 进入指定发行版：

```powershell
wsl -d Ubuntu-24.04
```

下面开始在 **Ubuntu 终端**执行：

```bash
cat /etc/os-release
uname -s
docker version
docker compose version
docker run --rm hello-world
```

应看到 Ubuntu 24.04、Linux、Docker Client 与 Server、Compose 版本，以及 Hello from Docker。找不到 docker 时先回 Desktop 检查 Integration，不使用 apt 安装另一套 Docker。

安装 Git 和编辑工具即可：

```bash
sudo apt update
sudo apt install -y git nano
```

这里 sudo 用于 Ubuntu 的软件安装，项目 Docker 命令不加 sudo。

## 4. 准备并运行项目

继续完整执行 [Docker Desktop 共用指南](DEVELOPMENT_DOCKER.md)，包括：

1. 在 Ubuntu 的 `~/projects` 克隆课程仓库。
2. 复制 `.env.example` 并填写两个本机密码。
3. `docker compose up --build -d --wait --wait-timeout 600`，等待三个服务 healthy。
4. Windows 浏览器打开 [http://localhost:5173](http://localhost:5173)。
5. 发布、刷新、容器内 SQL 查询与 down / up 后记录保留检查。

Docker Desktop 的 Containers 页面应看到 `campus-message-board` 及三个服务。安装、镜像下载与首次编译在课前完成，课堂不等待下载。

## 5. 编辑器与 Agent

源码放在 Ubuntu 的 `~/projects/campus-message-board`，编辑器与 Agent 打开这份 WSL 工作区。集成终端中 `pwd` 应指向此目录，`uname -s` 显示 Linux。

支持 WSL 的编辑器可连接该发行版；例如 [VS Code WSL 官方文档](https://code.visualstudio.com/docs/remote/wsl)。在项目目录运行 `explorer.exe .` 可查看文件。不要给 Agent 再建一个 Windows 项目副本，避免编辑的代码与容器挂载不同。

## 6. 下次上课怎样打开

先启动 Windows 的 Docker Desktop，等待 Engine 就绪，然后在 PowerShell：

```powershell
wsl -d Ubuntu-24.04
```

在 Ubuntu：

```bash
cd ~/projects/campus-message-board
docker compose up -d --wait --wait-timeout 600
docker compose ps
```

Windows 浏览器访问 localhost:5173。课后在项目目录执行 `docker compose down`，保留数据卷。

## 课前确认

- [ ] Docker Desktop 已启动，Ubuntu-24.04 使用 WSL 2 且 Integration 已启用。
- [ ] Ubuntu 内 docker version 显示 Client / Server，Compose 与 hello-world 通过。
- [ ] 编辑器与 Agent 打开 WSL 中同一份项目。
- [ ] 三个服务 healthy，Windows 浏览器能发布并刷新读取。
- [ ] MySQL 查询到记录，down / up 后仍在。

完整操作与排错见 [Docker Desktop 共用指南](DEVELOPMENT_DOCKER.md)，工程说明见 [README](../README.md)。
