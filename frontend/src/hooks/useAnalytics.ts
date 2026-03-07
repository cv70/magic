import { useCallback, useEffect, useState } from 'react'
import { analyticsApi } from '../utils/api'

interface AnalyticsSummary {
  totalPublished: number
  totalViews: number
  totalLikes: number
  totalComments: number
  totalShares: number
  avgViewsPerContent: number
  topContent: Array<{
    id: number
    title: string
    views: number
    likes: number
    comments: number
  }>
}

interface PlatformStats {
  platform: string
  contentCount: number
  totalViews: number
  totalLikes: number
  totalComments: number
  avgEngagement: number // (likes + comments + shares) / views
}

interface RankingData {
  rank: number
  contentId: number
  title: string
  views: number
  engagement: number
  trend: 'up' | 'down' | 'stable'
}

interface UseAnalyticsOptions {
  autoFetch?: boolean
  pollInterval?: number
  onError?: (error: Error) => void
}

/**
 * useAnalytics Hook
 * Manages analytics data for published content
 *
 * Features:
 * - Fetch analytics summary
 * - Track platform-specific stats
 * - Monitor content ranking
 * - Compare platform performance
 * - Find best publishing times
 */
export function useAnalytics(options?: UseAnalyticsOptions) {
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null)
  const [ranking, setRanking] = useState<RankingData[]>([])
  const [platformStats, setPlatformStats] = useState<PlatformStats[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  const {
    autoFetch = true,
    pollInterval = 30000, // 30 seconds
  } = options || {}

  // Fetch analytics summary
  const fetchSummary = useCallback(async (days: number = 30) => {
    setIsLoading(true)
    setError(null)

    try {
      const data = await analyticsApi.summary({ days })
      setSummary(data)
      return data
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to fetch analytics summary')
      setError(error)
      options?.onError?.(error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }, [options])

  // Fetch ranking data
  const fetchRanking = useCallback(async (limit: number = 10) => {
    setIsLoading(true)
    setError(null)

    try {
      const data = await analyticsApi.ranking({ limit })
      setRanking(data)
      return data
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to fetch ranking')
      setError(error)
      options?.onError?.(error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }, [options])

  // Fetch platform comparison
  const fetchPlatformComparison = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const data = await analyticsApi.platformComparison()
      setPlatformStats(data)
      return data
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to fetch platform stats')
      setError(error)
      options?.onError?.(error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }, [options])

  // Find best publishing time
  const findBestPublishTime = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const data = await analyticsApi.bestPublishTime()
      return data
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to fetch best publish time')
      setError(error)
      options?.onError?.(error)
      throw error
    } finally {
      setIsLoading(false)
    }
  }, [options])

  // Calculate engagement rate
  const getEngagementRate = useCallback((contentId: number): number => {
    const data = ranking.find((r) => r.contentId === contentId)
    return data?.engagement ?? 0
  }, [ranking])

  // Get top performing platform
  const getTopPlatform = useCallback((): PlatformStats | null => {
    if (platformStats.length === 0) return null
    return platformStats.reduce((prev, current) =>
      current.totalViews > prev.totalViews ? current : prev
    )
  }, [platformStats])

  // Get growth trends
  const getGrowthTrends = useCallback(async (days: number = 30) => {
    try {
      const current = await analyticsApi.summary({ days })
      const previous = await analyticsApi.summary({ days, offset: days })

      if (!summary || !current) return null

      return {
        viewGrowth: ((current.totalViews - (previous?.totalViews || 0)) / (previous?.totalViews || 1)) * 100,
        likeGrowth: ((current.totalLikes - (previous?.totalLikes || 0)) / (previous?.totalLikes || 1)) * 100,
        commentGrowth: ((current.totalComments - (previous?.totalComments || 0)) / (previous?.totalComments || 1)) * 100,
      }
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to calculate growth trends')
      setError(error)
      throw error
    }
  }, [summary])

  // Auto-fetch on mount
  useEffect(() => {
    if (!autoFetch) return

    const fetchAll = async () => {
      try {
        await Promise.all([
          fetchSummary(),
          fetchRanking(),
          fetchPlatformComparison(),
        ])
      } catch {
        // Errors handled by individual fetch functions
      }
    }

    fetchAll()

    // Set up polling if pollInterval is specified
    if (pollInterval > 0) {
      const interval = setInterval(fetchAll, pollInterval)
      return () => clearInterval(interval)
    }
  }, [autoFetch, pollInterval, fetchSummary, fetchRanking, fetchPlatformComparison])

  return {
    // State
    summary,
    ranking,
    platformStats,
    isLoading,
    error,

    // Methods
    fetchSummary,
    fetchRanking,
    fetchPlatformComparison,
    findBestPublishTime,
    getEngagementRate,
    getTopPlatform,
    getGrowthTrends,
    refresh: () => Promise.all([
      fetchSummary(),
      fetchRanking(),
      fetchPlatformComparison(),
    ]),
  }
}

export default useAnalytics
