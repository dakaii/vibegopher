import { clearLegacyToken, getAccessToken } from './auth'
import type { ApiErrorBody, Comment, Post, User } from './types'

const API_BASE = (import.meta.env.VITE_API_BASE_URL || '').replace(/\/$/, '')

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers)
  if (!headers.has('Content-Type') && options.body) {
    headers.set('Content-Type', 'application/json')
  }

  const token = await getAccessToken()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  })

  if (res.status === 401) {
    clearLegacyToken()
  }

  const text = await res.text()
  let data: unknown = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = { error: text } satisfies ApiErrorBody
    }
  }

  if (!res.ok) {
    const body = data as ApiErrorBody | null
    throw new Error(body?.error || `Request failed (${res.status})`)
  }

  return data as T
}

export function postCursor(post: Post): string {
  return `${post.created_at}|${post.id}`
}

export const api = {
  me: () => request<User>('/api/me'),
  posts: (opts?: { limit?: number; before?: string }) => {
    const params = new URLSearchParams()
    if (opts?.limit) params.set('limit', String(opts.limit))
    if (opts?.before) params.set('before', opts.before)
    const q = params.toString()
    return request<Post[]>(`/api/posts${q ? `?${q}` : ''}`)
  },
  createPost: (content: string) =>
    request<Post>('/api/posts', {
      method: 'POST',
      body: JSON.stringify({ content }),
    }),
  comments: (postId: string) => request<Comment[]>(`/api/comments/post/${postId}`),
  createComment: (postId: string, content: string) =>
    request<Comment>('/api/comments', {
      method: 'POST',
      body: JSON.stringify({ post_id: postId, content }),
    }),
}
