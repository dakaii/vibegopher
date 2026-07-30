import { clearToken, getToken } from './auth'

const API_BASE = (import.meta.env.VITE_API_BASE_URL || '').replace(/\/$/, '')

async function request(path, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  }
  const token = getToken()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  })

  if (res.status === 401) {
    clearToken()
  }

  const text = await res.text()
  let data = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = { error: text }
    }
  }

  if (!res.ok) {
    throw new Error(data?.error || `Request failed (${res.status})`)
  }
  return data
}

export const api = {
  googleAuth: (idToken) =>
    request('/api/auth/google', {
      method: 'POST',
      body: JSON.stringify({ id_token: idToken }),
    }),
  me: () => request('/api/me'),
  posts: () => request('/api/posts'),
  createPost: (content) =>
    request('/api/posts', {
      method: 'POST',
      body: JSON.stringify({ content }),
    }),
  comments: (postId) => request(`/api/comments/post/${postId}`),
  createComment: (postId, content) =>
    request('/api/comments', {
      method: 'POST',
      body: JSON.stringify({ post_id: postId, content }),
    }),
}
