// 对运行中的真实服务验收；每次保留 3 条新留言，不删除已有数据。
import assert from 'node:assert/strict'

const base = (process.env.API_BASE || 'http://127.0.0.1:5173').replace(/\/$/, '')
const url = `${base}/api/messages`
const tag = `验收-${Date.now()}`

async function list() {
  const response = await fetch(url)
  assert.equal(response.status, 200, 'GET 应返回 200')
  const messages = await response.json()
  assert.ok(Array.isArray(messages), 'GET 应返回 JSON 数组')
  return messages
}

async function post(body, expectedStatus, contentType = 'application/json') {
  const response = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': contentType },
    body: typeof body === 'string' ? body : JSON.stringify(body),
  })
  const data = await response.json()
  assert.equal(response.status, expectedStatus, JSON.stringify(data))
  if (expectedStatus !== 201) assert.equal(typeof data.message, 'string')
  return data
}

const before = await list()
const created = []
for (const input of [
  { nickname: '  新生同学  ', content: `${tag}：第一条校园留言` },
  { nickname: '😀'.repeat(30), content: '🌱'.repeat(200) },
  { nickname: "O'Reilly", content: `${tag}：' OR 1=1 -- <script>alert(1)</script>` },
]) {
  const saved = await post(input, 201)
  assert.equal(saved.nickname, input.nickname.trim())
  assert.equal(saved.content, input.content)
  assert.ok(Number.isInteger(saved.id) && saved.id > 0)
  assert.equal(typeof saved.createdAt, 'string')
  created.push(saved)
}

const invalidInputs = [
  {}, { nickname: null, content: '留言' },
  { nickname: ' \t\u00a0\ufeff ', content: '留言' },
  { nickname: '名'.repeat(31), content: '留言' },
  { nickname: '😀'.repeat(31), content: '留言' },
  { nickname: '同学', content: '' },
  { nickname: '同学', content: ' \n\u00a0\ufeff ' },
  { nickname: '同学', content: '字'.repeat(201) },
  { nickname: '同学', content: '🌱'.repeat(201) },
  { nickname: 123, content: '留言' },
  { nickname: '同学', content: false },
  'null', '{invalid json',
]
for (const input of invalidInputs) await post(input, 400)
await post('nickname=test', 415, 'text/plain')

const after = await list()
assert.equal(after.length, before.length + 3, '无效请求不能新增记录；验收期间请勿并发发布')
for (const saved of created) assert.deepEqual(after.find(item => item.id === saved.id), saved)
for (let i = 1; i < after.length; i++) assert.ok(after[i - 1].id > after[i].id, '列表按 id 倒序')
console.log(JSON.stringify({ result: 'PASS', base, createdIds: created.map(item => item.id), rejectedRequests: invalidInputs.length + 1, messageCount: after.length }, null, 2))
console.log('HTTP 验收通过。请继续在指定 MySQL 库中 SELECT 核对这些 id，并重启后端验证数据保留。')
