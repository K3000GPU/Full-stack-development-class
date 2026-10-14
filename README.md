# 校园留言板

大学新生第二节基础培训用的最小全栈工程：输入昵称和内容，发布留言，查看真实 MySQL 中保存的留言。课堂重点是认识目录分层，以及一次操作如何从 Vue 页面经过 HTTP、Spring Boot 和 JDBC 到达数据库。

只提供查看和发布功能，不包含登录。课堂导读见 [docs/TEACHING.md](docs/TEACHING.md)。

## 1. 运行环境

| 工具 | 本次使用版本 | 用途 |
| --- | --- | --- |
| Windows | Windows 11，build 26200；PowerShell | 执行以下启动命令 |
| Git | 2.51.2.windows.1 | clone 工程 |
| Node.js / npm | 24.21.0 / 11.19.0 | 安装、运行和构建前端 |
| JDK | Microsoft OpenJDK 21.0.8+9 LTS | 编译、运行后端 |
| Maven | 3.9.11，由 Wrapper 下载 | 不需要另行安装 Maven |
| MySQL | Community Server 8.4.9 | 保存留言 |
| Vue / Vite | 3.5.44 / 8.3.4 | 页面与前端开发服务 |
| Spring Boot | 3.5.16 | 后端 Web 与数据库支持 |

前端依赖版本由 `frontend/package-lock.json` 固定。后端依赖由 `backend/pom.xml` 固定；Maven Wrapper 3.3.4 使用官方 `only-script` 分发方式，`mvnw`、`mvnw.cmd` 和 `.mvn/wrapper/maven-wrapper.properties` 均随项目提供，首次运行会下载 Maven，无需 `maven-wrapper.jar`。

首次运行需要访问 npm 和 Maven 仓库。教师应在课前安装环境并完成依赖下载。JDK 必须包含 `javac`；`JAVA_HOME` 应指向 JDK 根目录，而不是 `bin`。

## 2. clone 并检查环境

在 PowerShell 中，将仓库地址替换为教师实际发放的 Git URL：

```powershell
$repoUrl = Read-Host '粘贴教师提供的真实 Git URL'
git clone $repoUrl campus-message-board
Set-Location .\campus-message-board
git --version
node --version
npm.cmd --version
java -version
javac -version
```

PowerShell 中使用 `npm.cmd` 可避免 `npm.ps1` 被脚本执行策略拦截。后面的三个终端都以克隆后的 `campus-message-board` 为起点。

## 3. 教师课前准备 MySQL

**只连接教师明确指定的教学 MySQL 实例和 `campus_message_board` 数据库。不要对不明实例、其他数据库或已有数据执行初始化。** 本工程不会执行 `DROP`、`TRUNCATE` 或清空数据。

教师先确认 MySQL 的主机、端口和管理账号。下面示例连接本机 `3306`；本次项目验收使用新建的独立本机实例、端口 `3307`，不是机器上未知的 MySQL 服务。学生按教师提供的端口配置。

在确认过的教学实例上，使用 MySQL 客户端登录。`-p` 会提示输入管理密码，不要把密码写进命令或 Git：

```powershell
mysql --protocol=TCP -h 127.0.0.1 -P 3306 -u root -p
```

如果 `mysql` 不在 PATH，使用安装目录下 `mysql.exe` 的完整路径。仅在上述已确认的实例中执行以下 SQL，并将示例密码替换为本机教学专用密码：

```sql
CREATE DATABASE IF NOT EXISTS campus_message_board
  CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;

CREATE USER IF NOT EXISTS 'campus_student'@'127.0.0.1'
  IDENTIFIED BY '请替换为本机教学专用密码';

GRANT SELECT, INSERT, CREATE ON campus_message_board.*
  TO 'campus_student'@'127.0.0.1';
```

这是本机 TCP 连接账号示例；如果教师已经准备了专用账号，直接使用所提供的配置，无需重复创建。`IF NOT EXISTS` 不会重设已有账号的密码。远程课堂数据库的账号来源地址由教师设置，不要自行扩大授权范围。

数据库本身由教师预先创建；`backend/src/main/resources/schema.sql` 只在已配置的数据库中执行 `CREATE TABLE IF NOT EXISTS messages`。后端启动时自动初始化表，重复启动保留原有留言，不需要手工再导入表结构。

