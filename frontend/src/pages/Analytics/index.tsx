import { useEffect } from 'react'
import { useAnalytics } from '../../hooks/useAnalytics'
import './Analytics.css'

/**
 * Analytics Page
 * Analyze and track content performance
 *
 * Features:
 * - Summary statistics
 * - Content ranking
 * - Platform comparison
 * - Engagement metrics
 * - Growth trends
 */
export function Analytics() {
  const {
    summary,
    ranking,
    platformStats,
    isLoading,
    error,
    fetchSummary,
    fetchRanking,
    fetchPlatformComparison,
    refresh,
  } = useAnalytics({
    autoFetch: true,
    pollInterval: 30000,
  })

  if (isLoading && !summary) {
    return <div className="analytics-loading">加载分析数据中...</div>
  }

  return (
    <div className="analytics-container">
      <div className="analytics-header">
        <h1>📊 内容分析</h1>
        <button onClick={refresh} className="btn btn-primary">
          刷新数据
        </button>
      </div>

      {error && (
        <div className="analytics-error">
          <strong>错误:</strong> {error.message}
        </div>
      )}

      {/* Summary Cards */}
      {summary && (
        <div className="analytics-summary">
          <div className="summary-card">
            <div className="summary-label">发布内容</div>
            <div className="summary-value">{summary.totalPublished}</div>
          </div>

          <div className="summary-card">
            <div className="summary-label">总浏览</div>
            <div className="summary-value">{summary.totalViews.toLocaleString()}</div>
            <div className="summary-sub">平均: {Math.round(summary.avgViewsPerContent)}</div>
          </div>

          <div className="summary-card">
            <div className="summary-label">总点赞</div>
            <div className="summary-value">{summary.totalLikes.toLocaleString()}</div>
          </div>

          <div className="summary-card">
            <div className="summary-label">总评论</div>
            <div className="summary-value">{summary.totalComments.toLocaleString()}</div>
          </div>

          <div className="summary-card">
            <div className="summary-label">总分享</div>
            <div className="summary-value">{summary.totalShares.toLocaleString()}</div>
          </div>
        </div>
      )}

      <div className="analytics-main">
        {/* Platform Comparison */}
        {platformStats.length > 0 && (
          <div className="analytics-section">
            <div className="section-header">
              <h2>平台对比</h2>
              <button onClick={fetchPlatformComparison} className="btn-refresh">
                刷新
              </button>
            </div>

            <div className="platform-stats">
              {platformStats.map((platform) => (
                <div key={platform.platform} className="platform-card">
                  <h3>{platform.platform}</h3>

                  <div className="stat-row">
                    <span className="stat-label">内容数</span>
                    <span className="stat-value">{platform.contentCount}</span>
                  </div>

                  <div className="stat-row">
                    <span className="stat-label">总浏览</span>
                    <span className="stat-value">{platform.totalViews.toLocaleString()}</span>
                  </div>

                  <div className="stat-row">
                    <span className="stat-label">总点赞</span>
                    <span className="stat-value">{platform.totalLikes.toLocaleString()}</span>
                  </div>

                  <div className="stat-row">
                    <span className="stat-label">总评论</span>
                    <span className="stat-value">{platform.totalComments.toLocaleString()}</span>
                  </div>

                  <div className="stat-row highlight">
                    <span className="stat-label">参与度</span>
                    <span className="stat-value">
                      {(platform.avgEngagement * 100).toFixed(2)}%
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Top Content */}
        {summary && summary.topContent.length > 0 && (
          <div className="analytics-section">
            <div className="section-header">
              <h2>热门内容</h2>
              <button onClick={fetchRanking} className="btn-refresh">
                刷新
              </button>
            </div>

            <div className="ranking-list">
              {summary.topContent.map((content, index) => (
                <div key={content.id} className="ranking-item">
                  <div className="ranking-badge">{index + 1}</div>

                  <div className="ranking-content">
                    <h3>{content.title}</h3>
                    <div className="ranking-stats">
                      <span>
                        <strong>浏览:</strong> {content.views.toLocaleString()}
                      </span>
                      <span>
                        <strong>点赞:</strong> {content.likes.toLocaleString()}
                      </span>
                      <span>
                        <strong>评论:</strong> {content.comments.toLocaleString()}
                      </span>
                    </div>
                  </div>

                  <div className="ranking-score">
                    {content.views > 0 && (
                      <div className="score-bar">
                        <div
                          className="score-fill"
                          style={{
                            width: `${Math.min(
                              (content.views / (summary.topContent[0]?.views || 1)) * 100,
                              100
                            )}%`,
                          }}
                        />
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* All Rankings */}
        {ranking.length > 0 && (
          <div className="analytics-section">
            <div className="section-header">
              <h2>内容排名</h2>
              <button onClick={fetchRanking} className="btn-refresh">
                刷新
              </button>
            </div>

            <div className="table-wrapper">
              <table className="ranking-table">
                <thead>
                  <tr>
                    <th>排名</th>
                    <th>标题</th>
                    <th>浏览量</th>
                    <th>参与度</th>
                    <th>趋势</th>
                  </tr>
                </thead>
                <tbody>
                  {ranking.slice(0, 10).map((item) => (
                    <tr key={item.contentId}>
                      <td className="rank-number">#{item.rank}</td>
                      <td>{item.title}</td>
                      <td>{item.views.toLocaleString()}</td>
                      <td>{(item.engagement * 100).toFixed(2)}%</td>
                      <td>
                        <span
                          className={`trend-badge trend-${item.trend}`}
                        >
                          {item.trend === 'up' && '↑ 上升'}
                          {item.trend === 'down' && '↓ 下降'}
                          {item.trend === 'stable' && '→ 稳定'}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default Analytics
