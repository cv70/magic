import { useEffect, useState } from 'react'
import { Card, Form, Input, Button, Select, Space, message, Spin, Row, Col, Modal, Table } from 'antd'
import { SendOutlined, SaveOutlined } from '@ant-design/icons'
import { aiApi, promptApi, draftApi } from '../../utils/api'

export function AIStudio() {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(true)
  const [generating, setGenerating] = useState(false)
  const [generatedContent, setGeneratedContent] = useState('')
  const [selectedGenerator, setSelectedGenerator] = useState('')
  const [selectedPrompt, setSelectedPrompt] = useState('')
  const [generationHistory, setGenerationHistory] = useState<any[]>([])
  const [showHistoryModal, setShowHistoryModal] = useState(false)
  const [generators, setGenerators] = useState<any[]>([])
  const [promptTemplates, setPromptTemplates] = useState<any[]>([])

  useEffect(() => {
    loadGenerators()
    loadPromptTemplates()
    // Load generation history from localStorage
    const history = localStorage.getItem('ai_generation_history')
    if (history) {
      setGenerationHistory(JSON.parse(history))
    }
  }, [])

  const loadGenerators = async () => {
    try {
      const result = await aiApi.listGenerators()
      setGenerators(result.data || [])
    } catch (error) {
      console.error('加载生成器失败')
      setGenerators([])
    } finally {
      setLoading(false)
    }
  }

  const loadPromptTemplates = async () => {
    try {
      const result = await promptApi.list()
      setPromptTemplates(result.data || [])
    } catch (error) {
      console.error('加载提示词模板失败')
    }
  }

  const pollGenerationStatus = async (taskId: number, maxAttempts = 120): Promise<string> => {
    let attempts = 0
    while (attempts < maxAttempts) {
      try {
        const task = await aiApi.getTask(taskId)
        if (task.status === 'completed' || task.status === 'success') {
          return task.content || ''
        }
        if (task.status === 'failed' || task.status === 'error') {
          throw new Error(task.error || '生成失败')
        }
        // 待处理，等待后重试
        await new Promise((resolve) => setTimeout(resolve, 1000))
        attempts++
      } catch (error) {
        if (attempts >= maxAttempts - 1) {
          throw error
        }
        await new Promise((resolve) => setTimeout(resolve, 1000))
        attempts++
      }
    }
    throw new Error('生成超时')
  }

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
      // 调用 AI 生成 API
      const response = await aiApi.generate({
        generatorID: parseInt(selectedGenerator),
        input: values.prompt || '请生成高质量的内容',
        title: values.title,
      })

      // 轮询获取生成结果
      const content = await pollGenerationStatus(response.id)
      setGeneratedContent(content)

      // Save to history
      const newItem = {
        id: Date.now(),
        title: values.title || '未命名内容',
        generator: selectedGenerator,
        template: selectedPrompt,
        prompt: values.prompt,
        content,
        createdAt: new Date().toISOString(),
      }

      const updatedHistory = [newItem, ...generationHistory].slice(0, 50)
      setGenerationHistory(updatedHistory)
      localStorage.setItem('ai_generation_history', JSON.stringify(updatedHistory))

      message.success('内容已生成！')
    } catch (error) {
      console.error('生成失败:', error)
      message.error(error instanceof Error ? error.message : '生成失败，请重试')
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

      // Save as draft via API
      await draftApi.create({
        title,
        content: generatedContent,
        content_type: 'html',
      })
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
