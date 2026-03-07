import { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { AppLayout } from './components/layout/AppLayout'
import { ProtectedRoute } from './components/layout/ProtectedRoute'
import { LoginPage } from './pages/Login'
import { RegisterPage } from './pages/Register'
import { Dashboard } from './pages/Dashboard'
import { DraftCenter } from './pages/DraftCenter'
import { Editor } from './pages/Editor'
import { PublishManager } from './pages/PublishManager'
import { useAuthStore } from './stores/authStore'

// 临时占位符页面
function AnalyticsPlaceholder() {
  return <div style={{ padding: 24 }}>分析页面 - 开发中</div>
}

function App() {
  const { init } = useAuthStore()

  useEffect(() => {
    init()
  }, [init])

  return (
    <Routes>
      {/* 不需要认证的路由 */}
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      {/* 需要认证的路由 */}
      <Route
        path="/*"
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="drafts" element={<DraftCenter />} />
        <Route path="drafts/:id/edit" element={<Editor />} />
        <Route path="drafts/new" element={<Editor />} />
        <Route path="analytics" element={<AnalyticsPlaceholder />} />
        <Route path="publish" element={<PublishManager />} />
      </Route>

      {/* 默认重定向 */}
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  )
}

export default App
