import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Form, Input, Button, Space, message, Spin, Select, Card, Modal, Timeline, DatePicker } from 'antd'
import { SaveOutlined, SendOutlined, ArrowLeftOutlined, HistoryOutlined, ScheduleOutlined } from '@ant-design/icons'
import { draftApi, publishApi } from '../../utils/api'
import { publisherApi } from '../../utils/api'
import { TipTapEditor } from '../../components/TipTapEditor'
import { AutoSaveIndicator } from '../../components/AutoSaveIndicator'
import { useAutoSave } from '../../hooks/useAutoSave'
import dayjs from 'dayjs'

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
  const [versions, setVersions] = useState<any[]>([])
  const [showVersionModal, setShowVersionModal] = useState(false)
  const [showScheduleModal, setShowScheduleModal] = useState(false)
  const [scheduledTime, setScheduledTime] = useState<any>(null)

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
      const result = await publisherApi.search()
      setPublishers(result.data || [])
    } catch (error) {
      console.error('加载发布平台失败')
      // Fallback to empty array on error
      setPublishers([])
    }
  }

  const loadVersions = async () => {
    if (!id) return
    try {
      const result = await draftApi.getVersions(parseInt(id))
      setVersions(result.data || [])
    } catch (error) {
      message.error('加载版本历史失败')
    }
  }

  const handleRevertVersion = async (versionNum: number) => {
    if (!id) return
    try {
      await draftApi.revertVersion(parseInt(id), versionNum)
      message.success('已恢复到该版本')
      loadDraft()
      setShowVersionModal(false)
    } catch (error) {
      message.error('恢复版本失败')
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
      // 创建发布任务（支持定时发布）
      if (id) {
        await publishApi.create({
          draft_id: parseInt(id),
          target_platforms: selectedPublishers.map(p => p.toString()),
          scheduled_at: scheduledTime ? scheduledTime.toISOString() : undefined,
        })
        message.success(scheduledTime ? '发布任务已创建，将在指定时间发布' : '已发布')
        setShowScheduleModal(false)
        navigate('/publish')
      }
    } catch (error) {
      message.error('发布失败')
    } finally {
      setIsPublishing(false)
    }
  }

  const openScheduleModal = () => {
    if (selectedPublishers.length === 0) {
      message.error('请先选择要发布的平台')
      return
    }
    setShowScheduleModal(true)
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
                <>
                  <Button onClick={manualSave} loading={autoSaveStatus.isSaving}>
                    立即保存
                  </Button>
                  <Button
                    icon={<HistoryOutlined />}
                    onClick={() => {
                      loadVersions()
                      setShowVersionModal(true)
                    }}
                  >
                    版本历史
                  </Button>
                </>
              )}

              <div style={{ borderLeft: '1px solid #ddd', paddingLeft: 16, paddingRight: 16 }}>
                <div style={{ marginBottom: 8 }}>发布到平台：</div>
                <div style={{ display: 'flex', gap: 8, alignItems: 'center', flexWrap: 'wrap' }}>
                  <Select
                    mode="multiple"
                    placeholder="选择要发布的平台"
                    style={{ width: 300, minWidth: 200 }}
                    value={selectedPublishers}
                    onChange={setSelectedPublishers}
                    options={publishers.map((p) => ({
                      value: p.id,
                      label: p.name || p.platform,
                    }))}
                  />
                  <Button
                    type="primary"
                    icon={<SendOutlined />}
                    loading={isPublishing}
                    onClick={handlePublish}
                  >
                    立即发布
                  </Button>
                  <Button
                    icon={<ScheduleOutlined />}
                    onClick={openScheduleModal}
                  >
                    定时发布
                  </Button>
                </div>
              </div>
            </Space>
          </Form.Item>
        </Form>

        {/* 版本历史模态框 */}
        <Modal
          title="版本历史"
          open={showVersionModal}
          onCancel={() => setShowVersionModal(false)}
          footer={null}
          width={500}
        >
          <Timeline
            items={versions.map((v) => ({
              label: new Date(v.created_at).toLocaleString('zh-CN'),
              children: (
                <div>
                  <p><strong>版本 {v.version_num}</strong></p>
                  <p>{v.change_summary || '无描述'}</p>
                  <Button
                    size="small"
                    type="link"
                    onClick={() => handleRevertVersion(v.version_num)}
                  >
                    恢复此版本
                  </Button>
                </div>
              ),
            }))}
          />
        </Modal>

        {/* 定时发布模态框 */}
        <Modal
          title="定时发布"
          open={showScheduleModal}
          onOk={handlePublish}
          onCancel={() => setShowScheduleModal(false)}
          okText="发布"
          cancelText="取消"
          confirmLoading={isPublishing}
        >
          <div style={{ marginBottom: 16 }}>
            <strong>选择发布时间：</strong>
            <DatePicker
              showTime
              format="YYYY-MM-DD HH:mm:ss"
              placeholder="选择发布时间"
              value={scheduledTime}
              onChange={setScheduledTime}
              style={{ width: '100%', marginTop: 8 }}
              disabledDate={(current) => current && current < dayjs().startOf('day')}
            />
          </div>
          <div>
            <strong>目标平台：</strong>
            <div style={{ marginTop: 8 }}>
              {selectedPublishers.map((id) => {
                const pub = publishers.find((p) => p.id === id)
                return <div key={id}> ✓ {pub?.name || pub?.platform}</div>
              })}
            </div>
          </div>
        </Modal>
      </div>
    </Spin>
  )
}
