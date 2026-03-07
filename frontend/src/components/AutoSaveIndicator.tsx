import { Tag, Tooltip, Spin } from 'antd'
import { CheckCircleOutlined, ClockCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons'
import { formatDistanceToNow } from 'date-fns'
import { zhCN } from 'date-fns/locale'

interface AutoSaveIndicatorProps {
  isSaving: boolean
  lastSavedAt: Date | null
  hasPendingChanges: boolean
  error: string | null
}

export function AutoSaveIndicator({
  isSaving,
  lastSavedAt,
  hasPendingChanges,
  error,
}: AutoSaveIndicatorProps) {
  if (error) {
    return (
      <Tooltip title={error}>
        <Tag icon={<ExclamationCircleOutlined />} color="red">
          保存失败
        </Tag>
      </Tooltip>
    )
  }

  if (isSaving) {
    return (
      <Tag icon={<Spin size="small" />} color="processing">
        保存中...
      </Tag>
    )
  }

  if (hasPendingChanges) {
    return (
      <Tag icon={<ClockCircleOutlined />} color="warning">
        有未保存的改动
      </Tag>
    )
  }

  if (lastSavedAt) {
    const timeAgo = formatDistanceToNow(new Date(lastSavedAt), {
      locale: zhCN,
      addSuffix: true,
    })
    return (
      <Tooltip title={`最后保存于: ${new Date(lastSavedAt).toLocaleString('zh-CN')}`}>
        <Tag icon={<CheckCircleOutlined />} color="success">
          已保存 {timeAgo}
        </Tag>
      </Tooltip>
    )
  }

  return (
    <Tag icon={<CheckCircleOutlined />} color="default">
      就绪
    </Tag>
  )
}
