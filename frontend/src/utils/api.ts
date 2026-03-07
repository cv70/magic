import type { ApiResponse, Draft, Content, PublishTask, PromptTemplate, Tag, Category } from '../types'

const API_BASE = '/api/v1'

// ===================== 通用请求函数 =====================

export async function request<T>(
  method: 'GET' | 'POST' | 'PUT' | 'DELETE',
  endpoint: string,
  data?: any
): Promise<T> {
  const url = `${API_BASE}${endpoint}`

  const options: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json',
    },
  }

  // 添加认证token（如果有）
  const token = localStorage.getItem('auth_token')
  if (token) {
    options.headers = {
      ...options.headers,
      'Authorization': `Bearer ${token}`,
    }
  }

  if (data && (method === 'POST' || method === 'PUT')) {
    options.body = JSON.stringify(data)
  }

  const response = await fetch(url, options)

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.message || `HTTP ${response.status}`)
  }

  const result: ApiResponse<T> = await response.json()

  if (result.code !== 200) {
    throw new Error(result.message || '请求失败')
  }

  return result.data as T
}

// ===================== Auth API =====================

export const authApi = {
  // 用户登录
  login: (username: string, password: string) =>
    request<{ token: string; user: { id: number; username: string; email: string } }>(
      'POST',
      '/auth/login',
      { username, password }
    ),

  // 用户注册
  register: (username: string, email: string, password: string) =>
    request<{ token: string; user: { id: number; username: string; email: string } }>(
      'POST',
      '/auth/register',
      { username, email, password }
    ),

  // 获取当前用户信息
  me: () =>
    request<{ id: number; username: string; email: string }>('POST', '/auth/me', {}),
}

// ===================== Draft API =====================

export const draftApi = {
  // 创建草稿
  create: (data: { title: string; content?: string; content_type?: string; tags?: string[]; metadata?: Record<string, any> }) =>
    request<{ id: number; title: string; created_at: string }>('POST', '/drafts/create', data),

  // 获取草稿列表
  list: (page = 1, limit = 20, keyword?: string, status?: string) =>
    request<{ data: Draft[]; total: number }>('POST', '/drafts/search', {
      keyword,
      status,
      page,
      limit,
    }),

  // 获取单个草稿
  get: (id: number) => request<Draft>('POST', '/drafts/get', { id }),

  // 更新草稿
  update: (id: number, data: Partial<Draft> & { change_summary?: string }) =>
    request<{ id: number; updated_at: string; version_num: number }>('POST', '/drafts/update', {
      id,
      ...data,
    }),

  // 删除草稿
  delete: (id: number) => request<{ success: boolean }>('POST', '/drafts/delete', { id }),

  // 获取版本历史
  getVersions: (id: number) =>
    request<{ data: Array<{ version_num: number; change_summary: string; created_at: string }> }>(
      'POST',
      '/drafts/versions',
      { id }
    ),

  // 恢复版本
  revertVersion: (id: number, versionNum: number) =>
    request<{ id: number; version_num: number; reverted_from: number; updated_at: string }>(
      'POST',
      `/drafts/${id}/revert/${versionNum}`,
      {}
    ),

  // 版本对比
  getDiff: (id: number, fromVersion: number, toVersion: number) =>
    request<{ from_version: number; to_version: number; diff: any }>('POST', '/drafts/diff', {
      id,
      from_version: fromVersion,
      to_version: toVersion,
    }),

  // 草稿发布为正式内容
  publish: (id: number, status: 'published' | 'archived') =>
    request<{ content_id: number; status: string; published_at: string }>('POST', `/drafts/${id}/publish`, {
      status,
    }),
}

// ===================== Content API =====================

export const contentApi = {
  // 获取内容列表
  list: (page = 1, limit = 20, filters?: any) =>
    request<{ data: Content[]; total: number }>('POST', '/content/search', {
      ...filters,
      page,
      limit,
    }),

  // 获取单个内容
  get: (id: number) => request<Content>('POST', '/content/get', { id }),

  // 删除内容
  delete: (id: number) => request<{ success: boolean }>('POST', '/content/delete', { id }),
}

// ===================== Publisher API =====================

export const publisherApi = {
  // 获取发布者列表
  search: (filters?: any) =>
    request<{ data: any[]; total: number }>('POST', '/publish/publisher/search', filters || {}),

  // 获取单个发布者
  get: (id: number) => request<any>('POST', '/publish/publisher/get', { id }),

  // 创建发布者
  add: (data: any) => request<any>('POST', '/publish/publisher/add', data),

  // 更新发布者
  update: (id: number, data: any) => request<any>('POST', '/publish/publisher/update', { id, ...data }),
}

// ===================== PublishTask API =====================

export const publishApi = {
  // 创建发布任务
  create: (data: {
    content_id?: number
    draft_id?: number
    target_platforms: string[]
    scheduled_at?: string
    max_retries?: number
  }) => request<PublishTask>('POST', '/publish-tasks', data),

  // 获取发布任务列表
  list: (page = 1, limit = 20, filters?: any) =>
    request<{ data: PublishTask[]; total: number }>('POST', '/publish-tasks/search', {
      ...filters,
      page,
      limit,
    }),

  // 获取单个发布任务
  get: (id: number) => request<PublishTask>('POST', '/publish/task/get', { id }),

  // 修改发布任务
  update: (id: number, data: Partial<PublishTask>) =>
    request<PublishTask>('POST', `/publish-tasks/${id}`, data),

  // 立即执行发布
  execute: (id: number) => request<{ id: number; status: string }>('POST', `/publish-tasks/${id}/execute`, {}),

  // 取消发布任务
  cancel: (id: number) => request<{ success: boolean }>('POST', `/publish-tasks/${id}/cancel`, {}),

  // 删除发布任务
  delete: (id: number) => request<{ success: boolean }>('POST', `/publish-tasks/${id}`, { _method: 'DELETE' }),
}

