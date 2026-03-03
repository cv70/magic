import { useEffect, useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import { api } from './api'
import type { Content, Generator, PublishTask, Publisher } from './api'
import './App.css'

function App() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const [contents, setContents] = useState<Content[]>([])
  const [generators, setGenerators] = useState<Generator[]>([])
  const [publishers, setPublishers] = useState<Publisher[]>([])
  const [publishTasks, setPublishTasks] = useState<PublishTask[]>([])

  const [contentForm, setContentForm] = useState({ title: '', content: '', contentType: 'text' })
  const [editingId, setEditingId] = useState<number | null>(null)

  const [generatorId, setGeneratorId] = useState<number | ''>('')
  const [prompt, setPrompt] = useState('')
  const [generatedTaskId, setGeneratedTaskId] = useState<number | null>(null)

  const [publishContentId, setPublishContentId] = useState<number | ''>('')
  const [publisherId, setPublisherId] = useState<number | ''>('')

  const isEditing = editingId !== null

  async function bootstrap() {
    setLoading(true)
    setError('')
    try {
      const [contentRes, generatorRes, publisherRes, taskRes] = await Promise.all([
        api.searchContents(),
        api.searchGenerators(),
        api.searchPublishers(),
        api.searchPublishTasks(),
      ])

      setContents(contentRes.contents)
      setGenerators(generatorRes.generators)
      setPublishers(publisherRes.publishers)
      setPublishTasks(taskRes.publish_tasks)

      if (generatorRes.generators.length > 0 && generatorId === '') {
        setGeneratorId(generatorRes.generators[0].id)
      }

      if (publisherRes.publishers.length > 0 && publisherId === '') {
        setPublisherId(publisherRes.publishers[0].id)
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : '加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void bootstrap()
  }, [])

  const contentOptions = useMemo(() => contents.map((c) => ({ value: c.id, label: `${c.id} - ${c.title}` })), [contents])

  async function handleSubmitContent(event: FormEvent) {
    event.preventDefault()
    if (!contentForm.title.trim() || !contentForm.content.trim()) {
      setError('标题和内容不能为空')
      return
    }

    setLoading(true)
    setError('')
    try {
      if (isEditing && editingId !== null) {
        await api.updateContent(editingId, contentForm.title, contentForm.content, contentForm.contentType)
      } else {
        await api.addContent(contentForm.title, contentForm.content, contentForm.contentType)
      }

      setContentForm({ title: '', content: '', contentType: 'text' })
      setEditingId(null)
      const contentRes = await api.searchContents()
      setContents(contentRes.contents)
    } catch (e) {
      setError(e instanceof Error ? e.message : '保存失败')
    } finally {
      setLoading(false)
    }
  }

  async function handleDeleteContent(id: number) {
    setLoading(true)
    setError('')
    try {
      await api.deleteContent(id)
      const contentRes = await api.searchContents()
      setContents(contentRes.contents)
      if (publishContentId === id) setPublishContentId('')
    } catch (e) {
      setError(e instanceof Error ? e.message : '删除失败')
    } finally {
      setLoading(false)
    }
  }

  function handleEditContent(content: Content) {
    setEditingId(content.id)
    setContentForm({
      title: content.title,
      content: content.content,
      contentType: content.content_type ?? 'text',
    })
  }

  async function handleGenerate(event: FormEvent) {
    event.preventDefault()
    if (generatorId === '' || !prompt.trim()) {
      setError('请选择生成器并输入提示词')
      return
    }

    setLoading(true)
    setError('')
    try {
      const res = await api.generate(generatorId, prompt)
      setGeneratedTaskId(res.id)
    } catch (e) {
      setError(e instanceof Error ? e.message : '生成失败')
    } finally {
      setLoading(false)
    }
  }

  async function handleSavePromptAsContent() {
    if (!prompt.trim()) {
      setError('提示词为空，无法保存')
      return
    }

    setLoading(true)
    setError('')
    try {
      await api.addContent(`AI Prompt ${new Date().toLocaleString()}`, prompt, 'text')
      const contentRes = await api.searchContents()
      setContents(contentRes.contents)
    } catch (e) {
      setError(e instanceof Error ? e.message : '保存为内容失败')
    } finally {
      setLoading(false)
    }
  }

  async function handlePublish(event: FormEvent) {
    event.preventDefault()
    if (publisherId === '' || publishContentId === '') {
      setError('请选择发布器和内容')
      return
    }

    setLoading(true)
    setError('')
    try {
      await api.publishContent(publisherId, publishContentId)
      const tasksRes = await api.searchPublishTasks()
      setPublishTasks(tasksRes.publish_tasks)
    } catch (e) {
      setError(e instanceof Error ? e.message : '发布失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="page">
      <header className="hero">
        <h1>AI Content Creation Engine</h1>
        <p>全链路 MVP 控制台：内容管理、AI 生成、发布任务</p>
        <div className="hero-actions">
          <button onClick={() => void bootstrap()} disabled={loading}>刷新数据</button>
          {loading && <span className="hint">处理中...</span>}
        </div>
        {error && <p className="error">{error}</p>}
      </header>

      <section className="grid">
        <article className="card">
          <h2>内容管理</h2>
          <form onSubmit={handleSubmitContent} className="form">
            <input
              placeholder="标题"
              value={contentForm.title}
              onChange={(e) => setContentForm((v) => ({ ...v, title: e.target.value }))}
            />
            <select
              value={contentForm.contentType}
              onChange={(e) => setContentForm((v) => ({ ...v, contentType: e.target.value }))}
            >
              <option value="text">text</option>
              <option value="code">code</option>
              <option value="image">image</option>
              <option value="audio">audio</option>
              <option value="video">video</option>
            </select>
            <textarea
              placeholder="内容"
              rows={6}
              value={contentForm.content}
              onChange={(e) => setContentForm((v) => ({ ...v, content: e.target.value }))}
            />
            <div className="row">
              <button type="submit" disabled={loading}>{isEditing ? '更新内容' : '新建内容'}</button>
              {isEditing && (
                <button
                  type="button"
                  onClick={() => {
                    setEditingId(null)
                    setContentForm({ title: '', content: '', contentType: 'text' })
                  }}
                >
                  取消编辑
                </button>
              )}
            </div>
          </form>

          <ul className="list">
            {contents.map((content) => (
              <li key={content.id}>
                <div>
                  <strong>{content.title}</strong>
                  <p>{content.content.slice(0, 80)}</p>
                </div>
                <div className="row">
                  <button onClick={() => handleEditContent(content)}>编辑</button>
                  <button onClick={() => void handleDeleteContent(content.id)}>删除</button>
                </div>
              </li>
            ))}
          </ul>
        </article>

        <article className="card">
          <h2>AI 生成</h2>
          <form onSubmit={handleGenerate} className="form">
            <select
              value={generatorId}
              onChange={(e) => setGeneratorId(e.target.value ? Number(e.target.value) : '')}
            >
              <option value="">选择生成器</option>
              {generators.map((generator) => (
                <option key={generator.id} value={generator.id}>
                  {generator.name} ({generator.provider}/{generator.model})
                </option>
              ))}
            </select>
            <textarea
              placeholder="输入提示词"
              rows={8}
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
            />
            <div className="row">
              <button type="submit" disabled={loading}>发起生成</button>
              <button type="button" onClick={() => void handleSavePromptAsContent()} disabled={loading}>
                保存提示词为内容
              </button>
            </div>
          </form>
          <p className="hint">最近生成任务 ID: {generatedTaskId ?? '-'}</p>
        </article>

        <article className="card">
          <h2>发布管理</h2>
          <form onSubmit={handlePublish} className="form">
            <select
              value={publisherId}
              onChange={(e) => setPublisherId(e.target.value ? Number(e.target.value) : '')}
            >
              <option value="">选择发布器</option>
              {publishers.map((publisher) => (
                <option key={publisher.id} value={publisher.id}>
                  {publisher.name} ({publisher.platform})
                </option>
              ))}
            </select>

            <select
              value={publishContentId}
              onChange={(e) => setPublishContentId(e.target.value ? Number(e.target.value) : '')}
            >
              <option value="">选择内容</option>
              {contentOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>

            <button type="submit" disabled={loading}>发起发布</button>
          </form>

          <ul className="list">
            {publishTasks.map((task) => (
              <li key={task.id}>
                <div>
                  <strong>Task #{task.id}</strong>
                  <p>
                    content_id={task.content_id}, publisher_id={task.publisher_id}, status={task.status ?? 'unknown'}
                  </p>
                </div>
              </li>
            ))}
          </ul>
        </article>
      </section>
    </main>
  )
}

export default App
