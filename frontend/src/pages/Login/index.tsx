import { Form, Input, Button, Card, message, Spin } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { useState } from 'react'
import { authApi } from '../../utils/api'
import { useAuthStore } from '../../stores/authStore'

export function LoginPage() {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { login } = useAuthStore()

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true)
    try {
      const result = await authApi.login(values.username, values.password)
      login(result.token, result.user)
      message.success('登录成功')
      navigate('/dashboard')
    } catch (error) {
      message.error(error instanceof Error ? error.message : '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f5f5f5' }}>
      <Card style={{ width: 400 }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <h1>🚀 Magic Engine</h1>
          <p style={{ color: '#999' }}>AI 内容生产平台</p>
        </div>

        <Spin spinning={loading}>
          <Form onFinish={onFinish} layout="vertical">
            <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
              <Input prefix={<UserOutlined />} placeholder="用户名" />
            </Form.Item>

            <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
              <Input.Password prefix={<LockOutlined />} placeholder="密码" />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" block loading={loading}>
                登录
              </Button>
            </Form.Item>

            <div style={{ textAlign: 'center' }}>
              没有账户？<Link to="/register">立即注册</Link>
            </div>
          </Form>
        </Spin>
      </Card>
    </div>
  )
}