## 4. 启动后端（终端 A）

```powershell
Set-Location .\backend

# 默认主机与端口；教师使用其他端口时在这里修改。
$env:DB_HOST = '127.0.0.1'
$env:DB_PORT = '3306'
$env:DB_NAME = 'campus_message_board'
$env:DB_USERNAME = 'campus_student'
$dbSecret = Read-Host '请输入教师提供的教学账号密码' -AsSecureString
$env:DB_PASSWORD = [System.Net.NetworkCredential]::new('', $dbSecret).Password
Remove-Variable dbSecret

.\mvnw.cmd --version
.\mvnw.cmd spring-boot:run
```

看到后端启动成功后，服务地址为 `http://localhost:8080`。保持终端 A 运行。`DB_NAME`、`DB_USERNAME` 和 `DB_PASSWORD` 必须提供；应用不内置真实账号或密码，也不会自动读取 `.env`。这些变量仅对当前 PowerShell 和它启动的进程生效，关闭终端后重新启动需要再次设置。

数据库连接配置在 `backend/src/main/resources/application.yaml`。默认 `DB_HOST=127.0.0.1`、`DB_PORT=3306`；本次隔离验收时设置 `DB_PORT=3307`。不要把数据库管理账号当作应用账号。

## 5. 启动前端（终端 B）

从项目根目录打开第二个 PowerShell：

```powershell
Set-Location .\frontend
npm.cmd ci
npm.cmd run dev
```

