import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Space, Modal, Form, Input, message, Popconfirm } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons'
import { draftApi } from '../../utils/api'
import type { Draft } from '../../types'

export function DraftCenter() {
  const [drafts, setDrafts] = useState<Draft[]>([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20, total: 0 })
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()

  useEffect(() => {
    loadDrafts()
  }, [pagination.current])

  const loadDrafts = async () => {
    try {
      setLoading(true)
      const result = await draftApi.list(pagination.current, pagination.pageSize)
      setDrafts(result.data || [])
      setPagination({ ...pagination, total: result.total })
    } catch (error) {
      message.error('加载草稿失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateDraft = async (values: any) => {
    try {
      const result = await draftApi.create({
        title: values.title,
        content: values.content || '',
      })
      message.success('草稿已创建')
      setIsModalVisible(false)
      form.resetFields()
      navigate(`/drafts/${result.id}/edit`)
    } catch (error) {
      message.error('创建草稿失败')
    }
  }

  const handleEditDraft = (draft: Draft) => {
    navigate(`/drafts/${draft.id}/edit`)
  }

  const handleDeleteDraft = async (id: number) => {
    try {
      await draftApi.delete(id)
      message.success('草稿已删除')
      loadDrafts()
    } catch (error) {
      message.error('删除草稿失败')
    }
  }

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: Draft) => (
        <Space size="small">
          <Button
            type="primary"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEditDraft(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定删除？"
            onConfirm={() => handleDeleteDraft(record.id)}
          >
            <Button danger size="small" icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setIsModalVisible(true)}
        >
          新建草稿
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={drafts}
        rowKey="id"
        loading={loading}
        pagination={{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          onChange: (page) => setPagination({ ...pagination, current: page }),
        }}
      />

      <Modal
        title="新建草稿"
        open={isModalVisible}
        onOk={() => form.submit()}
        onCancel={() => {
          setIsModalVisible(false)
          form.resetFields()
        }}
      >
        <Form form={form} layout="vertical" onFinish={handleCreateDraft}>
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input placeholder="输入草稿标题" />
          </Form.Item>
          <Form.Item name="content" label="内容（可选）">
            <Input.TextArea rows={4} placeholder="输入初始内容" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
