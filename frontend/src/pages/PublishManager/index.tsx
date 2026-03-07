import { useEffect, useState } from 'react'
import { Table, Button, Space, Modal, Form, Select, message, Tag, Spin } from 'antd'
import { ReloadOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons'
import { publishApi, contentApi, publisherApi } from '../../utils/api'

const statusColorMap: { [key: string]: string } = {
  pending: 'orange',
  publishing: 'blue',
  published: 'green',
  failed: 'red',
}

export function PublishManager() {
  const [tasks, setTasks] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [publishers, setPublishers] = useState<any[]>([])
  const [contents, setContents] = useState<any[]>([])
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20, total: 0 })

  useEffect(() => {
    loadTasks()
    loadPublishers()
    loadContents()
  }, [pagination.current])

  const loadTasks = async () => {
    try {
      setLoading(true)
      const result = await publishApi.list(pagination.current, pagination.pageSize)
      setTasks(result.data || [])
      setPagination({ ...pagination, total: result.total })
    } catch (error) {
      message.error('加载发布任务失败')
    } finally {
      setLoading(false)
    }
  }

  const loadPublishers = async () => {
    try {
      const result = await publisherApi.search()
      setPublishers(result.data || [])
    } catch (error) {
      console.error('加载发布平台失败')
      setPublishers([])
    }
  }

  const loadContents = async () => {
    try {
      const result = await contentApi.list(1, 1000)
      setContents(result.data || [])
    } catch (error) {
      console.error('加载内容失败')
    }
  }

  const handleCreateTask = async (values: any) => {
    try {
      await publishApi.create({
        content_id: values.content_id,
        target_platforms: [values.publisher_id],
      })
      message.success('发布任务已创建')
      setIsModalVisible(false)
      form.resetFields()
      loadTasks()
    } catch (error) {
      message.error('创建发布任务失败')
    }
  }

  const handleExecuteTask = async (taskId: number) => {
    try {
      await publishApi.execute(taskId)
      message.success('任务已执行')
      loadTasks()
    } catch (error) {
      message.error('执行任务失败')
    }
  }

  const handleDeleteTask = async (taskId: number) => {
    try {
      await publishApi.delete(taskId)
      message.success('任务已删除')
      loadTasks()
    } catch (error) {
      message.error('删除任务失败')
    }
  }

  const columns = [
    {
      title: '内容ID',
      dataIndex: 'content_id',
      key: 'content_id',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => (
        <Tag color={statusColorMap[status] || 'default'}>{status}</Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '完成时间',
      dataIndex: 'completed_at',
      key: 'completed_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render: (_: any, record: any) => (
        <Space size="small">
          {record.status === 'pending' && (
            <Button
              type="primary"
              size="small"
              onClick={() => handleExecuteTask(record.id)}
            >
              执行
            </Button>
          )}
          {record.status === 'failed' && (
            <Button
              size="small"
              onClick={() => handleExecuteTask(record.id)}
            >
              重试
            </Button>
          )}
          <Button
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => {
              Modal.confirm({
                title: '确定删除？',
                onOk: () => handleDeleteTask(record.id),
              })
            }}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', gap: 8 }}>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setIsModalVisible(true)}
        >
          新建发布任务
        </Button>
        <Button icon={<ReloadOutlined />} onClick={loadTasks}>
          刷新
        </Button>
      </div>

      <Spin spinning={loading}>
        <Table
          columns={columns}
          dataSource={tasks}
          rowKey="id"
          pagination={{
            current: pagination.current,
            pageSize: pagination.pageSize,
            total: pagination.total,
            onChange: (page) => setPagination({ ...pagination, current: page }),
          }}
        />
      </Spin>

      <Modal
        title="新建发布任务"
        open={isModalVisible}
        onOk={() => form.submit()}
        onCancel={() => {
          setIsModalVisible(false)
          form.resetFields()
        }}
      >
        <Form form={form} layout="vertical" onFinish={handleCreateTask}>
          <Form.Item
            name="content_id"
            label="选择内容"
            rules={[{ required: true, message: '请选择内容' }]}
          >
            <Select placeholder="选择要发布的内容">
              {contents.map((c) => (
                <Select.Option key={c.id} value={c.id}>
                  {c.title}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="publisher_id"
            label="选择发布平台"
            rules={[{ required: true, message: '请选择平台' }]}
          >
            <Select placeholder="选择发布平台">
              {publishers.map((p) => (
                <Select.Option key={p.id} value={p.id}>
                  {p.name || p.platform}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
