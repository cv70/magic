import { useCallback, useRef, useEffect, useState } from 'react'

interface AutoSaveConfig {
  interval?: number // 毫秒，默认30000（30秒）
  onSave?: (data: any) => Promise<void>
}

interface AutoSaveStatus {
  isSaving: boolean
  lastSavedAt: Date | null
  hasPendingChanges: boolean
  error: string | null
}

export function useAutoSave(config: AutoSaveConfig = {}) {
  const { interval = 30000, onSave } = config
  const [status, setStatus] = useState<AutoSaveStatus>({
    isSaving: false,
    lastSavedAt: null,
    hasPendingChanges: false,
    error: null,
  })

  const dataRef = useRef<any>(null)
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null)

  // 标记有改动
  const markAsDirty = useCallback(() => {
    setStatus((prev) => ({ ...prev, hasPendingChanges: true }))
  }, [])

  // 执行保存
  const performSave = useCallback(async () => {
    if (!dataRef.current || !onSave) return

    try {
      setStatus((prev) => ({ ...prev, isSaving: true, error: null }))
      await onSave(dataRef.current)
      setStatus((prev) => ({
        ...prev,
        isSaving: false,
        lastSavedAt: new Date(),
        hasPendingChanges: false,
        error: null,
      }))
      return true
    } catch (error) {
      const errorMsg = error instanceof Error ? error.message : '保存失败'
      setStatus((prev) => ({
        ...prev,
        isSaving: false,
        error: errorMsg,
        hasPendingChanges: true, // 保留待保存状态
      }))
      console.error('Auto-save error:', error)
      return false
    }
  }, [onSave])

  // 设置自动保存定时器
  useEffect(() => {
    if (!onSave) return

    timerRef.current = setInterval(async () => {
      if (status.hasPendingChanges && !status.isSaving) {
        await performSave()
      }
    }, interval)

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current)
      }
    }
  }, [interval, onSave, status.hasPendingChanges, status.isSaving, performSave])

  // 手动保存
  const manualSave = useCallback(async () => {
    return await performSave()
  }, [performSave])

  // 更新数据
  const updateData = useCallback((newData: any) => {
    dataRef.current = newData
    markAsDirty()
  }, [markAsDirty])

  // 清理
  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current)
      }
    }
  }, [])

  return {
    status,
    updateData,
    manualSave,
    markAsDirty,
  }
}
