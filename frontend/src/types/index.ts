// Draft相关类型
export interface Draft {
  id: number
  user_id: number
  title: string
  content: string
  content_type: 'text' | 'code' | 'image' | 'audio' | 'video'
  status: 'draft' | 'archived'
  tags: string[]
  metadata?: Record<string, any>
  last_edited_at?: string
  saved_version_id?: number
  created_at?: string
  updated_at?: string
}

export interface ContentVersion {
  id: number
  content_id: number
  draft_id?: number
  version_num: number
  title: string
  content: string
  change_summary: string
  created_by: number
  created_at: string
}

export interface Content {
  id: number
  user_id: number
  title: string
  body: string
  content_type: string
  status: string
  tags: string[]
  created_at?: string
  updated_at?: string
  published_at?: string
}

export interface Tag {
  id: number
  user_id: number
  name: string
  color: string
  count: number
  created_at?: string
}

export interface Category {
  id: number
  user_id: number
  name: string
  icon: string
  color: string
  order: number
  created_at?: string
}

// PublishTask相关类型
export interface Metrics {
  views: number
  likes: number
  comments: number
  shares: number
  follows: number
  updated_at?: string
}

export interface PlatformResult {
  platform: string
  status: 'pending' | 'publishing' | 'published' | 'failed'
  post_id?: string
  post_url?: string
  error?: string
  published_at?: string
  metrics?: Metrics
}

export interface PublishTask {
  id: number
  user_id: number
  content_id: number
  draft_id?: number
  target_platforms: string[]
  status: 'pending' | 'publishing' | 'published' | 'failed'
  scheduled_at?: string
  published_at?: string
  retry_count: number
  max_retries: number
  error_message?: string
  platform_results: Record<string, PlatformResult>
  created_at?: string
  updated_at?: string
}

export interface PublishAnalytic {
  id: number
  publish_task_id: number
  content_id: number
  user_id: number
  platform: string
  post_id: string
  title: string
  post_url: string
  metrics: Record<string, any>
  last_synced_at?: string
  created_at?: string
  updated_at?: string
}

// AI Template相关类型
export interface PromptTemplate {
  id: number
  user_id?: number
  category: string
  sub_category: string
  name: string
  description: string
  prompt_text: string
  variables: string[]
  example_input?: Record<string, any>
  example_output?: string
  usage_count: number
  rating: number
  tags: string[]
  status: 'active' | 'archived'
  version: number
  created_by: number
  created_at?: string
  updated_at?: string
}

// API响应通用格式
export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
  error?: {
    code: string
    message: string
    details?: Record<string, any>
  }
}

// 编辑器相关
export interface EditorState {
  currentDraft: Draft | null
  isAutoSaving: boolean
  lastSavedAt?: string
  versionNum: number
  isDirty: boolean
}

// 发布相关
export interface PublishState {
  tasks: PublishTask[]
  selectedTask?: PublishTask
  filters: {
    status?: string
    dateRange?: [string, string]
  }
}

// 分析相关
export interface AnalyticsSummary {
  period: string
  total_published: number
  total_views: number
  total_likes: number
  total_comments: number
  avg_likes_per_post: number
  total_new_followers: number
}

export interface ContentRanking {
  rank: number
  content_id: number
  title: string
  views: number
  likes: number
  comments: number
  shares: number
}

export interface PlatformComparison {
  platform: string
  posts: number
  avg_views: number
  avg_likes: number
  success_rate: string
}