// ===================== Analytics API =====================

export const analyticsApi = {
  // 获取统计摘要
  summary: (options?: { days?: number }) =>
    request<{
      period: string
      total_published: number
      total_views: number
      total_likes: number
      total_comments: number
      avg_likes_per_post: number
      total_new_followers: number
    }>('POST', '/analytics/summary', { days: options?.days || 30 }),

  // 获取内容排行
  ranking: (options?: { metric?: string; days?: number; limit?: number }) =>
    request<{ data: any[] }>('POST', '/analytics/ranking', {
      metric: options?.metric || 'views',
      days: options?.days || 30,
      limit: options?.limit || 10,
    }),

  // 平台对比
  platformComparison: () => request<{ data: any[] }>('POST', '/analytics/platform-comparison', {}),

  // 最佳发布时间
  bestPublishTime: () =>
    request<{ recommendations: Array<{ time: string; reason: string; performance: any }> }>(
      'POST',
      '/analytics/best-publish-time',
      {}
    ),
}

// ===================== Tag/Category API =====================

export const tagApi = {
  // 获取所有标签
  list: () => request<{ data: Tag[] }>('POST', '/tags', {}),

  // 创建标签
  create: (data: { name: string; color: string }) =>
    request<Tag>('POST', '/tags/create', data),

  // 删除标签
  delete: (id: number) => request<{ success: boolean }>('POST', `/tags/${id}/delete`, {}),
}

export const categoryApi = {
  // 获取所有分类
  list: () => request<{ data: Category[] }>('POST', '/categories', {}),

  // 创建分类
  create: (data: { name: string; icon: string; color: string }) =>
    request<Category>('POST', '/categories/create', data),

  // 修改分类顺序
  updateOrder: (id: number, order: number) =>
    request<Category>('POST', `/categories/${id}/order`, { order }),
}

// ===================== AI Generation API =====================

export const aiApi = {
  // 生成内容
  generate: (data: { generatorID: number; input: string; title?: string }) =>
    request<{ id: number }>('POST', '/ai/generate', data),

  // 获取生成任务状态
  getTask: (id: number) =>
    request<{
      id: number
      status: string
      content?: string
      error?: string
      createdAt: string
      completedAt?: string
    }>('POST', '/ai/task/get', { id }),

  // 列出生成器
  listGenerators: () =>
    request<{ data: any[] }>('POST', '/ai/generator/search', {}),
}

// ===================== PromptTemplate API =====================

export const promptApi = {
  // 获取提示词模板列表
  list: (category?: string, subCategory?: string, userTemplates = false) =>
    request<{ data: PromptTemplate[] }>('POST', '/prompt-templates', {
      category,
      sub_category: subCategory,
      user_templates: userTemplates,
    }),

  // 获取单个模板
  get: (id: number) => request<PromptTemplate>('POST', '/prompt-templates/get', { id }),

  // 创建自定义模板
  create: (data: Partial<PromptTemplate>) =>
    request<PromptTemplate>('POST', '/prompt-templates/create', data),

  // 删除模板
  delete: (id: number) => request<{ success: boolean }>('POST', `/prompt-templates/${id}/delete`, {}),
}

// ===================== WebSocket 连接 =====================

export function connectPublishStatus(taskId: number, onUpdate: (data: any) => void): () => void {
  // 创建WebSocket连接用于实时推送发布状态
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/api/v1/publish-tasks/${taskId}/stream`

  const ws = new WebSocket(wsUrl)

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    onUpdate(data)
  }

  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
  }

  // 返回断开连接的函数
  return () => {
    ws.close()
  }
}

// ===================== 内容广场 API =====================

export const squareApi = {
  // 获取广场内容列表
  list: (params: {
    domain?: string
    keyword?: string
    sort?: 'newest' | 'hottest' | 'trending'
    page: number
    page_size: number
  }) =>
    request<{ data: any[]; total: number }>('POST', '/square/posts', params),

  // 获取内容详情
  get: (id: number) =>
    request<any>('POST', '/square/posts/get', { id }),

  // 将草稿发布到广场
  publish: (draftId: number) =>
    request<{ id: number }>('POST', '/square/publish', { draft_id: draftId }),

  // 点赞内容
  like: (postId: number) =>
    request<{ success: boolean }>('POST', '/square/like', { post_id: postId }),

  // 取消点赞
  unlike: (postId: number) =>
    request<{ success: boolean }>('POST', '/square/unlike', { post_id: postId }),

  // 收藏内容
  collect: (postId: number) =>
    request<{ success: boolean }>('POST', '/square/collect', { post_id: postId }),

  // 取消收藏
  uncollect: (postId: number) =>
    request<{ success: boolean }>('POST', '/square/uncollect', { post_id: postId }),

  // 发表评论
  comment: (postId: number, content: string) =>
    request<{ id: number }>('POST', '/square/comment', { post_id: postId, content }),

  // 获取评论列表
  getComments: (postId: number, page: number = 1, pageSize: number = 20) =>
    request<{ data: any[]; total: number }>('POST', '/square/comments', {
      post_id: postId,
      page,
      page_size: pageSize
    }),
}

