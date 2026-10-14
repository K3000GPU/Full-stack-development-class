// 请求层：统一发送 HTTP 请求，并把网络和后端错误变成可展示的 Error。
export async function request(path, options = {}) {
  let response

  try {
    response = await fetch(path, {
      ...options,
      headers: {
        Accept: 'application/json',
        ...(options.body ? { 'Content-Type': 'application/json' } : {}),
        ...options.headers,
      },
    })
  } catch {
    throw new Error('无法连接服务，请检查前后端是否已启动。')
  }

  let data
  try {
    data = await response.json()
  } catch {
    throw new Error(`服务返回了无法读取的响应（HTTP ${response.status}）。`)
  }

  if (!response.ok) {
    throw new Error(data.message || `请求失败（HTTP ${response.status}）。`)
  }

  return data
}
