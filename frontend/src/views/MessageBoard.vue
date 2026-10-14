<script setup>
import { onMounted, ref } from 'vue'
import MessageForm from '../components/MessageForm.vue'
import MessageList from '../components/MessageList.vue'
import { createMessage, listMessages } from '../api/messages'

const messages = ref([])
const loading = ref(true)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const messageForm = ref(null)

// 页面层：管理列表、请求状态，协调子组件与接口层。
async function loadMessages() {
  loading.value = true
  try {
    messages.value = await listMessages()
  } finally {
    loading.value = false
  }
}

async function publishMessage(payload) {
  if (submitting.value) return
  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''

  let createdMessage
  try {
    createdMessage = await createMessage(payload)
  } catch (error) {
    errorMessage.value = error.message
    submitting.value = false
    return
  }

  // POST 已成功：立即确认并清空内容，GET 失败也不会让用户误以为未发布。
  messageForm.value.clearContent()
  messages.value = [createdMessage, ...messages.value]
  successMessage.value = '留言已发布，谢谢你的分享。'

  try {
    await loadMessages()
  } catch (error) {
    errorMessage.value = `留言已发布，但列表刷新失败：${error.message} 请刷新页面查看，请勿重复发布。`
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    await loadMessages()
  } catch (error) {
    errorMessage.value = `读取留言失败：${error.message}`
  }
})
</script>

<template>
  <div class="board-shell">
    <header class="board-header">
      <a class="brand" href="/" aria-label="校园留言板首页">
        <span class="brand-mark" aria-hidden="true">言</span>
        <span>校园留言板</span>
      </a>
      <span class="header-note">你好，新同学</span>
    </header>

    <main>
      <section class="intro" aria-labelledby="page-heading">
        <p class="eyebrow">分享 · 相遇 · 校园生活</p>
        <h1 id="page-heading">让想法在这里相遇<span>。</span></h1>
        <p class="intro-description">一声问候、一点灵感，或是今天的小发现。<br class="mobile-break" />写下来，和同学们分享。</p>
      </section>

      <div class="board-layout">
        <section class="compose-panel" aria-labelledby="form-heading">
          <div class="section-heading">
            <div>
              <p class="eyebrow">留下一点你的声音</p>
              <h2 id="form-heading">写一条留言</h2>
            </div>
            <span class="compose-mark" aria-hidden="true">✎</span>
          </div>
          <MessageForm ref="messageForm" :submitting="submitting" @submit="publishMessage" />
          <p v-if="successMessage" class="notice notice-success" role="status">{{ successMessage }}</p>
          <p v-if="errorMessage" class="notice notice-error" role="alert">{{ errorMessage }}</p>
        </section>

        <MessageList :messages="messages" :loading="loading" />
      </div>
    </main>

    <footer>从一句「你好」开始，认识彼此。</footer>
  </div>
</template>
