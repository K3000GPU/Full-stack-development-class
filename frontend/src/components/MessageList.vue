<script setup>
defineProps({
  messages: { type: Array, required: true },
  loading: { type: Boolean, default: false },
  error: { type: String, default: '' },
})
const emit = defineEmits(['retry'])

// 展示层：只接收留言并渲染，不发送请求、不修改数据。
function formatCreatedAt(value) {
  // 后端返回本地时间的 ISO 字符串，保留其时间，避免浏览器时区转换。
  return value.replace('T', ' ').slice(0, 19)
}
</script>

<template>
  <section id="messages" class="pm-panel message-list" aria-labelledby="list-heading" :aria-busy="loading" data-motion-layer>
    <div class="section-heading">
      <h2 id="list-heading">留言列表</h2>
      <span class="message-count">{{ loading ? '读取中' : error ? '暂不可用' : `${messages.length} 条` }}</span>
    </div>

    <div v-if="error" class="list-state list-error" role="alert">
      <p>{{ error }}</p>
      <button class="retry-button" type="button" @click="emit('retry')">重新加载</button>
    </div>
    <p v-else-if="loading" class="list-state" role="status">正在读取留言…</p>
    <p v-else-if="messages.length === 0" class="list-state empty-state">
      暂无留言
    </p>

    <ul v-else class="messages">
      <li v-for="message in messages" :key="message.id" class="message-card">
        <div class="message-meta">
          <span class="message-index" aria-hidden="true">#{{ message.id }}</span>
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
