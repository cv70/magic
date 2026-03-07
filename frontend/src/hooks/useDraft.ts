import { useState, useCallback, useRef, useEffect } from 'react'
import { draftApi } from '../utils/api'
import type { Draft } from '../types'

interface UseDraftOptions {
  autoSaveInterval?: number // 自动保存间隔（毫秒），默认30秒
  onSaving?: () => void
  onSaved?: (version: number) => void
  onError?: (error: Error) => void
}

export function useDraft(draftId?: number, options: UseDraftOptions = {}) {
  const {
    autoSaveInterval = 30000,
    onSaving,
    onSaved,
    onError,
  } = options

  const [draft, setDraft] = useState<Draft | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const [isSaving, setIsSaving] = useState(false)
  const [lastSavedAt, setLastSavedAt] = useState<string | null>(null)
  const [isDirty, setIsDirty] = useState(false)

  const saveTimeoutRef = useRef<ReturnType<typeof setTimeout>>()
  const draftRef = useRef<Draft | null>(null)

  // 加载草稿
  const fetchDraft = useCallback(async () => {
    if (!draftId) return

    setLoading(true)
    setError(null)

    try {
      const data = await draftApi.get(draftId)
      setDraft(data)
      draftRef.current = data
      setIsDirty(false)
    } catch (err) {
      const error = err as Error
      setError(error)
      onError?.(error)
    } finally {
      setLoading(false)
    }
  }, [draftId, onError])

  // 保存草稿
  const save = useCallback(
    async (data?: Partial<Draft>) => {
      if (!draftRef.current) return

      const toSave = data || draftRef.current
      setIsSaving(true)
      onSaving?.()

      try {
        const result = await draftApi.update(draftRef.current.id, {
          title: toSave.title,
          content: toSave.content,
          content_type: toSave.content_type,
          tags: toSave.tags,
        })

        const now = new Date().toISOString()
        setLastSavedAt(now)
        setIsDirty(false)
        onSaved?.(result.version_num)

        // 更新本地草稿对象
        if (draftRef.current) {
          draftRef.current.updated_at = result.updated_at
        }

        return result
      } catch (err) {
        const error = err as Error
        setError(error)
        onError?.(error)
        throw error
      } finally {
        setIsSaving(false)
      }
    },
    [onSaving, onSaved, onError]
  )

  // 自动保存逻辑
  const triggerAutoSave = useCallback((newDraft: Draft) => {
    draftRef.current = newDraft
    setIsDirty(true)

    // 清除之前的定时器
    if (saveTimeoutRef.current) {
      clearTimeout(saveTimeoutRef.current)
    }

    // 设置新的定时器
    saveTimeoutRef.current = setTimeout(() => {
      save(newDraft).catch((err) => {
        console.error('Auto save failed:', err)
      })
    }, autoSaveInterval)
  }, [autoSaveInterval, save])

  // 更新草稿内容
  const updateDraft = useCallback((updates: Partial<Draft>) => {
    setDraft((current) => {
      const updated = { ...current, ...updates } as Draft
      triggerAutoSave(updated)
      return updated
    })
  }, [triggerAutoSave])

  // 创建新草稿
  const createNew = useCallback(
    async (title: string, content = '', contentType = 'text') => {
      setLoading(true)
      setError(null)

      try {
        const result = await draftApi.create({
          title,
          content,
          content_type: contentType,
        })

        // 需要重新获取完整的草稿对象
        await fetchDraft()
        return result
      } catch (err) {
        const error = err as Error
        setError(error)
        onError?.(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [fetchDraft, onError]
  )

  // 删除草稿
  const deleteDraft = useCallback(async () => {
    if (!draft?.id) return

    setLoading(true)
    setError(null)

    try {
      await draftApi.delete(draft.id)
      setDraft(null)
      draftRef.current = null
      return true
    } catch (err) {
      const error = err as Error
      setError(error)
      onError?.(error)
      throw error
    } finally {
      setLoading(false)
    }
  }, [draft?.id, onError])

  // 获取版本历史
  const getVersions = useCallback(async () => {
    if (!draft?.id) return []

    try {
      const result = await draftApi.getVersions(draft.id)
      return result.data
    } catch (err) {
      const error = err as Error
      onError?.(error)
      throw error
    }
  }, [draft?.id, onError])

  // 恢复版本
  const revertVersion = useCallback(
    async (versionNum: number) => {
      if (!draft?.id) return

      setLoading(true)
      setError(null)

      try {
        const result = await draftApi.revertVersion(draft.id, versionNum)
        await fetchDraft() // 重新加载更新后的草稿
        onSaved?.(result.version_num)
        return result
      } catch (err) {
        const error = err as Error
        setError(error)
        onError?.(error)
        throw error
      } finally {
        setLoading(false)
      }
    },
    [draft?.id, fetchDraft, onSaved, onError]
  )

  // 获取版本对比
  const getDiff = useCallback(
    async (fromVersion: number, toVersion: number) => {
      if (!draft?.id) return null

      try {
        const result = await draftApi.getDiff(draft.id, fromVersion, toVersion)
        return result
      } catch (err) {
        const error = err as Error
        onError?.(error)
        throw error
      }
    },
    [draft?.id, onError]
  )

  // 初始化时加载草稿
  useEffect(() => {
    if (draftId) {
      fetchDraft()
    }
  }, [draftId, fetchDraft])

  // 清理定时器
  useEffect(() => {
    return () => {
      if (saveTimeoutRef.current) {
        clearTimeout(saveTimeoutRef.current)
      }
    }
  }, [])

  return {
    draft,
    loading,
    error,
    isSaving,
    isDirty,
    lastSavedAt,
    createNew,
    updateDraft,
    save,
    deleteDraft,
    fetchDraft,
    getVersions,
    revertVersion,
    getDiff,
  }
}
