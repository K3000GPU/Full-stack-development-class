import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ command }) => ({
  // GitHub Pages 项目站点位于 /<仓库名>/，开发服务器仍使用根路径。
  base: command === 'build' ? '/Full-stack-development-class/' : '/',
  plugins: [vue()],
  server: {
    host: process.env.DEV_SERVER_HOST || '127.0.0.1',
    port: 5173,
    strictPort: true,
    proxy: {
      // 浏览器请求 /api；开发服务转发到后端，保留完整的 /api 路径。
      '/api': {
        target: process.env.API_PROXY_TARGET || 'http://127.0.0.1:8080',
        changeOrigin: true,
      },
    },
  },
}))
