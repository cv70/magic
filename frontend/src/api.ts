export type ApiResponse<T> = {
  code: number
  message?: string
  data?: T
}

export type Content = {
  id: number
  title: string
  content: string
  content_type?: string
  status?: string
  tags?: string[]
  created_at?: string
  updated_at?: string
}

export type Generator = {
  id: number
  name: string
  provider: string
  model: string
  enabled?: boolean
}

export type Publisher = {
  id: number
  name: string
  platform: string
  enabled?: boolean
}

export type PublishTask = {
  id: number
  publisher_id: number
  content_id: number
  content: string
  status?: string
  created_at?: string
}

const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8888/api/v1'

async function request<T>(path: string, payload: unknown): Promise<T> {
  let response: Response
  try {
    response = await fetch(`${API_BASE}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
  } catch {
    throw new Error('连接后端失败，请检查服务是否启动')
  }

  const json = (await response.json()) as ApiResponse<T>
  if (json.code !== 200) {
    throw new Error(json.message ?? '请求失败')
  }

  if (json.data === undefined) {
    throw new Error('响应缺少 data 字段')
  }

  return json.data
}

export const api = {
  searchContents: (query?: string) =>
    request<{ contents: Content[]; total: number }>('/content/search', {
      query: query || null,
      content_type: null,
      status: null,
      tag: null,
      page: 1,
      limit: 50,
    }),

  addContent: (title: string, content: string, contentType?: string, tags?: string[]) =>
    request<{ id: number }>('/content/add', {
      title,
      content,
      content_type: contentType || null,
      tags: tags && tags.length > 0 ? tags : null,
    }),

  updateContent: (id: number, title: string, content: string, contentType?: string) =>
    request<{ id: number }>('/content/update', {
      id,
      title,
      content,
      content_type: contentType || null,
      tags: null,
      status: null,
    }),

  deleteContent: (id: number) => request<{ id: number }>('/content/delete', { id }),

  searchGenerators: () =>
    request<{ generators: Generator[]; total: number }>('/ai/generator/search', {
      provider: null,
      model: null,
      enabled: true,
      page: 1,
      limit: 50,
    }),

  generate: (generatorId: number, input: string) =>
    request<{ id: number }>('/ai/generate', {
      generator_id: generatorId,
      input,
    }),

  searchPublishers: () =>
    request<{ publishers: Publisher[]; total: number }>('/publishing/publisher/search', {
      platform: null,
      enabled: true,
      page: 1,
      limit: 50,
    }),

  publishContent: (publisherId: number, contentId: number) =>
    request<{ id: number }>('/publishing/content/publish', {
      publisher_id: publisherId,
      content_id: contentId,
    }),

  searchPublishTasks: () =>
    request<{ publish_tasks: PublishTask[]; total: number }>('/publishing/task/search', {
      status: null,
      platform: null,
      created_at: null,
      page: 1,
      limit: 50,
    }),
}
