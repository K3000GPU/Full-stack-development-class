# 教学导读：跟踪一次“发布留言”

面向大学新生，建议用 45 分钟完成。先按 [README](../README.md) 启动 MySQL、Spring Boot 和 Vite，再讲工程分层。课堂目标是能找到文件、说清职责，并用页面、Network 和数据库结果解释一次完整请求。

## 1. 先认目录

| 位置 | 负责什么 |
| --- | --- |
| `frontend/src/views/MessageBoard.vue` | 页面状态、发布过程、发布后重新加载列表 |
| `frontend/src/components/MessageForm.vue` | 输入、基础提示、发出提交事件 |
| `frontend/src/components/MessageList.vue` | 展示留言、加载与错误状态，发出重试事件 |
| `frontend/src/components/portrait/PortraitLayout.vue` | 人物背景、锚点导航与 GSAP 入场；slot 保持表单和列表实例 |
| `frontend/src/api/messages.js` | 将“查看留言”“发布留言”映射到接口 |
| `frontend/src/api/http.js` | 统一发送 HTTP 请求并处理失败响应 |
| `frontend/vite.config.js` | Vite 开发端口与 `/api` 代理 |
| `backend/src/main/java/` | Java 入口、Controller、Service、Repository、Model 和 DTO |
| `backend/src/main/resources/application.yaml` | 后端端口、数据库连接和表初始化配置 |
| `backend/src/main/resources/schema.sql` | 在已指定数据库中建表，不清空已有数据 |
| `frontend/package.json` / `package-lock.json` | 前端脚本、依赖与固定解析结果 |
| `backend/pom.xml` / `mvnw.cmd` / `.mvn/wrapper/` | 后端依赖、构建配置与 Maven Wrapper |

前端组件负责界面，Controller 负责 HTTP，Service 负责业务规则，Repository 负责 SQL。分层的作用是让每个文件的职责清楚，方便从一个功能找到相关代码。

## 2. 按一次发布的顺序找文件与方法

下表中的链接均指向本工程真实源码；行号从 1 开始。方法名称可在编辑器中直接搜索，代码修改后应重新核对行号。

