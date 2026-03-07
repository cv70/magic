import { create } from 'zustand'

interface User {
  id: number
  username: string
  email: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  login: (token: string, user: User) => void
  logout: () => void
  setUser: (user: User) => void
  init: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,
  isAuthenticated: false,

  init: () => {
    // 从 localStorage 恢复状态
    const token = localStorage.getItem('auth_token')
    if (token) {
      set({ token, isAuthenticated: true })
    }
  },

  login: (token: string, user: User) => {
    localStorage.setItem('auth_token', token)
    set({ token, user, isAuthenticated: true })
  },

  logout: () => {
    localStorage.removeItem('auth_token')
    set({ token: null, user: null, isAuthenticated: false })
  },

  setUser: (user: User) => {
    set({ user })
  },
}))

