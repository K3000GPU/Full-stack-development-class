# 开发环境文档入口已更新

项目现已统一使用 Docker Desktop，完整操作请阅读 [Docker Desktop 开发环境共用指南](DEVELOPMENT_DOCKER.md)。本文件保留为旧链接的入口。

- Windows：Docker Desktop + WSL 2 Ubuntu，通过 WSL Integration 使用 Docker，见 [Windows 准备指南](PREPARATION.md)。
- Mac：直接在 macOS 安装 Docker Desktop，使用 Mac Terminal 与本机源码，见 [Mac 准备指南](PREPARATION_MACOS.md)。

两种平台均不再按旧方案在 Ubuntu 中单独安装 Docker Engine，浏览器统一访问 localhost:5173。此前已有的虚拟机、独立 Engine 数据与数据库卷保留，切换不会自动迁移它们。
