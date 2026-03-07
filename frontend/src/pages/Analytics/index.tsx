import { useEffect, useState } from 'react'
import { Row, Col, Card, Select, Spin, message, Statistic } from 'antd'
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { analyticsApi } from '../../utils/api'

const COLORS = ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1']

export function Analytics() {
  const [days, setDays] = useState(30)
  const [loading, setLoading] = useState(true)
  const [summary, setSummary] = useState<any>(null)
  const [ranking, setRanking] = useState<any[]>([])
  const [platformComparison, setPlatformComparison] = useState<any[]>([])
  const [trendData, setTrendData] = useState<any[]>([])

  useEffect(() => {
    loadAnalytics()
  }, [days])

  const loadAnalytics = async () => {
    try {
      setLoading(true)

      // 加载统计摘要
      const summaryRes = await analyticsApi.summary({ days })
      setSummary(summaryRes)

      // 加载排行数据
      const rankingRes = await analyticsApi.ranking({ days, limit: 10 })
      setRanking((rankingRes.data as any) || [])

      // 加载平台对比数据
      const platformRes = await analyticsApi.platformComparison()
      setPlatformComparison((platformRes.data as any) || [])

      // 生成趋势数据（模拟：实际应从后端获取）
      generateTrendData()
    } catch (error) {
      message.error('加载分析数据失败')
    } finally {
      setLoading(false)
    }
  }

  // 生成模拟的30天趋势数据
  const generateTrendData = () => {
    const data = []
    for (let i = 29; i >= 0; i--) {
      const date = new Date()
      date.setDate(date.getDate() - i)
      data.push({
        date: `${date.getMonth() + 1}/${date.getDate()}`,
        published: Math.floor(Math.random() * 10) + 1,
        views: Math.floor(Math.random() * 1000) + 100,
      })
    }
    setTrendData(data)
  }

  return (
    <Spin spinning={loading}>
      <div>
        <div style={{ marginBottom: 24 }}>
          <h1>📊 数据分析</h1>
          <Select
            value={days}
            onChange={setDays}
            style={{ width: 150 }}
            options={[
              { label: '最近7天', value: 7 },
              { label: '最近30天', value: 30 },
              { label: '最近90天', value: 90 },
            ]}
          />
        </div>

        {/* 统计卡片 */}
        {summary && (
          <Row gutter={16} style={{ marginBottom: 24 }}>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Statistic title="总发布数" value={summary.total_published || 0} />
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Statistic title="总浏览数" value={summary.total_views || 0} />
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Statistic title="总点赞数" value={summary.total_likes || 0} />
              </Card>
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Card>
                <Statistic title="总评论数" value={summary.total_comments || 0} />
              </Card>
            </Col>
          </Row>
        )}

        {/* 发布趋势图表 */}
        <Row gutter={16} style={{ marginBottom: 24 }}>
          <Col xs={24}>
            <Card title="发布趋势（最近30天）">
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={trendData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Line type="monotone" dataKey="published" stroke="#1890ff" name="发布数" />
                  <Line type="monotone" dataKey="views" stroke="#52c41a" name="浏览数" yAxisId="right" />
                </LineChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>

        {/* 平台对比图表 */}
        {platformComparison.length > 0 && (
          <Row gutter={16} style={{ marginBottom: 24 }}>
            <Col xs={24} lg={12}>
              <Card title="平台发布数对比">
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={platformComparison}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="platform" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="count" fill="#1890ff" />
                  </BarChart>
                </ResponsiveContainer>
              </Card>
            </Col>

            {/* 平台成功率饼图 */}
            <Col xs={24} lg={12}>
              <Card title="平台成功率分布">
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={platformComparison}
                      dataKey="success_rate"
                      nameKey="platform"
                      cx="50%"
                      cy="50%"
                      outerRadius={100}
                      label
                    >
                      {platformComparison.map((_: any, index: number) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip formatter={(value: any) => {
                      const numValue = typeof value === 'number' ? value : 0
                      return `${numValue.toFixed(1)}%`
                    }} />
                  </PieChart>
                </ResponsiveContainer>
              </Card>
            </Col>
          </Row>
        )}

        {/* 内容排行榜 */}
        {ranking.length > 0 && (
          <Row gutter={16}>
            <Col xs={24}>
              <Card title="热门内容TOP 10">
                <table
                  style={{
                    width: '100%',
                    borderCollapse: 'collapse',
                  }}
                >
                  <thead>
                    <tr style={{ borderBottom: '1px solid #f0f0f0' }}>
                      <th style={{ padding: 8, textAlign: 'left' }}>排名</th>
                      <th style={{ padding: 8, textAlign: 'left' }}>内容ID</th>
                      <th style={{ padding: 8, textAlign: 'left' }}>发布次数</th>
                      <th style={{ padding: 8, textAlign: 'left' }}>总浏览</th>
                      <th style={{ padding: 8, textAlign: 'left' }}>总点赞</th>
                    </tr>
                  </thead>
                  <tbody>
                    {ranking.map((item: any, index) => (
                      <tr key={index} style={{ borderBottom: '1px solid #f0f0f0' }}>
                        <td style={{ padding: 8 }}>#{index + 1}</td>
                        <td style={{ padding: 8 }}>{item.content_id}</td>
                        <td style={{ padding: 8 }}>{item.publish_count || 0}</td>
                        <td style={{ padding: 8 }}>{item.total_views || 0}</td>
                        <td style={{ padding: 8 }}>{item.total_likes || 0}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </Card>
            </Col>
          </Row>
        )}
      </div>
    </Spin>
  )
}
