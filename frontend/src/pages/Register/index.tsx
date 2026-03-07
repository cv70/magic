import { Form, Input, Button, Card, message, Spin } from 'antd'
import { UserOutlined, MailOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { useState } from 'react'
import { authApi } from '../../utils/api'
import { useAuthStore } from '../../stores/authStore'

export function RegisterPage() {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { login } = useAuthStore()

  const onFinish = async (values: any) => {
    if (values.password !== values.confirm_password) {
      message.error('两次输入的密码不一致')
      return
    }

    setLoading(true)
    try {
      const result = await authApi.register(values.username, values.email, values.password)
      login(result.token, result.user)
      message.success('注册成功')
      navigate('/dashboard')
    } catch (error) {
      message.error(error instanceof Error ? error.message : '注册失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f5f5f5' }}>
      <Card style={{ width: 400 }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <h1>🚀 Magic Engine</h1>
          <p style={{ color: '#999' }}>创建账户</p>
        </div>

        <Spin spinning={loading}>
          <Form onFinish={onFinish} layout="vertical">
            <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
              <Input prefix={<UserOutlined />} placeholder="用户名" />
            </Form.Item>

            <Form.Item name="email" rules={[{ required: true, type: 'email', message: '请输入有效的邮箱' }]}>
              <Input prefix={<MailOutlined />} placeholder="邮箱" />
            </Form.Item>

            <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
              <Input.Password prefix={<LockOutlined />} placeholder="密码" />
            </Form.Item>

            <Form.Item name="confirm_password" rules={[{ required: true, message: '请确认密码' }]}>
              <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" block loading={loading}>
                注册
              </Button>
            </Form.Item>

            <div style={{ textAlign: 'center' }}>
              已有账户？<Link to="/login">立即登录</Link>
            </div>
          </Form>
        </Spin>
      </Card>
    </div>
  )
}