| 顺序 | 文件与定位方法 | 发生的事情 | 调用类型 |
| --- | --- | --- | --- |
| 1 | [MessageForm.vue](../frontend/src/components/MessageForm.vue)：第 16 行 `handleSubmit()` | 表单 `@submit.prevent` 触发方法，检查昵称与内容，`emit('submit', payload)` 交给父页面 | Vue 组件事件，浏览器内部 |
| 2 | [MessageBoard.vue](../frontend/src/views/MessageBoard.vue)：第 30 行 `publishMessage(payload)` | 接收 `@submit`，设置提交状态，调用 `createMessage(payload)` | JavaScript 内部方法调用 |
| 3 | [messages.js](../frontend/src/api/messages.js)：第 8 行 `createMessage(message)` | 指定 `/api/messages`、`POST` 和 JSON 请求体，调用 `request()` | JavaScript 内部方法调用 |
| 4 | [http.js](../frontend/src/api/http.js)：第 2 行 `request(path, options)` | 第 6 行调用 `fetch()`，设置 JSON 请求头；等待并解析响应 | **HTTP POST**，浏览器 → Vite `5173` |
| 5 | [vite.config.js](../frontend/vite.config.js)：第 10 行 `server.proxy` | 将请求转发到 Spring Boot `8080`，路径仍为 `/api/messages` | **HTTP 转发**，Vite → Spring Boot |
| 6 | [MessageController.java](../backend/src/main/java/com/example/campus/controller/MessageController.java)：第 33 行 `createMessage(request)` | `@PostMapping` 接收请求；JSON 被框架转换成 `CreateMessageRequest`，然后调用 Service | HTTP 入口，之后是 Java 内部方法调用 |
| 7 | [MessageService.java](../backend/src/main/java/com/example/campus/service/MessageService.java)：第 30 行 `createMessage(nickname, content)` | 第 53 行 `trimWhitespace()` 去除昵称首尾空白；第 48 行 `characterCount()` 按 code point 计数；检查规则，再调用 `repository.insert()` | Java 内部方法调用 |
| 8 | [MessageRepository.java](../backend/src/main/java/com/example/campus/repository/MessageRepository.java)：第 37 行 `insert(nickname, content)` | 第 39 行 `jdbcTemplate.update()` 执行带两个 `?` 参数的 INSERT，取得生成的主键 | **JDBC / MySQL 通信**，不是 HTTP |
| 9 | 同一文件第 54 行 `queryForObject()` 与第 17 行 `ROW_MAPPER` | 使用 `WHERE id = ?` 回读新记录，把数据库行转换为 `Message`，包括真实 `created_at` | JDBC 查询，然后返回 Java 对象 |
| 10 | [MessageResponse.java](../backend/src/main/java/com/example/campus/dto/MessageResponse.java)：第 9 行 `from(message)`；`MessageController.createMessage()` | Repository → Service → Controller 原路返回；转换为响应 DTO，框架序列化为 JSON，状态为 `201` | Java 内部返回，随后 **HTTP 响应** |
| 11 | [MessageBoard.vue](../frontend/src/views/MessageBoard.vue)：第 46 行；[MessageForm.vue](../frontend/src/components/MessageForm.vue)：第 33 行 `clearContent()` | POST 成功后清空内容、保留昵称，将新记录显示在列表，随后调用 `loadMessages()` | JavaScript 内部方法调用 |
| 12 | [MessageBoard.vue](../frontend/src/views/MessageBoard.vue)：第 17 行 `loadMessages()` → [messages.js](../frontend/src/api/messages.js)：第 4 行 `listMessages()` → [http.js](../frontend/src/api/http.js)：第 2 行 `request()` | 重新请求数据库中的完整留言列表，同样经过 Vite 代理 | **新的 HTTP GET** `/api/messages` |
| 13 | [MessageController.java](../backend/src/main/java/com/example/campus/controller/MessageController.java)：第 27 行 `listMessages()` → [MessageService.java](../backend/src/main/java/com/example/campus/service/MessageService.java)：第 24 行 `listMessages()` → [MessageRepository.java](../backend/src/main/java/com/example/campus/repository/MessageRepository.java)：第 29 行 `findAll()` | SELECT 查询全部留言，按 `created_at DESC, id DESC` 排序；返回响应 DTO 数组 | Java 内部调用 + JDBC，随后 HTTP `200` |
| 14 | [MessageBoard.vue](../frontend/src/views/MessageBoard.vue)：第 21 行 → [MessageList.vue](../frontend/src/components/MessageList.vue)：第 33 行 `v-for` | 用 GET 返回结果更新页面状态，列表组件重新展示 | Vue 响应式更新，浏览器内部 |

首次进入页面时，`MessageBoard.vue` 第 66 行的 `onMounted(refreshMessages)` 经第 59 行 `refreshMessages()` 调用 `loadMessages()`，走表中第 12～14 步。POST 保存一条记录后的回读 SELECT 属于该 POST 的内部过程；页面刷新列表则是另一次独立的 GET 请求。

### 出错时走哪里

Service 校验失败会抛出 [InvalidMessageException.java](../backend/src/main/java/com/example/campus/exception/InvalidMessageException.java)，[ApiExceptionHandler.java](../backend/src/main/java/com/example/campus/exception/ApiExceptionHandler.java) 的 `handleInvalidMessage()` 统一返回 HTTP `400` 与 `ErrorResponse(message)`。JSON 格式问题由同文件的 `handleExceptionInternal()` 统一处理。

前端 `http.js` 的 `request()` 遇到非成功响应会抛出带提示的 `Error`，`MessageBoard.publishMessage()` 捕获后展示，保留输入。若 POST 已成功但随后的 GET 失败，页面会明确提示“留言已发布，但列表刷新失败”，避免学生把已经保存的内容重复发布。

