import { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Table, Empty, Spin, message } from 'antd'
import { FileTextOutlined, CheckCircleOutlined, SendOutlined, BgColorsOutlined } from '@ant-design/icons'
import { draftApi, contentApi, publishApi } from '../../utils/api'

interface DashboardStats {
  drafts: number
  published: number
  publishTasks: number
  aiGenerations: number
}

interface RecentItem {
  id: number
  title: string
  created_at?: string
  status?: string
}

export function Dashboard() {
  const [stats, setStats] = useState<DashboardStats>({
    drafts: 0,
    published: 0,
    publishTasks: 0,
    aiGenerations: 0,
  })
  const [recentDrafts, setRecentDrafts] = useState<RecentItem[]>([])
  const [recentTasks, setRecentTasks] = useState<RecentItem[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadDashboard()
  }, [])

  const loadDashboard = async () => {
    try {
      setLoading(true)
      const [draftsRes, contentRes, tasksRes] = await Promise.all([
        draftApi.list(1, 1000),
        contentApi.list(1, 1000),
        publishApi.list(1, 100),
      ])

      setStats({
        drafts: draftsRes.total,
        published: contentRes.total,
        publishTasks: tasksRes.total,
        aiGenerations: 0, // TODO: 从 AI domain 获取
      })

      setRecentDrafts((draftsRes.data || []).slice(0, 5))
      setRecentTasks((tasksRes.data as any || []).slice(0, 5))
    } catch (error) {
      message.error('加载数据失败')
    } finally {
      setLoading(false)
    }
  }

  const draftColumns = [
    { title: '标题', dataIndex: 'title', key: 'title' },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  ]

  const taskColumns = [
    { title: '内容ID', dataIndex: 'content_id', key: 'content_id' },
    { title: '状态', dataIndex: 'status', key: 'status' },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  ]

  return (
    <Spin spinning={loading}>
      <div>
        <h1>📊 数据总览</h1>

        <Row gutter={16} style={{ marginBottom: 24 }}>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="草稿数"
                value={stats.drafts}
                prefix={<FileTextOutlined />}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="已发布"
                value={stats.published}
                prefix={<CheckCircleOutlined />}
                valueStyle={{ color: '#52c41a' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="发布任务"
                value={stats.publishTasks}
                prefix={<SendOutlined />}
                valueStyle={{ color: '#faad14' }}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} lg={6}>
            <Card>
              <Statistic
                title="AI生成"
                value={stats.aiGenerations}
                prefix={<BgColorsOutlined />}
                valueStyle={{ color: '#722ed1' }}
              />
            </Card>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} lg={12}>
            <Card title="最近草稿" extra={<a href="/drafts">查看全部</a>}>
              {recentDrafts.length > 0 ? (
                <Table
                  columns={draftColumns}
                  dataSource={recentDrafts}
                  pagination={false}
                  size="small"
                  rowKey="id"
                />
              ) : (
                <Empty description="暂无草稿" />
              )}
            </Card>
          </Col>
          <Col xs={24} lg={12}>
            <Card title="最近发布任务" extra={<a href="/publish">查看全部</a>}>
              {recentTasks.length > 0 ? (
                <Table
                  columns={taskColumns}
                  dataSource={recentTasks}
                  pagination={false}
                  size="small"
                  rowKey="id"
                />
              ) : (
                <Empty description="暂无任务" />
              )}
            </Card>
          </Col>
        </Row>
      </div>
    </Spin>
  )
}
