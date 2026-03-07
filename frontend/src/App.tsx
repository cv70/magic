import { useEffect, Suspense, lazy } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { Spin } from 'antd'
import { AppLayout } from './components/layout/AppLayout'
import { ProtectedRoute } from './components/layout/ProtectedRoute'
import { LoginPage } from './pages/Login'
import { RegisterPage } from './pages/Register'
import { useAuthStore } from './stores/authStore'

// 懒加载页面
const Dashboard = lazy(() => import('./pages/Dashboard').then(m => ({ default: m.Dashboard })))
const DraftCenter = lazy(() => import('./pages/DraftCenter').then(m => ({ default: m.DraftCenter })))
const Editor = lazy(() => import('./pages/Editor').then(m => ({ default: m.Editor })))
const PublishManager = lazy(() => import('./pages/PublishManager').then(m => ({ default: m.PublishManager })))
const Analytics = lazy(() => import('./pages/Analytics').then(m => ({ default: m.Analytics })))
const AIStudio = lazy(() => import('./pages/AIStudio').then(m => ({ default: m.AIStudio })))
const Settings = lazy(() => import('./pages/Settings').then(m => ({ default: m.Settings })))

// 加载中组件
const LoadingSpinner = () => (
  <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
    <Spin size="large" tip="加载中..." />
  </div>
)

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
        <Route
          path="dashboard"
          element={
            <Suspense fallback={<LoadingSpinner />}>
              <Dashboard />
            </Suspense>
          }
        />
        <Route
          path="drafts"
          element={
            <Suspense fallback={<LoadingSpinner />}>
              <DraftCenter />
            </Suspense>
          }
        />
        <Route
          path="drafts/:id/edit"
          element={
            <Suspense fallback={<LoadingSpinner />}>
              <Editor />
            </Suspense>
          }
        />
        <Route
          path="drafts/new"
          element={
            <Suspense fallback={<LoadingSpinner />}>
              <Editor />
            </Suspense>
          }
        />
        <Route
          path="analytics"
          element={
            <Suspense fallback={<LoadingSpinner />}>
              <Analytics />
            </Suspense>
          }
        />
        <Route
          path="publish"
          element={
            <Suspense fallback={<LoadingSpinner />}>
              <PublishManager />
            </Suspense>
          }
        />
        <Route
          path="ai-studio"
          element={
            <Suspense fallback={<LoadingSpinner />}>
              <AIStudio />
            </Suspense>
          }
        />
        <Route
          path="settings"
          element={
            <Suspense fallback={<LoadingSpinner />}>
              <Settings />
            </Suspense>
          }
        />
      </Route>

      {/* 默认重定向 */}
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  )
}

export default App
