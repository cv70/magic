import { useCallback, useEffect, useState } from 'react'
import { publishApi } from '../utils/api'
import type { PublishTask, PlatformResult } from '../types'

interface UsePublishOptions {
  onPublishing?: () => void
  onPublished?: (task: PublishTask) => void
  onError?: (error: Error) => void
}

interface UsePublishState {
  publishTask: PublishTask | null
  isPublishing: boolean
  isPending: boolean
  error: Error | null
  progress: number // 0-100
}

/**
 * usePublish Hook
 * Manages content publishing to multiple platforms
 *
 * Features:
 * - Publish content to configured platforms
 * - Track publishing progress
 * - Monitor platform-specific results
 * - Handle scheduled publishing
 */
export function usePublish(options?: UsePublishOptions) {
  const [state, setState] = useState<UsePublishState>({
    publishTask: null,
    isPublishing: false,
    isPending: false,
    error: null,
    progress: 0,
  })

  // Publish content to specified platforms
  const publish = useCallback(
    async (
      contentId: number,
      platforms: string[],
      scheduledAt?: Date
    ) => {
      setState((prev) => ({
        ...prev,
        isPublishing: true,
        error: null,
      }))

      options?.onPublishing?.()

      try {
        const result = await publishApi.create({
          content_id: contentId,
          target_platforms: platforms,
          scheduled_at: scheduledAt,
          status: scheduledAt ? 'pending' : 'publishing',
        })

        setState((prev) => ({
          ...prev,
          publishTask: result,
          isPublishing: false,
          isPending: scheduledAt ? true : false,
          progress: scheduledAt ? 0 : 100,
        }))

        options?.onPublished?.(result)

        return result
      } catch (err) {
        const error = err instanceof Error ? err : new Error('Failed to publish content')
        setState((prev) => ({
          ...prev,
          isPublishing: false,
          error,
        }))
        options?.onError?.(error)
        throw error
      }
    },
    [options]
  )

  // Get publish task details
  const getTask = useCallback(async (taskId: number) => {
    try {
      const task = await publishApi.get(taskId)
      setState((prev) => ({
        ...prev,
        publishTask: task,
      }))
      return task
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to fetch task')
      setState((prev) => ({
        ...prev,
        error,
      }))
      throw error
    }
  }, [])

  // Update publish task (e.g., reschedule)
  const updateTask = useCallback(
    async (taskId: number, updates: Partial<PublishTask>) => {
      try {
        const updated = await publishApi.update(taskId, updates)
        setState((prev) => ({
          ...prev,
          publishTask: updated,
        }))
        return updated
      } catch (err) {
        const error = err instanceof Error ? err : new Error('Failed to update task')
        setState((prev) => ({
          ...prev,
          error,
        }))
        throw error
      }
    },
    []
  )

  // Execute pending task (start publishing)
  const executeTask = useCallback(async (taskId: number) => {
    setState((prev) => ({
      ...prev,
      isPublishing: true,
      error: null,
    }))

    try {
      const result = await publishApi.execute(taskId)
      setState((prev) => ({
        ...prev,
        publishTask: result,
        isPublishing: false,
        isPending: false,
        progress: 100,
      }))
      return result
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to execute publish task')
      setState((prev) => ({
        ...prev,
        isPublishing: false,
        error,
      }))
      throw error
    }
  }, [])

  // Cancel pending task
  const cancelTask = useCallback(async (taskId: number) => {
    try {
      await publishApi.cancel(taskId)
      setState((prev) => ({
        ...prev,
        publishTask: null,
        isPending: false,
      }))
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to cancel task')
      setState((prev) => ({
        ...prev,
        error,
      }))
      throw error
    }
  }, [])

  // Delete publish task
  const deleteTask = useCallback(async (taskId: number) => {
    try {
      await publishApi.delete(taskId)
      setState((prev) => ({
        ...prev,
        publishTask: null,
      }))
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to delete task')
      setState((prev) => ({
        ...prev,
        error,
      }))
      throw error
    }
  }, [])

  // List all publish tasks
  const listTasks = useCallback(
    async (page: number = 1, limit: number = 20) => {
      try {
        const tasks = await publishApi.list(page, limit)
        return tasks
      } catch (err) {
        const error = err instanceof Error ? err : new Error('Failed to list tasks')
        setState((prev) => ({
          ...prev,
          error,
        }))
        throw error
      }
    },
    []
  )

  // Get platform results
  const getPlatformResults = useCallback((taskId: number) => {
    if (!state.publishTask || state.publishTask.id !== taskId) {
      return []
    }

    try {
      if (typeof state.publishTask.platform_results === 'string') {
        return JSON.parse(state.publishTask.platform_results) as PlatformResult[]
      }
      return (state.publishTask.platform_results as PlatformResult[]) || []
    } catch {
      return []
    }
  }, [state.publishTask])

  // Poll task status (for real-time updates)
  useEffect(() => {
    if (!state.publishTask || state.publishTask.status === 'published' || state.publishTask.status === 'failed') {
      return
    }

    const pollInterval = setInterval(() => {
      getTask(state.publishTask!.id).catch(() => {
        // Silently handle errors in polling
      })
    }, 5000) // Poll every 5 seconds

    return () => clearInterval(pollInterval)
  }, [state.publishTask, getTask])

  return {
    // State
    publishTask: state.publishTask,
    isPublishing: state.isPublishing,
    isPending: state.isPending,
    error: state.error,
    progress: state.progress,

    // Methods
    publish,
    getTask,
    updateTask,
    executeTask,
    cancelTask,
    deleteTask,
    listTasks,
    getPlatformResults,
  }
}

export default usePublish
