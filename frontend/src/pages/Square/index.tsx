import { useEffect, useState } from 'react'
import { Row, Col, Tabs, Input, Select, Button, Space, Card, Spin, Empty, Modal, Form, message } from 'antd'
import { SearchOutlined, ThunderboltOutlined, FireOutlined, ClockCircleOutlined } from '@ant-design/icons'
import { squareApi } from '../../utils/api'
import { getDomainsList } from '../../constants/domains'
import SquareCard from '../../components/SquareCard'
import type { SquarePost } from '../../types'

export function Square() {
  const [posts, setPosts] = useState<SquarePost[]>([])
  const [loading, setLoading] = useState(true)
  const [pagination, setPagination] = useState({ page: 1, pageSize: 20, total: 0 })
  const [selectedDomain, setSelectedDomain] = useState<string | undefined>()
  const [searchKeyword, setSearchKeyword] = useState('')
  const [sortOrder, setSortOrder] = useState<'newest' | 'hottest' | 'trending'>('newest')
  const [detailModal, setDetailModal] = useState(false)
  const [selectedPost, setSelectedPost] = useState<SquarePost | null>(null)
  const [commentForm] = Form.useForm()
  const [postingComment, setPostingComment] = useState(false)

  const domains = getDomainsList()

  useEffect(() => {
    loadPosts()
  }, [pagination.page, selectedDomain, searchKeyword, sortOrder])

  const loadPosts = async () => {
    try {
      setLoading(true)
      const result = await squareApi.list({
        domain: selectedDomain,
        keyword: searchKeyword,
        sort: sortOrder,
        page: pagination.page,
        page_size: pagination.pageSize,
      })
      setPosts(result.data || [])
      setPagination(prev => ({ ...prev, total: result.total }))
    } catch (error) {
      message.error('加载广场内容失败')
    } finally {
      setLoading(false)
    }
  }

  const handleLike = async (postId: number, isLiked: boolean) => {
    try {
      if (isLiked) {
        await squareApi.unlike(postId)
      } else {
        await squareApi.like(postId)
      }
      // 更新本地状态
      setPosts(posts.map(p => {
        if (p.id === postId) {
          return {
            ...p,
            is_liked: !isLiked,
            likes_count: isLiked ? p.likes_count - 1 : p.likes_count + 1
          }
        }
        return p
      }))
    } catch (error) {
      message.error('操作失败')
    }
  }

  const handleCollect = async (postId: number, isCollected: boolean) => {
    try {
      if (isCollected) {
        await squareApi.uncollect(postId)
      } else {
        await squareApi.collect(postId)
      }
      // 更新本地状态
      setPosts(posts.map(p => {
        if (p.id === postId) {
          return {
            ...p,
            is_collected: !isCollected,
            collects_count: isCollected ? p.collects_count - 1 : p.collects_count + 1
          }
        }
        return p
      }))
    } catch (error) {
      message.error('操作失败')
    }
  }

  const handleViewDetail = async (postId: number) => {
    try {
      const post = await squareApi.get(postId)
      setSelectedPost(post)
      setDetailModal(true)
    } catch (error) {
      message.error('加载详情失败')
    }
  }

  const handlePostComment = async (values: any) => {
    if (!selectedPost) return

    try {
      setPostingComment(true)
      await squareApi.comment(selectedPost.id, values.comment)
      message.success('评论已发表')
      commentForm.resetFields()
      // 重新加载详情
      const post = await squareApi.get(selectedPost.id)
      setSelectedPost(post)
    } catch (error) {
      message.error('评论失败')
    } finally {
      setPostingComment(false)
    }
  }

  const sortTabs = [
    {
      key: 'newest',
      label: <><ClockCircleOutlined /> 最新</>,
      sort: 'newest' as const
    },
    {
      key: 'hottest',
      label: <><FireOutlined /> 最热</>,
      sort: 'hottest' as const
    },
    {
      key: 'trending',
      label: <><ThunderboltOutlined /> 趋势</>,
      sort: 'trending' as const
    }
  ]

  return (
    <div>
      <h1>🎨 内容广场</h1>
      <div style={{ marginBottom: 24, color: '#666' }}>
        发现来自全球创作者的精美内容，点赞、收藏和评论你喜欢的作品
      </div>

      {/* 搜索和筛选栏 */}
      <Card style={{ marginBottom: 24 }}>
        <Space direction="vertical" style={{ width: '100%' }} size="middle">
          <Row gutter={16}>
            <Col xs={24} sm={12}>
              <Input
                placeholder="搜索广场内容..."
                prefix={<SearchOutlined />}
                value={searchKeyword}
                onChange={(e) => {
                  setSearchKeyword(e.target.value)
                  setPagination(prev => ({ ...prev, page: 1 }))
                }}
                style={{ width: '100%' }}
              />
            </Col>
            <Col xs={24} sm={12}>
              <Select
                placeholder="按领域筛选"
                allowClear
                style={{ width: '100%' }}
                value={selectedDomain}
                onChange={(value) => {
                  setSelectedDomain(value)
                  setPagination(prev => ({ ...prev, page: 1 }))
                }}
                options={[
                  { label: '所有领域', value: undefined },
                  ...domains.map((d) => ({ label: `${d.icon} ${d.name}`, value: d.id }))
                ]}
              />
            </Col>
          </Row>
        </Space>
      </Card>

      {/* 排序选项卡 */}
      <Card style={{ marginBottom: 24 }}>
        <Tabs
          activeKey={sortOrder}
          onChange={(key) => setSortOrder(key as any)}
          items={sortTabs.map(tab => ({
            key: tab.key,
            label: tab.label
          }))}
        />
      </Card>

      {/* 广场内容网格 */}
      <Spin spinning={loading}>
        {posts.length > 0 ? (
          <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
            {posts.map((post) => (
              <Col key={post.id} xs={24} sm={12} md={8} lg={6}>
                <SquareCard
                  post={post}
                  onLike={handleLike}
                  onCollect={handleCollect}
                  onViewDetail={handleViewDetail}
                />
              </Col>
            ))}
          </Row>
        ) : (
          <Empty description="暂无内容" style={{ marginTop: 48 }} />
        )}
      </Spin>

      {/* 分页 */}
      {posts.length > 0 && (
        <div style={{ textAlign: 'center', marginTop: 24 }}>
          <Space>
            <Button
              disabled={pagination.page === 1}
              onClick={() => setPagination(prev => ({ ...prev, page: prev.page - 1 }))}
            >
              上一页
            </Button>
            <span>
              第 {pagination.page} / {Math.ceil(pagination.total / pagination.pageSize)} 页
            </span>
            <Button
              disabled={pagination.page >= Math.ceil(pagination.total / pagination.pageSize)}
              onClick={() => setPagination(prev => ({ ...prev, page: prev.page + 1 }))}
            >
              下一页
            </Button>
          </Space>
        </div>
      )}

      {/* 内容详情模态框 */}
      <Modal
        title={selectedPost?.title}
        open={detailModal}
        onCancel={() => setDetailModal(false)}
        footer={null}
        width={800}
      >
        {selectedPost && (
          <div>
            {/* 内容预览 */}
            <div style={{ marginBottom: 24, padding: 16, background: '#f5f5f5', borderRadius: 4 }}>
              <p>{selectedPost.preview_text}</p>
            </div>

            {/* 互动信息 */}
            <Row gutter={16} style={{ marginBottom: 24 }}>
              <Col span={8}>
                <div style={{ textAlign: 'center' }}>
                  <div style={{ fontSize: 20, fontWeight: 600 }}>{selectedPost.likes_count}</div>
                  <div style={{ fontSize: 12, color: '#999' }}>赞</div>
                </div>
              </Col>
              <Col span={8}>
                <div style={{ textAlign: 'center' }}>
                  <div style={{ fontSize: 20, fontWeight: 600 }}>{selectedPost.comments_count}</div>
                  <div style={{ fontSize: 12, color: '#999' }}>评论</div>
                </div>
              </Col>
              <Col span={8}>
                <div style={{ textAlign: 'center' }}>
                  <div style={{ fontSize: 20, fontWeight: 600 }}>{selectedPost.collects_count}</div>
                  <div style={{ fontSize: 12, color: '#999' }}>收藏</div>
                </div>
              </Col>
            </Row>

            {/* 评论表单 */}
            <Form form={commentForm} onFinish={handlePostComment}>
              <Form.Item
                name="comment"
                rules={[{ required: true, message: '请输入评论' }]}
              >
                <Input.TextArea
                  placeholder="写下你的想法..."
                  rows={3}
                />
              </Form.Item>
              <Form.Item>
                <Button
                  type="primary"
                  htmlType="submit"
                  loading={postingComment}
                  style={{ width: '100%' }}
                >
                  发表评论
                </Button>
              </Form.Item>
            </Form>

            {/* 评论列表 */}
            {selectedPost.comments_count > 0 && (
              <div style={{ marginTop: 24, borderTop: '1px solid #f0f0f0', paddingTop: 16 }}>
                <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 12 }}>
                  评论（{selectedPost.comments_count}）
                </div>
                <div style={{ fontSize: 12, color: '#999', textAlign: 'center' }}>
                  加载评论功能开发中...
                </div>
              </div>
            )}
          </div>
        )}
      </Modal>
    </div>
  )
}

export default Square
