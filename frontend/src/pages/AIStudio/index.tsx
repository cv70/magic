import { useEffect, useState } from 'react'
import { Card, Form, Button, Select, Space, message, Spin, Row, Col, Modal, Table, Steps } from 'antd'
import { SaveOutlined } from '@ant-design/icons'
import { aiApi, draftApi } from '../../utils/api'
import { getDomainById } from '../../constants/domains'
import { getTemplatesByDomain, buildPromptFromTemplate } from '../../constants/domainTemplates'
import DomainSelector from '../../components/DomainSelector'
import DomainParamsForm from '../../components/DomainParamsForm'
import type { DomainTemplate } from '../../constants/domainTemplates'

export function AIStudio() {
  const [loading, setLoading] = useState(true)
  const [generating, setGenerating] = useState(false)
  const [generatedContent, setGeneratedContent] = useState('')
  const [selectedGenerator, setSelectedGenerator] = useState('')
  const [generationHistory, setGenerationHistory] = useState<any[]>([])
  const [showHistoryModal, setShowHistoryModal] = useState(false)
  const [generators, setGenerators] = useState<any[]>([])

  // 向导式流程相关状态
  const [currentStep, setCurrentStep] = useState(0)
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null)
  const [domainTemplates, setDomainTemplates] = useState<DomainTemplate[]>([])
  const [selectedTemplate, setSelectedTemplate] = useState<DomainTemplate | null>(null)
  const [templateParams, setTemplateParams] = useState<Record<string, any>>({})

  useEffect(() => {
    loadGenerators()
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

  // Step 1: 选择领域后，加载该领域的模板
  const handleDomainSelect = (domainId: string) => {
    setSelectedDomain(domainId)
    const templates = getTemplatesByDomain(domainId)
    setDomainTemplates(templates)
    setSelectedTemplate(templates.length > 0 ? templates[0] : null)
    setCurrentStep(1)
  }

  // Step 2: 选择模板
  const handleTemplateSelect = (templateId: string) => {
    const template = domainTemplates.find(t => t.id === templateId)
    setSelectedTemplate(template || null)
    setCurrentStep(2)
  }

  // Step 3: 填写参数后生成
  const handleGenerateFromTemplate = async (params: Record<string, any>) => {
    if (!selectedGenerator) {
      message.error('请先选择 AI 生成器')
      return
    }
    if (!selectedTemplate) {
      message.error('请先选择模板')
      return
    }

    setGenerating(true)
    try {
      // 根据模板和参数构建最终的提示词
      const finalPrompt = buildPromptFromTemplate(selectedTemplate, params)

      // 调用 AI 生成 API
      const response = await aiApi.generate({
        generatorID: parseInt(selectedGenerator),
        input: finalPrompt,
        title: params.title || selectedTemplate.name,
      })

      // 轮询获取生成结果
      const content = await pollGenerationStatus(response.id)
      setGeneratedContent(content)
      setTemplateParams(params)

      // Save to history
      const newItem = {
        id: Date.now(),
        title: params.title || selectedTemplate.name,
        domain: selectedDomain,
        template: selectedTemplate.id,
        prompt: finalPrompt,
        content,
        createdAt: new Date().toISOString(),
      }

      const updatedHistory = [newItem, ...generationHistory].slice(0, 50)
      setGenerationHistory(updatedHistory)
      localStorage.setItem('ai_generation_history', JSON.stringify(updatedHistory))

      setCurrentStep(3)
      message.success('内容已生成！')
    } catch (error) {
      console.error('生成失败:', error)
      message.error(error instanceof Error ? error.message : '生成失败，请重试')
    } finally {
      setGenerating(false)
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

  const handleSavedAsDraft = async () => {
    if (!generatedContent) {
      message.error('没有生成的内容可保存')
      return
    }

    try {
      setLoading(true)
      const title = selectedTemplate?.name || templateParams.title || '生成的内容'

      // Save as draft via API with domain metadata
      await draftApi.create({
        title,
        content: generatedContent,
        content_type: 'html',
        metadata: {
          domain: selectedDomain,
          template_id: selectedTemplate?.id,
          ...templateParams
        }
      })
      message.success(`已保存为草稿：${title}`)

      // Reset form
      setGeneratedContent('')
      setSelectedDomain(null)
      setSelectedTemplate(null)
      setTemplateParams({})
      setCurrentStep(0)
    } catch (error) {
      message.error('保存失败')
    } finally {
      setLoading(false)
    }
  }

  const historyColumns = [
    { title: '标题', dataIndex: 'title', key: 'title', width: 150 },
    {
      title: '领域',
      dataIndex: 'domain',
      key: 'domain',
      width: 120,
      render: (domain: string) => {
        const domainInfo = getDomainById(domain)
        return domainInfo ? `${domainInfo.icon} ${domainInfo.name}` : '-'
      }
    },
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
        <Button
          type="link"
          size="small"
          onClick={() => {
            setGeneratedContent(record.content)
            setSelectedDomain(record.domain)
            setTemplateParams(record.params || {})
            setCurrentStep(3)
            setShowHistoryModal(false)
          }}
        >
          恢复
        </Button>
      ),
    },
  ]

  const steps = [
    {
      title: '选择领域',
      description: '选择内容所属的垂直领域'
    },
    {
      title: '选择模板',
      description: '从领域专属模板中选择'
    },
    {
      title: '填写参数',
      description: '自定义模板参数'
    },
    {
      title: '查看结果',
      description: '预览和保存生成内容'
    }
  ]

  return (
    <Spin spinning={loading}>
      <div>
        <h1>✨ AI Studio - 领域专属内容生成</h1>
        <div style={{ marginBottom: 24, color: '#666' }}>
          按领域选择专用模板，自动填充参数，一键生成高质量内容
        </div>

        {/* 向导进度条 */}
        <Card style={{ marginBottom: 24 }}>
          <Steps current={currentStep} items={steps} />
        </Card>

        {/* Step 0: 选择领域 */}
        {currentStep === 0 && (
          <Card title="第 1 步：选择内容领域" style={{ marginBottom: 24 }}>
            <DomainSelector value={selectedDomain || undefined} onChange={handleDomainSelect} />
          </Card>
        )}

        {/* Step 1: 选择模板 */}
        {currentStep === 1 && (
          <Card
            title={`第 2 步：选择 ${selectedDomain ? getDomainById(selectedDomain)?.name : ''} 的专属模板`}
            extra={
              <Button type="link" onClick={() => setCurrentStep(0)}>
                重新选择领域
              </Button>
            }
            style={{ marginBottom: 24 }}
          >
            <Form layout="vertical">
              <Form.Item label="可用的模板">
                <Select
                  placeholder="选择模板"
                  value={selectedTemplate?.id}
                  onChange={handleTemplateSelect}
                  options={domainTemplates.map((t) => ({
                    label: t.name,
                    value: t.id,
                  }))}
                />
              </Form.Item>
              {selectedTemplate && (
                <div style={{ background: '#f5f5f5', padding: 12, borderRadius: 4 }}>
                  <div style={{ marginBottom: 8 }}>
                    <strong>{selectedTemplate.name}</strong>
                  </div>
                  <div style={{ fontSize: 12, color: '#666', marginBottom: 12 }}>
                    {selectedTemplate.description}
                  </div>
                  <div style={{ fontSize: 12, color: '#999' }}>
                    示例输出：{selectedTemplate.exampleOutput}
                  </div>
                </div>
              )}
            </Form>
          </Card>
        )}

        {/* Step 2: 填写参数 */}
        {currentStep === 2 && selectedTemplate && (
          <Card
            title={`第 3 步：填写 "${selectedTemplate.name}" 的参数`}
            extra={
              <Button type="link" onClick={() => setCurrentStep(1)}>
                重新选择模板
              </Button>
            }
            style={{ marginBottom: 24 }}
          >
            <Row gutter={16}>
              <Col xs={24} md={12}>
                <div style={{ marginBottom: 16 }}>
                  <label style={{ fontWeight: 600, display: 'block', marginBottom: 8 }}>
                    选择 AI 生成器
                  </label>
                  <Select
                    placeholder="选择 AI 模型"
                    value={selectedGenerator}
                    onChange={setSelectedGenerator}
                    options={generators.map((g) => ({
                      label: `${g.name} - ${g.description}`,
                      value: g.id,
                    }))}
                  />
                </div>
              </Col>
            </Row>

            <DomainParamsForm
              template={selectedTemplate}
              onSubmit={handleGenerateFromTemplate}
              loading={generating}
            />
          </Card>
        )}

        {/* Step 3: 预览和保存 */}
        {currentStep === 3 && (
          <Card
            title="第 4 步：预览和保存"
            extra={
              <Space>
                <Button onClick={() => setCurrentStep(2)}>返回编辑</Button>
                <Button
                  type="primary"
                  icon={<SaveOutlined />}
                  onClick={handleSavedAsDraft}
                  loading={loading}
                >
                  保存为草稿
                </Button>
              </Space>
            }
            style={{ marginBottom: 24 }}
          >
            <div style={{ marginBottom: 16 }}>
              <div style={{ marginBottom: 8 }}>
                <span style={{ color: '#666' }}>
                  领域：{selectedDomain ? getDomainById(selectedDomain)?.name : '-'}
                </span>
              </div>
              <div>
                <span style={{ color: '#666' }}>
                  模板：{selectedTemplate?.name}
                </span>
              </div>
            </div>

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
                正在生成内容...
              </div>
            )}

            <Space style={{ marginTop: 16 }}>
              <Button onClick={() => setShowHistoryModal(true)}>查看历史</Button>
            </Space>
          </Card>
        )}

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