浏览器打开 [http://localhost:5173](http://localhost:5173)。Vite 固定使用 `5173`，`/api` 代理到 `http://127.0.0.1:8080`，**不改写 `/api` 路径**。浏览器请求 `/api/messages`，后端收到的也是 `/api/messages`。代理配置见 `frontend/vite.config.js`。

前端和后端都需要持续运行；MySQL 也必须处于运行状态。结束演示时在前后端终端分别按 `Ctrl+C`。停止服务不会删除数据库记录。

## 6. 接口与规则

| 方法与路径 | 请求 | 成功响应 |
| --- | --- | --- |
| `GET /api/messages` | 无请求体 | `200`，留言数组；无记录时为 `[]` |
| `POST /api/messages` | JSON：`nickname`、`content` | `201`，新建的留言对象 |

发布示例：

```json
{"nickname":"  小林  ","content":"第一次完成前端到数据库的完整请求！"}
```

返回记录的形状如下，实际 ID 和时间由数据库生成：

```json
{
  "id": 1,
  "nickname": "小林",
  "content": "第一次完成前端到数据库的完整请求！",
  "createdAt": "2026-10-13T14:30:00"
}
```

昵称去除首尾空白后必须为 1～30 个字符；内容不能为 `null`、空字符串或纯空白，原文最多 200 个字符。内容保留原文，包括换行。长度按 Unicode code point 计算，一个普通 emoji 计为一个字符。前端先提示用户，后端 Service 再检查；绕过前端也不能跳过规则。

无效输入返回 `400` 和 `{"message":"具体错误说明"}`；昵称和内容必须为 JSON 字符串，数字、布尔值不会自动转成文字。非 JSON 请求返回 `415`，数据库连接不可用返回 `503`，其他服务错误返回 `500`。接口不向浏览器暴露数据库密码、SQL 或 Java 异常堆栈。

表结构：`messages(id, nickname, content, created_at)`；数据库下划线字段 `created_at` 转为 JSON 的 `createdAt`。

## 7. 自己验收

### 页面与完整链路

1. 打开浏览器开发者工具的 Network 面板，访问 `http://localhost:5173`，确认 `GET /api/messages` 成功。
2. 输入昵称 `  新生  ` 和一条独特内容，点击发布。确认 `POST /api/messages` 返回 `201`，记录中昵称为 `新生`。
3. 确认随后出现新的 `GET /api/messages`，页面显示刚发布的留言。
4. 刷新浏览器、重启后端后再次查询，确认记录仍存在。
5. 用教学账号连接同一实例，执行下方只读 SQL，核对 ID、昵称和内容；这一步用于确认真正写入的是 MySQL。

```sql
SELECT DATABASE(), VERSION();
SELECT id, nickname, content, created_at
FROM campus_message_board.messages
ORDER BY id DESC
LIMIT 10;
```

### 自动接口检查与编译

前后端都在运行时，从项目根目录执行接口检查：

```powershell
node .\scripts\verify-api.mjs
```

脚本默认访问 `http://127.0.0.1:5173`，覆盖 Vite 代理；将验证正常发布、读取和输入边界，每次在教学库中保留 3 条验收留言。它不会清空数据；运行前应确认后端指向教师指定的教学库。需要单独检查后端时，可先设置 `$env:API_BASE = 'http://127.0.0.1:8080'`；恢复默认用 `Remove-Item Env:API_BASE`。

在终端 C 中编译前后端；后端命令只运行不连接数据库的单元测试，真实数据库验收由上面的运行流程完成：

```powershell
# 从项目根目录开始。
Set-Location .\frontend
npm.cmd ci
npm.cmd run build
Set-Location ..\backend
.\mvnw.cmd clean test package
```

`npm run build` 产物在 `frontend/dist`，后端 JAR 在 `backend/target`。本课堂使用 Vite 开发服务的 `/api` 代理启动；`vite preview` 不代表完整部署，不能单独替代前后端运行步骤。

### 常见问题

| 现象 | 检查方法 |
| --- | --- |
| `java` 或 `javac` 找不到 | 安装 JDK 21，检查 `JAVA_HOME` 和 PATH，重新打开终端 |
| Maven/npm 首次下载失败 | 检查网络与所处学校网络的仓库访问策略，再重新执行安装命令 |
| 缺少 `DB_NAME` 等配置 | 在启动后端的同一个终端设置环境变量 |
| MySQL 拒绝连接或账号认证失败 | 核对教师指定的主机、端口、数据库、账号及来源地址授权 |
| 页面请求 `/api/messages` 失败 | 先确认后端已启动在 8080，再查看 Vite 终端和浏览器 Network |
| 5173 或 8080 已被占用 | 查清占用者；不要随意关闭未知服务，按教师指引处理 |

## 8. Git 发放范围与实测记录

应提交源码、`README.md`、`docs/TEACHING.md`、前端 lockfile 和完整 Maven Wrapper。`.gitignore` 排除依赖目录、构建产物、`.local`、本机配置和日志，不要提交数据库数据目录或密码。教师发布到自己的 Git 仓库后，将真实 clone 地址发给学生。

本机环境安装包和隔离数据库仅用于本次验收，位于被 Git 排除的 `.local`；学生 clone 后依照本文准备自己的教学环境。

2026-10-14 已在上述 Windows 环境实际验证：

| 验收项 | 结果 |
| --- | --- |
| 官方 Maven Wrapper 3.3.4 下载并运行 Maven 3.9.11 | 通过，发行包配置 SHA-256 校验 |
| `npm.cmd ci` 与 `npm.cmd run build` | 通过，使用提交的 lockfile |
| `mvnw.cmd clean test package` | 通过，28 项单元/MVC 测试，0 失败、0 错误 |
| Spring Boot 启动、schema.sql 初始化 | 通过，连接真实 MySQL 8.4.9 的指定教学库 |
| 直连 8080 与经过 5173 代理的 GET | 均返回 `200`，路径保留 `/api/messages` |
| 经过代理的自动接口检查 | 3 条有效留言返回 `201`，14 次无效请求被拒绝，GET 能回读 |
| MySQL 客户端独立 SELECT 与权限检查 | 通过，应用账号仅有教学库的 SELECT、INSERT、CREATE；emoji 昵称 30 字符、内容 200 字符真实入库 |
| 浏览器操作 | 空昵称校验、点击发布、成功后清空内容、重新 GET 列表和刷新后读取均通过；HTML 作为普通文字显示 |

单元/MVC 测试使用 mock 隔离 Repository，只验证规则和协议；上表 MySQL 结果来自真实 MySQL 客户端查询，使用的是 `127.0.0.1:3307/campus_message_board`，不是 mock 或其他数据库。验收留言保留在本机数据库中，不随 Git 发放。

远程 Git 仓库尚未发布：教师可将此目录推送到自己的仓库，发放实际 URL。本工程没有部署平台配置。
