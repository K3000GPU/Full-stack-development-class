<script setup>
import { onMounted, ref } from 'vue'
import MessageForm from '../components/MessageForm.vue'
import MessageList from '../components/MessageList.vue'
import PortraitLayout from '../components/portrait/PortraitLayout.vue'
import { createMessage, listMessages } from '../api/messages'

const messages = ref([])
const loading = ref(true)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const messageForm = ref(null)
const listError = ref('')
let listRequestId = 0

// 页面层：管理列表、请求状态，协调子组件与接口层。
async function loadMessages() {
  const requestId = ++listRequestId
  loading.value = true
  listError.value = ''
  try {
    const result = await listMessages()
    // 只让最新一次读取更新页面，避免较慢的旧响应覆盖发布后的列表。
    if (requestId === listRequestId) messages.value = result
  } catch (error) {
    if (requestId !== listRequestId) return
    listError.value = `读取留言失败：${error.message}`
    throw error
  } finally {
    if (requestId === listRequestId) loading.value = false
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
  successMessage.value = '留言已发布。'

  try {
    await loadMessages()
  } catch (error) {
    errorMessage.value = `留言已发布，但列表刷新失败：${error.message} 请刷新页面查看，请勿重复发布。`
  } finally {
    submitting.value = false
  }
}

async function refreshMessages() {
  try {
    await loadMessages()
  } catch {
    // 读取失败由列表显示，保留表单输入，允许用户重新加载。
  }
}
onMounted(refreshMessages)
</script>

<template>
  <PortraitLayout>
    <section id="publish" class="pm-panel compose-panel" aria-labelledby="form-heading" data-motion-layer>
      <div class="section-heading">
        <h2 id="form-heading">发布留言</h2>
        <span class="section-number" aria-hidden="true">01</span>
      </div>
      <MessageForm ref="messageForm" :submitting="submitting" @submit="publishMessage" />
      <p v-if="successMessage" class="notice notice-success" role="status">{{ successMessage }}</p>
      <p v-if="errorMessage" class="notice notice-error" role="alert">{{ errorMessage }}</p>
    </section>

    <MessageList :messages="messages" :loading="loading" :error="listError" @retry="refreshMessages" />
  </PortraitLayout>
</template>