页面读取失败时，`MessageList.vue` 发出 `retry` 组件事件，由 `MessageBoard.refreshMessages()` 重新请求 GET。顶部“发布留言”“查看留言”是 `#publish`、`#messages` 锚点跳转，只在当前页面滚动，不单独发送 API 请求。视觉外壳与素材来源见 [UI_REUSE.md](UI_REUSE.md)。

## 3. 分清三种通信

```text
浏览器中的 Vue
  MessageForm → MessageBoard → messages.js → http.js
                       JavaScript 调用和组件事件
                                   │
                  HTTP POST /api/messages（JSON）
                                   ↓
Vite :5173（开发代理，路径仍为 /api/messages）
                                   │
                              HTTP 转发
                                   ↓
Spring Boot :8080
  Controller → Service → Repository
                   Java 内部方法调用
                                   │
                        JDBC + 参数化 SQL
                                   ↓
教学 MySQL / campus_message_board.messages
```

组件发出的 `submit` 是浏览器进程内部的 Vue 事件，不是 HTTP。Controller 调用 Service、Service 调用 Repository 是同一后端进程内的 Java 方法调用，不会多产生几个 HTTP 请求。Repository 经 JDBC 连接 MySQL，这也不是 HTTP。

浏览器通常看到发布的 POST 和随后刷新列表的 GET。Network 面板不会显示 Java 内部每一次方法调用，也不会显示后端发送给 MySQL 的 SQL。

## 4. 两条 SQL 与一张表

打开 Repository，找到 `INSERT INTO messages ... VALUES (?, ?)` 和 `SELECT ... FROM messages`。`?` 是参数占位符，输入作为值交给 JDBC，不拼接成 SQL。昵称中有单引号也应作为普通文字保存。

| 数据库字段 | 含义 | API 字段 |
| --- | --- | --- |
| `id` | 自增主键，区分不同留言 | `id` |
| `nickname` | 校验并去掉首尾空白后的昵称 | `nickname` |
| `content` | 留言原文 | `content` |
| `created_at` | MySQL 生成的发布时间 | `createdAt` |

只读核对时必须连接教师指定实例：

```sql
SELECT DATABASE();
SELECT id, nickname, content, created_at
FROM campus_message_board.messages
ORDER BY id DESC
LIMIT 10;
```

页面上出现一条留言还不足以单独证明数据库写入成功。将 POST 返回的 `id`、GET 结果和 MySQL SELECT 结果对应起来，并在重启后端后再次读取，才能看到这次持久化的结果。

## 5. 课堂安排与理解检查

| 时间 | 演示内容 | 学生应能指出 |
| --- | --- | --- |
| 0～5 分钟 | clone、按 README 启动、发布一条留言 | 前端、后端和数据库分别在哪里运行 |
| 5～10 分钟 | 查看目录、依赖和配置 | 页面、组件、API 与后端三层的文件 |
| 10～18 分钟 | 从发布按钮进入前端代码 | 组件事件、页面状态和请求方法 |
| 18～28 分钟 | 顺着 POST 阅读 Java | Controller、Service、Repository 的职责 |
| 28～35 分钟 | 看表结构、INSERT 和 SELECT | 字段、记录、主键和参数占位符 |
| 35～42 分钟 | Network、响应、数据库互相核对 | 哪一步是 HTTP，哪一步是内部方法调用 |
| 42～45 分钟 | 用自己的话复述流程 | 为什么刷新或重启后还能看见留言 |

可以用以下三个问题收尾：

1. 如果昵称全是空格，哪个前端文件会提示？绕过前端后，哪个 Java 方法仍然会拒绝？
2. 浏览器发往哪个端口？Vite 转发到哪个端口？为什么两个地方的路径都含 `/api`？
3. 修改页面样式、修改昵称长度规则、修改列表排序，分别应先查看哪一层？

课堂不要清空表来“恢复现场”。需要再次演示时新增一条可辨认的留言，保留之前的记录。
