import { useState } from 'react'
import { Card, Form, Input, Button, Select, Space, message, Spin, Tabs, Table, Modal, Divider, Row, Col, Switch } from 'antd'
import { PlusOutlined, DeleteOutlined, EditOutlined } from '@ant-design/icons'

export function Settings() {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [activeTab, setActiveTab] = useState('ai')
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)

  // Mock AI Generators
  const [generators, setGenerators] = useState([
    {
      id: 1,
      name: 'GPT-4',
      provider: 'openai',
      model: 'gpt-4',
      api_key: '***',
      endpoint: 'https://api.openai.com/v1',
      enabled: true,
    },
    {
      id: 2,
      name: 'Claude',
      provider: 'anthropic',
      model: 'claude-3-opus',
      api_key: '***',
      endpoint: 'https://api.anthropic.com',
      enabled: false,
    },
  ])

  // Mock Publishers
  const [publishers, setPublishers] = useState([
    { id: 1, name: '微信公众号', platform: 'wechat', platform_id: 'gh_xxx', enabled: true },
    { id: 2, name: '新浪微博', platform: 'weibo', platform_id: '@user', enabled: true },
    { id: 3, name: '抖音', platform: 'douyin', platform_id: '@user123', enabled: false },
  ])

  // Mock Prompt Templates
  const [templates] = useState([
    { id: 1, name: '博客文章', category: '内容创作', description: '生成SEO友好的博客文章' },
    { id: 2, name: '社交媒体', category: '内容创作', description: '生成吸引人的社交帖子' },
    { id: 3, name: '产品描述', category: '电商', description: '生成详细的产品描述' },
  ])

  const handleAddGenerator = async (values: any) => {
    try {
      setLoading(true)
      if (editingId) {
        // Update existing
        setGenerators(
          generators.map((g) =>
            g.id === editingId ? { ...g, ...values } : g
          )
        )
        message.success('生成器已更新')
        setEditingId(null)
      } else {
        // Add new
        setGenerators([...generators, { ...values, id: Date.now(), enabled: true }])
        message.success('生成器已添加')
      }
      form.resetFields()
      setIsModalVisible(false)
    } catch (error) {
      message.error('操作失败')
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteGenerator = (id: number) => {
    Modal.confirm({
      title: '确定删除？',
      onOk: () => {
        setGenerators(generators.filter((g) => g.id !== id))
        message.success('已删除')
      },
    })
  }

  const handleDeletePublisher = (id: number) => {
    Modal.confirm({
      title: '确定删除？',
      onOk: () => {
        setPublishers(publishers.filter((p) => p.id !== id))
        message.success('已删除')
      },
    })
  }

  const generatorColumns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '供应商', dataIndex: 'provider', key: 'provider' },
    { title: '模型', dataIndex: 'model', key: 'model' },
    { title: 'API端点', dataIndex: 'endpoint', key: 'endpoint', width: 200 },
    {
      title: '启用',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean) => <Switch checked={enabled} />,
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: any) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => {
              setEditingId(record.id)
              form.setFieldsValue(record)
              setIsModalVisible(true)
            }}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => handleDeleteGenerator(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  const publisherColumns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '平台', dataIndex: 'platform', key: 'platform' },
    { title: '平台ID', dataIndex: 'platform_id', key: 'platform_id' },
    {
      title: '启用',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean) => <Switch checked={enabled} />,
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: any) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => {
              setEditingId(record.id)
              form.setFieldsValue(record)
              setIsModalVisible(true)
            }}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => handleDeletePublisher(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  const templateColumns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '分类', dataIndex: 'category', key: 'category' },
    { title: '描述', dataIndex: 'description', key: 'description' },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: () => (
        <Space size="small">
          <Button type="link" size="small" icon={<EditOutlined />}>
            编辑
          </Button>
          <Button type="link" danger size="small" icon={<DeleteOutlined />}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  const tabItems = [
    {
      key: 'ai',
      label: 'AI生成器',
      children: (
        <div>
          <div style={{ marginBottom: 16 }}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setEditingId(null)
                form.resetFields()
                setIsModalVisible(true)
              }}
            >
              添加生成器
            </Button>
          </div>
          <Table columns={generatorColumns} dataSource={generators} rowKey="id" />
        </div>
      ),
    },
    {
      key: 'publishers',
      label: '发布平台',
      children: (
        <div>
          <div style={{ marginBottom: 16 }}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setEditingId(null)
                form.resetFields()
                setIsModalVisible(true)
              }}
            >
              添加平台
            </Button>
          </div>
          <Table columns={publisherColumns} dataSource={publishers} rowKey="id" />
        </div>
      ),
    },
    {
      key: 'templates',
      label: '提示词模板',
      children: (
        <div>
          <div style={{ marginBottom: 16 }}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setEditingId(null)
                form.resetFields()
                setIsModalVisible(true)
              }}
            >
              添加模板
            </Button>
          </div>
          <Table columns={templateColumns} dataSource={templates} rowKey="id" />
        </div>
      ),
    },
    {
      key: 'general',
      label: '通用设置',
      children: (
        <Card>
          <Form layout="vertical">
            <Row gutter={16}>
              <Col xs={24} sm={12}>
                <Form.Item label="API请求超时（秒）">
                  <Input type="number" defaultValue={30} />
                </Form.Item>
              </Col>
              <Col xs={24} sm={12}>
                <Form.Item label="最大重试次数">
                  <Input type="number" defaultValue={3} />
                </Form.Item>
              </Col>
            </Row>

            <Divider />

            <div style={{ marginTop: 16 }}>
              <Form.Item label="">
                <Space direction="vertical" style={{ width: '100%' }}>
                  <Form.Item name="auto_save" label="自动保存草稿">
                    <Switch defaultChecked />
                  </Form.Item>
                  <Form.Item name="notifications" label="启用通知提醒">
                    <Switch defaultChecked />
                  </Form.Item>
                  <Form.Item name="dark_mode" label="暗黑模式">
                    <Switch />
                  </Form.Item>
                </Space>
              </Form.Item>
            </div>

            <Divider />

            <Space>
              <Button type="primary">保存设置</Button>
              <Button>恢复默认</Button>
            </Space>
          </Form>
        </Card>
      ),
    },
  ]

  return (
    <Spin spinning={loading}>
      <div>
        <h1>⚙️ 系统设置</h1>
        <div style={{ marginBottom: 24, color: '#666' }}>
          管理AI生成器、发布平台、提示词模板和系统偏好设置
        </div>

        <Card>
          <Tabs
            activeKey={activeTab}
            onChange={setActiveTab}
            items={tabItems}
          />
        </Card>

        {/* 模态框：添加/编辑 */}
        <Modal
          title={editingId ? '编辑' : '添加'}
          open={isModalVisible}
          onCancel={() => {
            setIsModalVisible(false)
            setEditingId(null)
            form.resetFields()
          }}
          onOk={() => form.submit()}
        >
          <Form form={form} layout="vertical" onFinish={handleAddGenerator}>
            <Form.Item
              name="name"
              label="名称"
              rules={[{ required: true, message: '请输入名称' }]}
            >
              <Input />
            </Form.Item>

            <Form.Item
              name="provider"
              label="供应商"
              rules={[{ required: true, message: '请选择供应商' }]}
            >
              <Select
                placeholder="选择供应商"
                options={[
                  { label: 'OpenAI', value: 'openai' },
                  { label: 'Anthropic', value: 'anthropic' },
                  { label: 'Google', value: 'google' },
                ]}
              />
            </Form.Item>

            <Form.Item
              name="model"
              label="模型"
              rules={[{ required: true, message: '请输入模型名称' }]}
            >
              <Input placeholder="例：gpt-4" />
            </Form.Item>

            <Form.Item
              name="api_key"
              label="API密钥"
              rules={[{ required: true, message: '请输入API密钥' }]}
            >
              <Input.Password />
            </Form.Item>

            <Form.Item name="endpoint" label="API端点">
              <Input />
            </Form.Item>
          </Form>
        </Modal>
      </div>
    </Spin>
  )
}
