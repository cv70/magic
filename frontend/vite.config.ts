import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    chunkSizeWarningLimit: 600, // 增加警告阈值到 600KB
    rollupOptions: {
      output: {
        manualChunks: {
          // 分离大依赖
          'vendor-antd': ['antd', '@ant-design/icons'],
          'vendor-recharts': ['recharts'],
          'vendor-tiptap': ['@tiptap/react', '@tiptap/starter-kit'],
          'vendor-misc': ['zustand', 'date-fns', 'react-router-dom'],
          // 按页面分割
          'pages-auth': [
            './src/pages/Login',
            './src/pages/Register',
          ],
          'pages-content': [
            './src/pages/DraftCenter',
            './src/pages/Editor',
          ],
          'pages-publish': [
            './src/pages/PublishManager',
          ],
          'pages-analytics': [
            './src/pages/Analytics',
          ],
          'pages-ai': [
            './src/pages/AIStudio',
          ],
        },
      },
    },
  },
})

