import { useEffect, useState } from 'react'
import { Card, Form, Input, Button, Select, Space, message, Spin, Row, Col, Modal, Table } from 'antd'
import { SendOutlined, SaveOutlined } from '@ant-design/icons'

export function AIStudio() {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [generating, setGenerating] = useState(false)
  const [generatedContent, setGeneratedContent] = useState('')
  const [selectedGenerator, setSelectedGenerator] = useState('')
  const [selectedPrompt, setSelectedPrompt] = useState('')
  const [generationHistory, setGenerationHistory] = useState<any[]>([])
  const [showHistoryModal, setShowHistoryModal] = useState(false)

  // Mock data for AI generators
  const generators = [
    { id: 'gpt4', name: 'GPT-4', description: '最强大的AI模型' },
    { id: 'gpt35', name: 'GPT-3.5', description: '快速高效的AI模型' },
    { id: 'claude', name: 'Claude', description: '优秀的文本生成模型' },
  ]

  // Mock data for prompt templates
  const promptTemplates = [
    { id: 1, name: '博客文章', category: '内容创作', description: '生成SEO友好的博客文章' },
    { id: 2, name: '社交媒体', category: '内容创作', description: '生成吸引人的社交媒体帖子' },
    { id: 3, name: '产品描述', category: '电商', description: '生成详细的产品描述' },
    { id: 4, name: '新闻稿', category: '新闻', description: '生成专业的新闻稿' },
    { id: 5, name: '电子邮件', category: '营销', description: '生成营销邮件内容' },
  ]

  useEffect(() => {
    // Load generation history from localStorage
    const history = localStorage.getItem('ai_generation_history')
    if (history) {
      setGenerationHistory(JSON.parse(history))
    }
  }, [])

  const handleGenerate = async (values: any) => {
    if (!selectedGenerator) {
      message.error('请选择AI生成器')
      return
    }
    if (!selectedPrompt) {
      message.error('请选择提示词模板')
      return
    }

    setGenerating(true)
    try {
      // Simulate AI generation API call
      // In production, this would call: POST /api/v1/ai/generate
      await new Promise((resolve) => setTimeout(resolve, 2000))

      // Mock generated content
      const mockContent = `# ${values.title || '生成的内容标题'}\n\n这是由AI生成的内容示例。实际内容会根据您的提示和选择的生成器而不同。\n\n## 主要特点\n\n- 内容根据您的需求自动生成\n- 支持多种生成器选择\n- 快速高效的生成速度\n\n${values.prompt || '您可以在此添加额外的提示词来定制生成的内容。'}`

      setGeneratedContent(mockContent)

      // Save to history
      const newItem = {
        id: Date.now(),
        title: values.title || '未命名内容',
        generator: selectedGenerator,
        template: selectedPrompt,
        prompt: values.prompt,
        content: mockContent,
        createdAt: new Date().toISOString(),
      }

      const updatedHistory = [newItem, ...generationHistory].slice(0, 50)
      setGenerationHistory(updatedHistory)
      localStorage.setItem('ai_generation_history', JSON.stringify(updatedHistory))

      message.success('内容已生成！')
    } catch (error) {
      message.error('生成失败，请重试')
    } finally {
      setGenerating(false)
    }
  }

  const handleSavedAsDraft = async () => {
    if (!generatedContent) {
      message.error('没有生成的内容可保存')
      return
    }

    try {
      setLoading(true)
      const title = form.getFieldValue('title') || '生成的内容'

      // Save as draft via API: POST /api/v1/drafts/create
      // For now, just show success message
      message.success(`已保存为草稿：${title}`)

      // Clear form
      form.resetFields()
      setGeneratedContent('')
      setSelectedGenerator('')
      setSelectedPrompt('')
    } catch (error) {
      message.error('保存失败')
    } finally {
      setLoading(false)
    }
  }

  const historyColumns = [
    { title: '标题', dataIndex: 'title', key: 'title', width: 150 },
    { title: '生成器', dataIndex: 'generator', key: 'generator', width: 100 },
    { title: '模板', dataIndex: 'template', key: 'template', width: 100 },
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
      render: (text: string) => new Date(text).toLocaleString(),
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
            onClick={() => {
              setGeneratedContent(record.content)
              form.setFieldsValue({
                title: record.title,
                prompt: record.prompt,
              })
              setSelectedGenerator(record.generator)
              setSelectedPrompt(record.template)
              setShowHistoryModal(false)
            }}
          >
            恢复
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <Spin spinning={loading}>
      <div>
        <h1>✨ AI Studio</h1>
        <div style={{ marginBottom: 24, color: '#666' }}>
          使用AI强大功能快速生成高质量内容
        </div>

        <Row gutter={16}>
          {/* 左侧配置区 */}
          <Col xs={24} lg={12}>
            <Card title="生成配置" style={{ marginBottom: 16 }}>
              <Form form={form} layout="vertical" onFinish={handleGenerate}>
                <Form.Item
                  name="generator"
                  label="选择AI生成器"
                  rules={[{ required: true, message: '请选择生成器' }]}
                >
                  <Select
                    placeholder="选择AI模型"
                    value={selectedGenerator}
                    onChange={setSelectedGenerator}
                    options={generators.map((g) => ({
                      label: `${g.name} - ${g.description}`,
                      value: g.id,
                    }))}
                  />
                </Form.Item>

                <Form.Item
                  name="template"
                  label="选择提示词模板"
                  rules={[{ required: true, message: '请选择模板' }]}
                >
                  <Select
                    placeholder="选择提示词模板"
                    value={selectedPrompt}
                    onChange={setSelectedPrompt}
                    options={promptTemplates.map((t) => ({
                      label: `${t.name} (${t.category})`,
                      value: t.id.toString(),
                    }))}
                  />
                </Form.Item>

                <Form.Item
                  name="title"
                  label="内容标题"
                  rules={[{ required: true, message: '请输入标题' }]}
                >
                  <Input placeholder="输入要生成的内容标题" />
                </Form.Item>

                <Form.Item
                  name="prompt"
                  label="额外提示词（可选）"
                >
                  <Input.TextArea
                    rows={5}
                    placeholder="输入额外的提示词来定制生成内容..."
                  />
                </Form.Item>

                <Form.Item>
                  <Space>
                    <Button
                      type="primary"
                      icon={<SendOutlined />}
                      loading={generating}
                      onClick={() => form.submit()}
                    >
                      生成内容
                    </Button>
                    <Button onClick={() => setShowHistoryModal(true)}>
                      查看历史
                    </Button>
                  </Space>
                </Form.Item>
              </Form>
            </Card>
          </Col>

          {/* 右侧预览区 */}
          <Col xs={24} lg={12}>
            <Card
              title="生成结果预览"
              extra={
                generatedContent && (
                  <Button
                    type="primary"
                    icon={<SaveOutlined />}
                    onClick={handleSavedAsDraft}
                  >
                    保存为草稿
                  </Button>
                )
              }
            >
              {generatedContent ? (
                <div
                  style={{
                    background: '#fafafa',
                    padding: 16,
                    borderRadius: 4,
                    maxHeight: 500,
                    overflowY: 'auto',
                    whiteSpace: 'pre-wrap',
                    fontFamily: 'monospace',
                    fontSize: 12,
                  }}
                >
                  {generatedContent}
                </div>
              ) : (
                <div
                  style={{
                    background: '#fafafa',
                    padding: 32,
                    borderRadius: 4,
                    textAlign: 'center',
                    color: '#999',
                  }}
                >
                  生成的内容将在此显示...
                </div>
              )}
            </Card>
          </Col>
        </Row>

        {/* 生成历史模态框 */}
        <Modal
          title="生成历史"
          open={showHistoryModal}
          onCancel={() => setShowHistoryModal(false)}
          width={1000}
          footer={null}
        >
          <Table
            columns={historyColumns}
            dataSource={generationHistory}
            rowKey="id"
            pagination={{ pageSize: 10 }}
            size="small"
          />
        </Modal>
      </div>
    </Spin>
  )
}
