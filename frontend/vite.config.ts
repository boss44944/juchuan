import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    outDir: 'dist',
    // 指定旧浏览器兼容目标：让 Lightning CSS 输出传统 max-width 媒体查询语法
    // （如 @media (max-width:780px)），而不是新式 range 语法 (width<=780px)，
    // 否则 Android 旧版 Chrome（如 Chrome 90）会忽略移动端适配样式。
    cssTarget: ['chrome90', 'safari14', 'firefox100', 'edge90']
  },
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8000',
      '/ws': 'ws://localhost:8000',
      '/upload': 'http://localhost:8000'
    }
  }
})
