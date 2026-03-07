import React from 'react'
import { Card, Avatar, Button, Tag, Row, Col } from 'antd'
import { HeartOutlined, HeartFilled, SaveOutlined, MessageOutlined } from '@ant-design/icons'
import { getDomainById } from '../constants/domains'
import type { SquarePost } from '../types'

interface SquareCardProps {
  post: SquarePost
  onLike?: (postId: number, isLiked: boolean) => void
  onCollect?: (postId: number, isCollected: boolean) => void
  onComment?: (postId: number) => void
  onViewDetail?: (postId: number) => void
}

/**
 * 内容广场卡片组件
 * 展示单个广场内容的卡片，包含点赞、收藏、评论等操作
 */
export const SquareCard: React.FC<SquareCardProps> = ({
  post,
  onLike,
  onCollect,
  onComment,
  onViewDetail
}) => {
  const domainInfo = getDomainById(post.domain)

  return (
    <Card
      hoverable
      onClick={() => onViewDetail?.(post.id)}
      style={{ height: '100%', display: 'flex', flexDirection: 'column' }}
      cover={
        post.cover_image && (
          <div
            style={{
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              height: 160,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#fff',
              fontSize: 48
            }}
          >
            {domainInfo?.icon || '📄'}
          </div>
        )
      }
    >
      <div style={{ flex: 1 }}>
        {/* 标题 */}
        <div style={{ marginBottom: 8 }}>
          <h3 style={{ margin: 0, fontSize: 14, lineHeight: 1.5 }}>
            {post.title}
          </h3>
        </div>

        {/* 预览文本 */}
        <div
          style={{
            fontSize: 12,
            color: '#666',
            marginBottom: 12,
            display: '-webkit-box',
            WebkitLineClamp: 2,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
            minHeight: 32
          }}
        >
          {post.preview_text}
        </div>

        {/* 领域标签和用户信息 */}
        <div style={{ marginBottom: 12 }}>
          {domainInfo && (
            <Tag color={domainInfo.color} style={{ marginBottom: 8 }}>
              {domainInfo.icon} {domainInfo.name}
            </Tag>
          )}
          {post.tags && post.tags.length > 0 && (
            <div style={{ fontSize: 12 }}>
              {post.tags.slice(0, 2).map((tag, idx) => (
                <Tag key={idx} style={{ marginRight: 4 }}>
                  {tag}
                </Tag>
              ))}
              {post.tags.length > 2 && <span style={{ fontSize: 11, color: '#999' }}>+{post.tags.length - 2}</span>}
            </div>
          )}
        </div>

        {/* 用户信息 */}
        {post.user && (
          <div style={{ marginBottom: 12, display: 'flex', alignItems: 'center', gap: 8 }}>
            <Avatar size={24} src={post.user.avatar}>{post.user.username?.charAt(0)}</Avatar>
            <span style={{ fontSize: 12, color: '#666' }}>{post.user.username}</span>
            <span style={{ fontSize: 11, color: '#999', marginLeft: 'auto' }}>
              {new Date(post.created_at).toLocaleDateString()}
            </span>
          </div>
        )}
      </div>

      {/* 底部操作栏 */}
      <Row gutter={8} style={{ borderTop: '1px solid #f0f0f0', paddingTop: 12 }}>
        <Col span={8}>
          <Button
            type="text"
            size="small"
            icon={post.is_liked ? <HeartFilled style={{ color: '#ff4d4f' }} /> : <HeartOutlined />}
            onClick={(e) => {
              e.stopPropagation()
              onLike?.(post.id, !post.is_liked)
            }}
            style={{ width: '100%' }}
          >
            {post.likes_count}
          </Button>
        </Col>
        <Col span={8}>
          <Button
            type="text"
            size="small"
            icon={post.is_collected ? <SaveOutlined style={{ color: '#1890ff' }} /> : <SaveOutlined />}
            onClick={(e) => {
              e.stopPropagation()
              onCollect?.(post.id, !post.is_collected)
            }}
            style={{ width: '100%' }}
          >
            {post.collects_count}
          </Button>
        </Col>
        <Col span={8}>
          <Button
            type="text"
            size="small"
            icon={<MessageOutlined />}
            onClick={(e) => {
              e.stopPropagation()
              onComment?.(post.id)
            }}
            style={{ width: '100%' }}
          >
            {post.comments_count}
          </Button>
        </Col>
      </Row>
    </Card>
  )
}

export default SquareCard
