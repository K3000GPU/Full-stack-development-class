import { request } from './http'

// 接口层：说明留言接口的地址、HTTP 方法和请求数据。
export function listMessages() {
  return request('/api/messages')
}

export function createMessage(message) {
  return request('/api/messages', {
    method: 'POST',
    body: JSON.stringify(message),
  })
}
