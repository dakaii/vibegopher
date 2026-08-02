export interface User {
  id: string
  username: string
  email?: string
  is_bot?: boolean
}

export interface AuthToken {
  tokenType: string
  token: string
  expiresIn: number
}

export interface Post {
  id: string
  content: string
  created_at: string
  updated_at?: string
  user_id: string
  user?: User
}

export interface Comment {
  id: string
  content: string
  created_at: string
  updated_at?: string
  user_id: string
  post_id: string
  user?: User
}

export interface ApiErrorBody {
  error?: string
}
