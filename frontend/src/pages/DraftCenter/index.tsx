import { useEffect, useState } from 'react'
import { draftApi } from '../../utils/api'
import type { Draft } from '../../types'
import './DraftCenter.css'

interface DraftItemProps {
  draft: Draft
  onEdit: (draft: Draft) => void
  onDelete: (id: number) => void
  onPublish: (id: number) => void
}

function DraftItem({ draft, onEdit, onDelete, onPublish }: DraftItemProps) {
  const lastEdited = draft.last_edited_at ? new Date(draft.last_edited_at).toLocaleString() : 'N/A'
  const preview = draft.content.replace(/<[^>]*>/g, '').substring(0, 100) + '...'

  return (
    <div className="draft-item">
      <div className="draft-content">
        <h3>{draft.title}</h3>
        <p className="preview">{preview}</p>
        <div className="meta">
          <span className="type">{draft.content_type}</span>
          <span className="date">上次编辑: {lastEdited}</span>
        </div>
      </div>
      <div className="draft-actions">
        <button className="btn btn-primary" onClick={() => onEdit(draft)}>
          编辑
        </button>
        <button className="btn btn-success" onClick={() => onPublish(draft.id)}>
          发布
        </button>
        <button className="btn btn-danger" onClick={() => onDelete(draft.id)}>
          删除
        </button>
      </div>
    </div>
  )
}

export function DraftCenter() {
  const [drafts, setDrafts] = useState<Draft[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [keyword, setKeyword] = useState('')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [showNewForm, setShowNewForm] = useState(false)
  const [newTitle, setNewTitle] = useState('')

  const limit = 20

  // 加载草稿列表
  const loadDrafts = async (searchPage = 1, searchKeyword = '') => {
    setLoading(true)
    setError(null)

    try {
      const result = await draftApi.list(searchPage, limit, searchKeyword)
      setDrafts(result.data)
      setTotal(result.total)
      setPage(searchPage)
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载失败')
    } finally {
      setLoading(false)
    }
  }

  // 初始化加载
  useEffect(() => {
    loadDrafts(1, keyword)
  }, [])

  // 搜索
  const handleSearch = (value: string) => {
    setKeyword(value)
    loadDrafts(1, value)
  }

  // 创建新草稿
  const handleCreateDraft = async () => {
    if (!newTitle.trim()) {
      setError('标题不能为空')
      return
    }

    setLoading(true)
    setError(null)

    try {
      await draftApi.create({
        title: newTitle,
        content: '',
        content_type: 'text',
      })
      setNewTitle('')
      setShowNewForm(false)
      await loadDrafts(1, keyword)
    } catch (err) {
      setError(err instanceof Error ? err.message : '创建失败')
    } finally {
      setLoading(false)
    }
  }

  // 编辑草稿
  const handleEditDraft = (draft: Draft) => {
    // TODO: 导航到编辑页面
    console.log('Edit draft:', draft.id)
    alert(`编辑草稿 ${draft.title} - 功能在编辑页面中实现`)
  }

  // 删除草稿
  const handleDeleteDraft = async (id: number) => {
    if (!confirm('确定要删除这篇草稿吗?')) return

    setLoading(true)
    setError(null)

    try {
      await draftApi.delete(id)
      await loadDrafts(page, keyword)
    } catch (err) {
      setError(err instanceof Error ? err.message : '删除失败')
    } finally {
      setLoading(false)
    }
  }

  // 发布草稿
  const handlePublishDraft = async (draftId: number) => {
    if (!confirm('确定要发布这篇草稿吗?')) return

    setLoading(true)
    setError(null)

    try {
      await draftApi.publish(draftId, 'published')
      await loadDrafts(page, keyword)
      alert('发布成功!')
    } catch (err) {
      setError(err instanceof Error ? err.message : '发布失败')
    } finally {
      setLoading(false)
    }
  }

  // 分页
  const handlePreviousPage = () => {
    if (page > 1) {
      loadDrafts(page - 1, keyword)
    }
  }

  const handleNextPage = () => {
    const totalPages = Math.ceil(total / limit)
    if (page < totalPages) {
      loadDrafts(page + 1, keyword)
    }
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <div className="draft-center">
      <div className="header">
        <h1>📝 草稿中心</h1>
        <div className="toolbar">
          <input
            type="text"
            placeholder="搜索草稿..."
            value={keyword}
            onChange={(e) => handleSearch(e.target.value)}
            className="search-input"
          />
          <button className="btn btn-primary" onClick={() => setShowNewForm(true)}>
            + 新建草稿
          </button>
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      {showNewForm && (
        <div className="new-draft-form">
          <input
            type="text"
            placeholder="输入标题..."
            value={newTitle}
            onChange={(e) => setNewTitle(e.target.value)}
            className="input-field"
          />
          <div className="form-actions">
            <button className="btn btn-primary" onClick={handleCreateDraft} disabled={loading}>
              {loading ? '创建中...' : '创建'}
            </button>
            <button className="btn btn-secondary" onClick={() => setShowNewForm(false)}>
              取消
            </button>
          </div>
        </div>
      )}

      {loading && <div className="loading">加载中...</div>}

      {drafts.length === 0 && !loading && (
        <div className="empty-state">
          <p>暂无草稿</p>
          <p>点击"新建草稿"开始创建</p>
        </div>
      )}

      <div className="drafts-list">
        {drafts.map((draft) => (
          <DraftItem
            key={draft.id}
            draft={draft}
            onEdit={handleEditDraft}
            onDelete={handleDeleteDraft}
            onPublish={handlePublishDraft}
          />
        ))}
      </div>

      {totalPages > 1 && (
        <div className="pagination">
          <button className="btn" onClick={handlePreviousPage} disabled={page === 1}>
            ← 上一页
          </button>
          <span className="page-info">
            第 {page} / {totalPages} 页
          </span>
          <button className="btn" onClick={handleNextPage} disabled={page === totalPages}>
            下一页 →
          </button>
        </div>
      )}

      <div className="stats">
        <span>总共 {total} 篇草稿</span>
      </div>
    </div>
  )
}

export default DraftCenter
