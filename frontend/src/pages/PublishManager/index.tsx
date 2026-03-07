import { useEffect, useState } from 'react'
import { usePublish } from '../../hooks/usePublish'
import type { PublishTask } from '../../types'
import './PublishManager.css'

interface PublishFilter {
  status: string
  platform: string
}

/**
 * PublishManager Page
 * Manage and track content publishing across platforms
 *
 * Features:
 * - View all publish tasks
 * - Filter by status and platform
 * - Create new publish tasks
 * - Schedule publishing
 * - Track publishing progress
 * - Monitor platform results
 */
export function PublishManager() {
  const [tasks, setTasks] = useState<PublishTask[]>([])
  const [filter, setFilter] = useState<PublishFilter>({
    status: 'all',
    platform: 'all',
  })
  const [isLoading, setIsLoading] = useState(false)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [selectedTask, setSelectedTask] = useState<PublishTask | null>(null)

  const {
    publishTask,
    isPublishing,
    listTasks,
    executeTask,
    cancelTask,
    deleteTask,
  } = usePublish({
    onError: (err) => {
      console.error('Publish error:', err)
    },
  })

  // Load tasks on mount
  useEffect(() => {
    loadTasks()
  }, [])

  const loadTasks = async () => {
    setIsLoading(true)
    try {
      const result = await listTasks(1, 50)
      setTasks(result.data || [])
    } catch (err) {
      console.error('Failed to load tasks:', err)
    } finally {
      setIsLoading(false)
    }
  }

  // Filter tasks based on selected filters
  const filteredTasks = tasks.filter((task) => {
    if (filter.status !== 'all' && task.status !== filter.status) {
      return false
    }
    if (filter.platform !== 'all') {
      // Check if platform is in target_platforms
      const platforms = typeof task.target_platforms === 'string'
        ? JSON.parse(task.target_platforms)
        : task.target_platforms || []
      if (!platforms.includes(filter.platform)) {
        return false
      }
    }
    return true
  })

  // Get status badge class
  const getStatusClass = (status: string) => {
    switch (status) {
      case 'published':
        return 'status-published'
      case 'publishing':
        return 'status-publishing'
      case 'pending':
        return 'status-pending'
      case 'failed':
        return 'status-failed'
      default:
        return ''
    }
  }

  // Get status label
  const getStatusLabel = (status: string) => {
    const labels: Record<string, string> = {
      published: '已发布',
      publishing: '发布中',
      pending: '待发布',
      failed: '失败',
    }
    return labels[status] || status
  }

  // Handle task execution
  const handleExecuteTask = async (taskId: number) => {
    try {
      await executeTask(taskId)
      loadTasks()
    } catch (err) {
      console.error('Failed to execute task:', err)
    }
  }

  // Handle task cancellation
  const handleCancelTask = async (taskId: number) => {
    if (!confirm('Are you sure you want to cancel this task?')) {
      return
    }
    try {
      await cancelTask(taskId)
      loadTasks()
    } catch (err) {
      console.error('Failed to cancel task:', err)
    }
  }

  // Handle task deletion
  const handleDeleteTask = async (taskId: number) => {
    if (!confirm('Are you sure you want to delete this task?')) {
      return
    }
    try {
      await deleteTask(taskId)
      loadTasks()
    } catch (err) {
      console.error('Failed to delete task:', err)
    }
  }

  return (
    <div className="publish-manager">
      <div className="publish-header">
        <h1>📤 发布管理</h1>
        <button
          className="btn btn-primary"
          onClick={() => setShowCreateForm(true)}
        >
          + 新建发布任务
        </button>
      </div>

      {/* Filters */}
      <div className="publish-filters">
        <div className="filter-group">
          <label>状态</label>
          <select
            value={filter.status}
            onChange={(e) =>
              setFilter({ ...filter, status: e.target.value })
            }
            className="filter-select"
          >
            <option value="all">全部</option>
            <option value="pending">待发布</option>
            <option value="publishing">发布中</option>
            <option value="published">已发布</option>
            <option value="failed">失败</option>
          </select>
        </div>

        <div className="filter-group">
          <label>平台</label>
          <select
            value={filter.platform}
            onChange={(e) =>
              setFilter({ ...filter, platform: e.target.value })
            }
            className="filter-select"
          >
            <option value="all">全部</option>
            <option value="juejin">掘金</option>
            <option value="csdn">CSDN</option>
            <option value="segmentfault">SegmentFault</option>
            <option value="medium">Medium</option>
            <option value="dev.to">Dev.to</option>
          </select>
        </div>

        <button onClick={loadTasks} className="btn btn-secondary">
          刷新
        </button>
      </div>

      {/* Loading state */}
      {isLoading && <div className="loading">加载中...</div>}

      {/* Empty state */}
      {!isLoading && filteredTasks.length === 0 && (
        <div className="empty-state">
          <p>暂无发布任务</p>
          <button
            className="btn btn-primary"
            onClick={() => setShowCreateForm(true)}
          >
            创建新任务
          </button>
        </div>
      )}

      {/* Task list */}
      <div className="tasks-list">
        {filteredTasks.map((task) => (
          <div key={task.id} className="task-item">
            <div className="task-header">
              <div className="task-title">
                <h3>{task.id}</h3>
                <span className={`status-badge ${getStatusClass(task.status)}`}>
                  {getStatusLabel(task.status)}
                </span>
              </div>

              <div className="task-meta">
                <span className="meta-item">
                  创建: {new Date(task.created_at).toLocaleString()}
                </span>
                {task.scheduled_at && (
                  <span className="meta-item">
                    计划: {new Date(task.scheduled_at).toLocaleString()}
                  </span>
                )}
              </div>
            </div>

            <div className="task-body">
              <div className="task-platforms">
                {typeof task.target_platforms === 'string'
                  ? JSON.parse(task.target_platforms).map((p: string) => (
                      <span key={p} className="platform-tag">
                        {p}
                      </span>
                    ))
                  : (task.target_platforms || []).map((p: any) => (
                      <span key={p} className="platform-tag">
                        {p}
                      </span>
                    ))}
              </div>

              {task.status === 'failed' && task.error_message && (
                <div className="task-error">
                  <strong>错误:</strong> {task.error_message}
                </div>
              )}
            </div>

            <div className="task-actions">
              {task.status === 'pending' && (
                <>
                  <button
                    className="btn btn-success"
                    onClick={() => handleExecuteTask(task.id)}
                    disabled={isPublishing}
                  >
                    发布
                  </button>
                  <button
                    className="btn btn-secondary"
                    onClick={() => handleCancelTask(task.id)}
                  >
                    取消
                  </button>
                </>
              )}

              <button
                className="btn btn-danger"
                onClick={() => handleDeleteTask(task.id)}
              >
                删除
              </button>

              <button
                className="btn btn-primary"
                onClick={() => setSelectedTask(task)}
              >
                详情
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Task detail modal */}
      {selectedTask && (
        <div className="publish-modal">
          <div className="modal-content">
            <div className="modal-header">
              <h2>发布任务详情</h2>
              <button
                className="modal-close"
                onClick={() => setSelectedTask(null)}
              >
                ×
              </button>
            </div>

            <div className="modal-body">
              <div className="detail-section">
                <h3>基本信息</h3>
                <dl>
                  <dt>任务ID:</dt>
                  <dd>{selectedTask.id}</dd>
                  <dt>状态:</dt>
                  <dd>{getStatusLabel(selectedTask.status)}</dd>
                  <dt>创建时间:</dt>
                  <dd>{new Date(selectedTask.created_at).toLocaleString()}</dd>
                  {selectedTask.scheduled_at && (
                    <>
                      <dt>计划发布:</dt>
                      <dd>{new Date(selectedTask.scheduled_at).toLocaleString()}</dd>
                    </>
                  )}
                  {selectedTask.published_at && (
                    <>
                      <dt>发布时间:</dt>
                      <dd>{new Date(selectedTask.published_at).toLocaleString()}</dd>
                    </>
                  )}
                </dl>
              </div>

              <div className="detail-section">
                <h3>目标平台</h3>
                <div className="platforms-list">
                  {typeof selectedTask.target_platforms === 'string'
                    ? JSON.parse(selectedTask.target_platforms).map((p: string) => (
                        <span key={p} className="platform-tag">
                          {p}
                        </span>
                      ))
                    : (selectedTask.target_platforms || []).map((p: any) => (
                        <span key={p} className="platform-tag">
                          {p}
                        </span>
                      ))}
                </div>
              </div>

              {selectedTask.status === 'failed' && selectedTask.error_message && (
                <div className="detail-section">
                  <h3>错误信息</h3>
                  <p className="error-message">{selectedTask.error_message}</p>
                </div>
              )}

              {selectedTask.platform_results && (
                <div className="detail-section">
                  <h3>平台结果</h3>
                  <div className="results-list">
                    {typeof selectedTask.platform_results === 'string'
                      ? JSON.parse(selectedTask.platform_results).map((result: any, i: number) => (
                          <div key={i} className="result-item">
                            <h4>{result.platform}</h4>
                            <p>状态: {result.status}</p>
                            {result.post_url && (
                              <p>
                                <a href={result.post_url} target="_blank" rel="noopener noreferrer">
                                  查看文章
                                </a>
                              </p>
                            )}
                          </div>
                        ))
                      : null}
                  </div>
                </div>
              )}
            </div>

            <div className="modal-footer">
              <button
                onClick={() => setSelectedTask(null)}
                className="btn btn-secondary"
              >
                关闭
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Create task form */}
      {showCreateForm && (
        <div className="publish-modal">
          <div className="modal-content">
            <div className="modal-header">
              <h2>创建发布任务</h2>
              <button
                className="modal-close"
                onClick={() => setShowCreateForm(false)}
              >
                ×
              </button>
            </div>

            <div className="modal-body">
              <p>TODO: 实现发布任务创建表单</p>
            </div>

            <div className="modal-footer">
              <button
                onClick={() => setShowCreateForm(false)}
                className="btn btn-secondary"
              >
                关闭
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default PublishManager
