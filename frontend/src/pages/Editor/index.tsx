import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Form, Input, Button, Space, message, Spin, Select } from 'antd'
import { SaveOutlined, SendOutlined, ArrowLeftOutlined } from '@ant-design/icons'
import { draftApi } from '../../utils/api'

export function Editor() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(!!id)
  const [saving, setSaving] = useState(false)
  const [publishers, setPublishers] = useState<any[]>([])
  const [selectedPublishers, setSelectedPublishers] = useState<number[]>([])
  const [isPublishing, setIsPublishing] = useState(false)

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
        content: draft.content,
        content_type: draft.content_type,
      })
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
      if (id) {
        await draftApi.update(parseInt(id), values)
        message.success('草稿已保存')
      } else {
        const result = await draftApi.create(values)
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
            initialValue="text"
          >
            <Select>
              <Select.Option value="text">文本</Select.Option>
              <Select.Option value="markdown">Markdown</Select.Option>
              <Select.Option value="html">HTML</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入内容' }]}
          >
            <Input.TextArea
              rows={15}
              placeholder="输入文章内容（当前版本使用纯文本编辑器，后续会升级为富文本编辑器）"
            />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                loading={saving}
                onClick={() => form.submit()}
              >
                保存草稿
              </Button>

              <div style={{ borderLeft: '1px solid #ddd', paddingLeft: 16 }}>
                <div style={{ marginBottom: 8 }}>选择发布平台：</div>
                <Select
                  mode="multiple"
                  placeholder="选择要发布的平台"
                  style={{ width: 300 }}
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
                  style={{ marginTop: 8 }}
                >
                  发布到平台
                </Button>
              </div>
            </Space>
          </Form.Item>
        </Form>
      </div>
    </Spin>
  )
}
