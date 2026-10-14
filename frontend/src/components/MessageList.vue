<script setup>
defineProps({
  messages: { type: Array, required: true },
  loading: { type: Boolean, default: false },
})

// 展示层：只接收留言并渲染，不发送请求、不修改数据。
function formatCreatedAt(value) {
  // 后端返回本地时间的 ISO 字符串，保留其时间，避免浏览器时区转换。
  return value.replace('T', ' ').slice(0, 19)
}
</script>

<template>
  <section class="message-list" aria-labelledby="list-heading" :aria-busy="loading">
    <div class="section-heading">
      <div>
        <p class="eyebrow">来自校园的声音</p>
        <h2 id="list-heading">最新留言</h2>
      </div>
      <span class="message-count">{{ messages.length }} 条</span>
    </div>

    <p v-if="loading" class="list-state" role="status">正在读取留言…</p>
    <p v-else-if="messages.length === 0" class="list-state empty-state">
      <span class="empty-mark" aria-hidden="true">“</span>
      这里还没有留言<br />
      <span>来留下校园里的第一声问候吧。</span>
    </p>

    <ul v-if="messages.length" class="messages">
      <li v-for="message in messages" :key="message.id" class="message-card">
        <div class="message-meta">
          <span class="avatar" aria-hidden="true">{{ Array.from(message.nickname)[0] }}</span>
          <div>
            <h3>{{ message.nickname }}</h3>
            <time :datetime="message.createdAt">{{ formatCreatedAt(message.createdAt) }}</time>
          </div>
        </div>
        <p class="message-content">{{ message.content }}</p>
      </li>
    </ul>
  </section>
</template>
