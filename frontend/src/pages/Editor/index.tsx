import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useDraft } from '../../hooks/useDraft'
import type { Draft } from '../../types'
import './Editor.css'

interface EditorState {
  title: string
  content: string
  contentType: string
  tags: string[]
}

/**
 * Editor Page Component
 * Full-featured editor for creating and editing drafts
 *
 * Features:
 * - Real-time editing with auto-save
 * - Version history and recovery
 * - Multiple content types (text, code, image, etc.)
 * - Tag management
 * - Auto-save status indicator
 * - Publish to content
 */
export function Editor() {
  const { draftId } = useParams<{ draftId?: string }>()
  const navigate = useNavigate()

  const [editorState, setEditorState] = useState<EditorState>({
    title: '',
    content: '',
    contentType: 'text',
    tags: [],
  })

  const [tagInput, setTagInput] = useState('')
  const [showVersionHistory, setShowVersionHistory] = useState(false)
  const [selectedVersion, setSelectedVersion] = useState<number | null>(null)

  const draftIdNum = draftId ? parseInt(draftId) : undefined

  // Use the draft hook for state management and auto-save
  const {
    draft,
    isDirty,
    isSaving,
    lastSavedAt,
    updateDraft,
    save,
    getVersions,
    revertVersion,
  } = useDraft(draftIdNum || 0, {
    autoSaveInterval: 30000, // Auto-save every 30 seconds
    onSaved: (version) => {
      console.log('Draft saved, version:', version)
    },
    onError: (err) => {
      console.error('Auto-save error:', err)
    },
  })

  // Initialize editor state from draft
  useEffect(() => {
    if (draft) {
      setEditorState({
        title: draft.title,
        content: draft.content,
        contentType: draft.content_type,
        tags: draft.tags || [],
      })
    }
  }, [draft])

  // Handle title change
  const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newTitle = e.target.value
    setEditorState((prev) => ({ ...prev, title: newTitle }))
    updateDraft({ title: newTitle })
  }

  // Handle content change
  const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value
    setEditorState((prev) => ({ ...prev, content: newContent }))
    updateDraft({ content: newContent })
  }

  // Handle content type change
  const handleContentTypeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newType = e.target.value
    setEditorState((prev) => ({ ...prev, contentType: newType }))
    updateDraft({ content_type: newType })
  }

  // Handle tag input
  const handleAddTag = () => {
    if (tagInput.trim() && !editorState.tags.includes(tagInput.trim())) {
      const newTags = [...editorState.tags, tagInput.trim()]
      setEditorState((prev) => ({ ...prev, tags: newTags }))
      updateDraft({ tags: newTags })
      setTagInput('')
    }
  }

  // Handle tag removal
  const handleRemoveTag = (tag: string) => {
    const newTags = editorState.tags.filter((t) => t !== tag)
    setEditorState((prev) => ({ ...prev, tags: newTags }))
    updateDraft({ tags: newTags })
  }

  // Handle manual save
  const handleSave = async () => {
    try {
      await save()
    } catch (err) {
      console.error('Failed to save:', err)
    }
  }

  // Handle publish
  const handlePublish = async () => {
    if (!draftIdNum) {
      alert('Please save the draft first')
      return
    }

    if (!confirm('Are you sure you want to publish this draft? It cannot be edited after publishing.')) {
      return
    }

    try {
      // TODO: Call publish API
      console.log('Publishing draft:', draftIdNum)
      alert('Draft published successfully!')
      navigate('/drafts')
    } catch (err) {
      alert('Failed to publish draft')
      console.error(err)
    }
  }

  // Handle version history
  const handleShowVersions = async () => {
    try {
      const versions = await getVersions()
      console.log('Versions:', versions)
      setShowVersionHistory(true)
    } catch (err) {
      console.error('Failed to fetch versions:', err)
    }
  }

  // Handle revert to version
  const handleRevertVersion = async (versionNum: number) => {
    if (!confirm(`Are you sure you want to revert to version ${versionNum}?`)) {
      return
    }

    try {
      await revertVersion(versionNum)
      setShowVersionHistory(false)
      setSelectedVersion(null)
    } catch (err) {
      console.error('Failed to revert version:', err)
    }
  }

  if (!draft && draftId) {
    return <div className="editor-loading">Loading draft...</div>
  }

  return (
    <div className="editor-container">
      {/* Header */}
      <div className="editor-header">
        <div className="editor-title-section">
          <input
            type="text"
            value={editorState.title}
            onChange={handleTitleChange}
            placeholder="Enter draft title..."
            className="editor-title-input"
          />
          <span className={`editor-status ${isDirty ? 'dirty' : 'saved'}`}>
            {isSaving ? '保存中...' : isDirty ? '未保存' : '已保存'}
          </span>
          {lastSavedAt && (
            <span className="editor-last-saved">
              最后保存于 {lastSavedAt.toLocaleTimeString()}
            </span>
          )}
        </div>

        <div className="editor-actions">
          <select value={editorState.contentType} onChange={handleContentTypeChange} className="editor-type-select">
            <option value="text">文本</option>
            <option value="code">代码</option>
            <option value="markdown">Markdown</option>
            <option value="html">HTML</option>
          </select>

          <button onClick={handleSave} disabled={isSaving || !isDirty} className="btn btn-primary">
            {isSaving ? '保存中...' : '保存'}
          </button>

          <button onClick={handleShowVersions} className="btn btn-secondary">
            版本历史
          </button>

          <button onClick={handlePublish} className="btn btn-success">
            发布
          </button>

          <button onClick={() => navigate('/drafts')} className="btn btn-secondary">
            返回
          </button>
        </div>
      </div>

      {/* Main editor */}
      <div className="editor-main">
        {/* Editor area */}
        <div className="editor-area">
          <textarea
            value={editorState.content}
            onChange={handleContentChange}
            placeholder="输入你的内容..."
            className="editor-textarea"
          />
        </div>

        {/* Sidebar */}
        <div className="editor-sidebar">
          {/* Tags section */}
          <div className="editor-section">
            <h3>标签</h3>
            <div className="tag-input-group">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleAddTag()}
                placeholder="添加标签..."
                className="tag-input"
              />
              <button onClick={handleAddTag} className="btn btn-small">
                添加
              </button>
            </div>

            <div className="tags-list">
              {editorState.tags.map((tag) => (
                <span key={tag} className="tag">
                  {tag}
                  <button
                    className="tag-remove"
                    onClick={() => handleRemoveTag(tag)}
                  >
                    ×
                  </button>
                </span>
              ))}
            </div>
          </div>

          {/* Info section */}
          <div className="editor-section">
            <h3>信息</h3>
            {draft && (
              <ul className="info-list">
                <li>
                  <label>创建于:</label>
                  <span>{new Date(draft.created_at).toLocaleString()}</span>
                </li>
                <li>
                  <label>更新于:</label>
                  <span>{new Date(draft.updated_at).toLocaleString()}</span>
                </li>
                <li>
                  <label>状态:</label>
                  <span>{draft.status}</span>
                </li>
                <li>
                  <label>字数:</label>
                  <span>{editorState.content.length}</span>
                </li>
              </ul>
            )}
          </div>
        </div>
      </div>

      {/* Version history modal */}
      {showVersionHistory && (
        <div className="editor-modal">
          <div className="modal-content">
            <div className="modal-header">
              <h2>版本历史</h2>
              <button
                className="modal-close"
                onClick={() => setShowVersionHistory(false)}
              >
                ×
              </button>
            </div>

            <div className="modal-body">
              {/* TODO: Render version list from getVersions() */}
              <p>加载版本历史中...</p>
            </div>

            <div className="modal-footer">
              <button
                onClick={() => setShowVersionHistory(false)}
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

export default Editor
