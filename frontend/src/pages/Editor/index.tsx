import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Form, Input, Button, Space, message, Spin, Select, Card } from 'antd'
import { SaveOutlined, SendOutlined, ArrowLeftOutlined } from '@ant-design/icons'
import { draftApi } from '../../utils/api'
import { TipTapEditor } from '../../components/TipTapEditor'
import { AutoSaveIndicator } from '../../components/AutoSaveIndicator'
import { useAutoSave } from '../../hooks/useAutoSave'

export function Editor() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(!!id)
  const [saving, setSaving] = useState(false)
  const [publishers, setPublishers] = useState<any[]>([])
  const [selectedPublishers, setSelectedPublishers] = useState<number[]>([])
  const [isPublishing, setIsPublishing] = useState(false)
  const [content, setContent] = useState('')

  // 自动保存
  const { status: autoSaveStatus, updateData, manualSave } = useAutoSave({
    interval: 30000, // 30秒
    onSave: async (data) => {
      if (id && form.getFieldValue('title')) {
        await draftApi.update(parseInt(id), {
          title: form.getFieldValue('title'),
          content_type: form.getFieldValue('content_type'),
          content: data.content,
        })
      }
    },
  })

  useEffect(() => {
    if (id) {
      loadDraft()
    }
    loadPublishers()
  }, [id])

  const loadDraft = async () => {
    try {
      const draft = await draftApi.get(parseInt(id!))
      form.setFieldsValue({
        title: draft.title,
        content_type: draft.content_type,
      })
      setContent(draft.content)
    } catch (error) {
      message.error('加载草稿失败')
      navigate('/drafts')
    } finally {
      setLoading(false)
    }
  }

  const loadPublishers = async () => {
    try {
      // TODO: Load publishers from API
      // For now, use mock data
      setPublishers([
        { id: '1', platform: 'WeChat' },
        { id: '2', platform: 'Weibo' },
        { id: '3', platform: 'Douyin' },
      ])
    } catch (error) {
      console.error('加载发布平台失败')
    }
  }

  const handleSave = async (values: any) => {
    setSaving(true)
    try {
      const data = {
        ...values,
        content,
      }
      if (id) {
        await draftApi.update(parseInt(id), data)
        message.success('草稿已保存')
      } else {
        const result = await draftApi.create(data)
        navigate(`/drafts/${result.id}/edit`)
        message.success('草稿已创建')
      }
    } catch (error) {
      message.error('保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handlePublish = async () => {
    if (selectedPublishers.length === 0) {
      message.error('请选择发布平台')
      return
    }

    setIsPublishing(true)
    try {
      // 先发布为正式内容
      if (id) {
        await draftApi.publish(parseInt(id), 'published')
        message.success('已发布为正式内容')
        navigate('/drafts')
      }
    } catch (error) {
      message.error('发布失败')
    } finally {
      setIsPublishing(false)
    }
  }

  return (
    <Spin spinning={loading}>
      <div>
        <div style={{ marginBottom: 24 }}>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/drafts')}>
            返回草稿中心
          </Button>
        </div>

        <Form form={form} layout="vertical" onFinish={handleSave}>
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input placeholder="输入文章标题" size="large" />
          </Form.Item>

          <Form.Item
            name="content_type"
            label="内容类型"
            initialValue="html"
          >
            <Select>
              <Select.Option value="text">纯文本</Select.Option>
              <Select.Option value="markdown">Markdown</Select.Option>
              <Select.Option value="html">富文本 (HTML)</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item label="内容编辑器">
            <Card style={{ marginBottom: 16 }}>
              <div style={{ marginBottom: 12, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span style={{ fontSize: '12px', color: '#666' }}>
                  自动保存每 30 秒
                </span>
                <AutoSaveIndicator {...autoSaveStatus} />
              </div>
              <TipTapEditor
                value={content}
                onChange={(newContent) => {
                  setContent(newContent)
                  updateData({ content: newContent })
                }}
                placeholder="开始输入内容...支持富文本、Markdown 和源码编辑（自动保存）"
                height="500px"
              />
            </Card>
          </Form.Item>

          <Form.Item>
            <Space wrap>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                loading={saving}
                onClick={() => form.submit()}
              >
                保存草稿
              </Button>

              {id && (
                <Button onClick={manualSave} loading={autoSaveStatus.isSaving}>
                  立即保存
                </Button>
              )}

              <div style={{ borderLeft: '1px solid #ddd', paddingLeft: 16, paddingRight: 16 }}>
                <div style={{ marginBottom: 8 }}>发布到平台：</div>
                <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                  <Select
                    mode="multiple"
                    placeholder="选择要发布的平台"
                    style={{ width: 300, minWidth: 200 }}
                    value={selectedPublishers}
                    onChange={setSelectedPublishers}
                    options={publishers.map((p) => ({
                      value: p.id,
                      label: p.platform,
                    }))}
                  />
                  <Button
                    type="primary"
                    danger
                    icon={<SendOutlined />}
                    loading={isPublishing}
                    onClick={handlePublish}
                  >
                    发布
                  </Button>
                </div>
              </div>
            </Space>
          </Form.Item>
        </Form>
      </div>
    </Spin>
  )
}
