import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Space, Modal, Form, Input, message, Popconfirm, Select, Card } from 'antd'
import { EditOutlined, DeleteOutlined, PlusOutlined, SearchOutlined } from '@ant-design/icons'
import { draftApi, tagApi, categoryApi } from '../../utils/api'
import type { Draft } from '../../types'

export function DraftCenter() {
  const [drafts, setDrafts] = useState<Draft[]>([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20, total: 0 })
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const [searchKeyword, setSearchKeyword] = useState('')
  const [selectedTags, setSelectedTags] = useState<number[]>([])
  const [selectedRows, setSelectedRows] = useState<Draft[]>([])
  const [tags, setTags] = useState<any[]>([])
  const [categories, setCategories] = useState<any[]>([])
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null)

  useEffect(() => {
    loadDrafts()
    loadTags()
    loadCategories()
  }, [pagination.current, searchKeyword, selectedTags, selectedCategory])

  const loadDrafts = async () => {
    try {
      setLoading(true)
      const result = await draftApi.list(
        pagination.current,
        pagination.pageSize,
        searchKeyword,
        selectedCategory?.toString()
      )
      setDrafts(result.data || [])
      setPagination({ ...pagination, total: result.total })
    } catch (error) {
      message.error('加载草稿失败')
    } finally {
      setLoading(false)
    }
  }

  const loadTags = async () => {
    try {
      const result = await tagApi.list()
      setTags(result.data || [])
    } catch (error) {
      console.error('加载标签失败')
    }
  }

  const loadCategories = async () => {
    try {
      const result = await categoryApi.list()
      setCategories(result.data || [])
    } catch (error) {
      console.error('加载分类失败')
    }
  }

  const handleBatchDelete = async () => {
    if (selectedRows.length === 0) {
      message.warning('请先选择要删除的草稿')
      return
    }

    Modal.confirm({
      title: `确定删除 ${selectedRows.length} 个草稿？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          for (const draft of selectedRows) {
            await draftApi.delete(draft.id)
          }
          message.success('批量删除成功')
          setSelectedRows([])
          loadDrafts()
        } catch (error) {
          message.error('批量删除失败')
        }
      },
    })
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
      {/* 搜索和过滤栏 */}
      <Card style={{ marginBottom: 16 }}>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', alignItems: 'center' }}>
            <Input
              placeholder="搜索草稿标题..."
              prefix={<SearchOutlined />}
              value={searchKeyword}
              onChange={(e) => setSearchKeyword(e.target.value)}
              style={{ width: 250 }}
            />
            <Select
              placeholder="按分类筛选"
              allowClear
              style={{ width: 150 }}
              value={selectedCategory}
              onChange={setSelectedCategory}
              options={[
                ...categories.map((c) => ({ label: c.name, value: c.id })),
              ]}
            />
            <Select
              mode="multiple"
              placeholder="按标签筛选"
              style={{ width: 200 }}
              value={selectedTags}
              onChange={setSelectedTags}
              options={tags.map((t) => ({ label: t.name, value: t.id }))}
            />
          </div>

          <Space wrap>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => setIsModalVisible(true)}
            >
              新建草稿
            </Button>
            {selectedRows.length > 0 && (
              <>
                <span style={{ color: '#666' }}>
                  已选 {selectedRows.length} 个草稿
                </span>
                <Button danger onClick={handleBatchDelete}>
                  批量删除
                </Button>
              </>
            )}
          </Space>
        </Space>
      </Card>

      {/* 草稿列表 */}
      <Table
        columns={columns}
        dataSource={drafts}
        rowKey="id"
        loading={loading}
        rowSelection={{
          selectedRowKeys: selectedRows.map((r) => r.id),
          onChange: (_, rows) => setSelectedRows(rows as Draft[]),
        }}
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
