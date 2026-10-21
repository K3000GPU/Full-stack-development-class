# UI 复用说明

前端采用用户指定的 `ui-portrait-music` skill 的 `portrait` 人物留白版。将原音乐组件的布局与视觉样式适配成留言板，保留 Vue 3 + JavaScript；表单、页面状态和 API 仍按原教学分层组织。

## 文件

| 文件 | 职责 |
| --- | --- |
| `frontend/src/components/portrait/PortraitLayout.vue` | 人物背景、页面外壳、锚点导航和 GSAP 入场动效；业务内容由 slot 传入 |
| `frontend/src/components/portrait/portrait-board.css` | 从原 `portrait-music.css` 适配的黑白面板、硬阴影、斜切导航与响应式布局 |
| `frontend/src/assets/portrait-gold.jpeg` | skill 随附的 gold 原图，未裁切或修改源文件 |
| `frontend/src/views/MessageBoard.vue` | 原有发布与读取流程，另管理列表读取错误和重试 |

桌面内容栏为 `clamp(340px, 35vw, 480px)`，背景固定覆盖整个视口；760px 以下切换为顶部独立人物区域和下方功能区。中文正文不倾斜、不拆字。导航只跳转到页面内的表单与列表，不卸载表单或清空输入。

仅引入 GSAP `3.15.0`，并更新现有 npm lockfile。入场动画使用带组件作用域的 `gsap.matchMedia()`，遵守 `prefers-reduced-motion`，卸载时 `revert()`。外壳提供 `duration`、`reduced`、`paused` props；暂停或减少动态时结束当前动画并呈现最终状态。动画不持续循环。

## 来源

布局来源为该 skill 的 `assets/PortraitMusic.vue`、`assets/portrait-music.css` 及其 GSAP 集成规格。这是 P3R 参考风格的网页适配实现，不是 ATLUS 官方 UI SDK。

图片来自该 skill 中用户提供的 `assets/backgrounds/gold.jpeg`，原始文件名 `1789983476206.jpeg`，尺寸 3305 × 2160，SHA-256：

```text
59c8d3237c3627a5c5533b1d21b4538e584f10295d0edba3af82dec911357643
```

图片沿用原素材的适用授权，本工程不宣称图片原创，也不附加新的图片授权。只复用当前页面需要的一张背景，不引入音乐数据或第三方图标。

## 2026-10-21 验收

- `npm ci`、生产构建与依赖锁定检查通过。
- Chromium 内嵌浏览器检查 1440px、390px、320px；320px 下可用宽度与 scrollWidth 均为 305px，无横向溢出。
- 人物脸部可见，手机人物区域与功能区分离；桌面滚动时背景仍从视口 top=0 开始覆盖。
- 空昵称校验、服务中断时的错误状态、恢复服务后重试均通过；失败时不显示“暂无留言”，重新加载保留表单草稿。
- 手机端发布成功后清空内容、保留昵称，并重新读取列表；真实 MySQL SELECT 确认新增 id=9，原有 8 条记录保留。
- 未实测真实触屏、Safari/Firefox、SSR hydration、系统减少动态设置切换和性能基准；减少动态与暂停保护已接入代码。

截图与原始数据库核对结果保存在被 Git 排除的 `.local/portrait-desktop.png`、`.local/portrait-mobile.png`、`.local/portrait-mysql-verification.txt`。
