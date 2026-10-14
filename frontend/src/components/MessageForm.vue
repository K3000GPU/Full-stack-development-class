<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  submitting: { type: Boolean, default: false },
})
const emit = defineEmits(['submit'])

const nickname = ref('')
const content = ref('')
const validationError = ref('')
const nicknameLength = computed(() => Array.from(nickname.value.trim()).length)
const contentLength = computed(() => Array.from(content.value).length)

// 表单层：管理输入和即时校验；向父页面提交事件，不直接调用 HTTP。
function handleSubmit() {
  if (props.submitting) return

  validationError.value = ''
  if (nicknameLength.value === 0 || nicknameLength.value > 30) {
    validationError.value = '昵称去除首尾空白后须为 1～30 个字符。'
    return
  }
  if (!content.value.trim() || contentLength.value > 200) {
    validationError.value = '留言内容不能为空，且最多为 200 个字符。'
    return
  }

  emit('submit', { nickname: nickname.value.trim(), content: content.value })
}

// 仅在父页面确认 POST 成功后清空内容，失败时保留输入。
function clearContent() {
  content.value = ''
  validationError.value = ''
}

defineExpose({ clearContent })
</script>

<template>
  <form class="message-form" novalidate @submit.prevent="handleSubmit">
    <fieldset :disabled="submitting">
      <legend class="sr-only">发布留言</legend>
      <div class="field">
        <label for="nickname">你的昵称</label>
        <input
          id="nickname"
          v-model="nickname"
          name="nickname"
          autocomplete="nickname"
          placeholder="让同学们认识你"
          aria-describedby="nickname-hint"
        />
        <span id="nickname-hint" class="field-hint">
          <span>去除首尾空白，最多 30 个字符</span>
          <span :class="{ 'over-limit': nicknameLength > 30 }">{{ nicknameLength }} / 30</span>
        </span>
      </div>

      <div class="field">
        <label for="content">留言内容</label>
        <textarea
          id="content"
          v-model="content"
          name="content"
          rows="6"
          placeholder="分享一个想法，或向新同学打个招呼……"
          aria-describedby="content-hint"
        ></textarea>
        <span id="content-hint" class="field-hint">
          <span>每一条分享，都值得被看见</span>
          <span :class="{ 'over-limit': contentLength > 200 }">{{ contentLength }} / 200</span>
        </span>
      </div>

      <p v-if="validationError" class="notice notice-error" role="alert">{{ validationError }}</p>
      <button class="publish-button" type="submit" :aria-busy="submitting">
        {{ submitting ? '发布中…' : '发布留言' }}
        <span aria-hidden="true">↗</span>
      </button>
    </fieldset>
  </form>
</template>
