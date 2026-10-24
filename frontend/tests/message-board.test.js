import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import test from 'node:test'
import vm from 'node:vm'
import { ref } from 'vue'

function deferred() {
  let resolve, reject
  const promise = new Promise((accept, fail) => { resolve = accept; reject = fail })
  return { promise, resolve, reject }
}

function createBoard(reads, createMessage = async () => ({ id: 2 })) {
  // 直接运行页面 script 中的真实方法；只替换 API 和挂载入口，不连接 HTTP 或数据库。
  const source = readFileSync(new URL('../src/views/MessageBoard.vue', import.meta.url), 'utf8')
    .match(/<script setup>([\s\S]*?)<\/script>/)[1].replace(/^import .*$/gm, '')
  let cleared = 0
  const context = {
    ref, onMounted() {}, createMessage,
    listMessages() {
      assert.ok(reads.length, '出现未预期的列表请求')
      return reads.shift().promise
    },
  }
  vm.runInNewContext(`${source}\nthis.board = { loadMessages, publishMessage, messages, messageForm, loading, listError, errorMessage, successMessage, submitting };`, context)
  context.board.messageForm.value = { clearContent() { cleared++ } }
  return { ...context.board, cleared: () => cleared }
}

test('首次 GET 晚于发布后 GET 返回，不覆盖新留言', async () => {
  const initial = deferred(), refreshed = deferred()
  const board = createBoard([initial, refreshed])
  const firstRead = board.loadMessages()
  const publish = board.publishMessage({ nickname: '同学', content: '新留言' })
  await Promise.resolve()
  refreshed.resolve([{ id: 2 }, { id: 1 }])
  await publish
  initial.resolve([{ id: 1 }])
  await firstRead
  assert.deepEqual(Array.from(board.messages.value, message => message.id), [2, 1])
  assert.equal(board.cleared(), 1)
})

test('旧 GET 失败不覆盖新结果，也不显示过期错误', async () => {
  const old = deferred(), current = deferred()
  const board = createBoard([old, current])
  const firstRead = board.loadMessages()
  const latestRead = board.loadMessages()
  current.resolve([{ id: 2 }])
  await latestRead
  old.reject(new Error('旧请求失败'))
  await firstRead
  assert.equal(board.messages.value[0].id, 2)
  assert.equal(board.listError.value, '')
  assert.equal(board.loading.value, false)
})

test('旧 GET 结束时，仍等待最新 GET 的加载状态保持不变', async () => {
  const old = deferred(), current = deferred()
  const board = createBoard([old, current])
  const firstRead = board.loadMessages()
  const latestRead = board.loadMessages()
  old.resolve([{ id: 1 }])
  await firstRead
  assert.equal(board.loading.value, true)
  assert.equal(board.messages.value.length, 0)
  current.resolve([{ id: 2 }])
  await latestRead
  assert.equal(board.loading.value, false)
  assert.equal(board.messages.value[0].id, 2)
})

test('最新 GET 失败仍显示错误并结束加载', async () => {
  const current = deferred()
  const board = createBoard([current])
  const read = board.loadMessages()
  current.reject(new Error('连接失败'))
  await assert.rejects(read, /连接失败/)
  assert.match(board.listError.value, /连接失败/)
  assert.equal(board.loading.value, false)
})

test('POST 成功而列表刷新失败，保留新记录并明确提示已发布', async () => {
  const current = deferred()
  const board = createBoard([current])
  const publish = board.publishMessage({ nickname: '同学', content: '新留言' })
  await Promise.resolve()
  current.reject(new Error('读取失败'))
  await publish
  assert.equal(board.messages.value[0].id, 2)
  assert.equal(board.cleared(), 1)
  assert.equal(board.successMessage.value, '留言已发布。')
  assert.match(board.errorMessage.value, /留言已发布，但列表刷新失败/)
  assert.equal(board.submitting.value, false)
})

test('POST 失败保留输入，不请求列表', async () => {
  const board = createBoard([], async () => { throw new Error('发布失败') })
  await board.publishMessage({ nickname: '同学', content: '草稿' })
  assert.equal(board.cleared(), 0)
  assert.equal(board.messages.value.length, 0)
  assert.equal(board.errorMessage.value, '发布失败')
  assert.equal(board.submitting.value, false)
})
