import { create } from 'zustand'
import type { User } from '../types/auth'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  initialized: boolean
  setUser: (user: User) => void
  setInitialized: () => void
  clearAuth: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  initialized: false,

  setUser: (user) =>
    set({
      user,
      isAuthenticated: true,
    }),

  setInitialized: () =>
    set({
      initialized: true,
    }),

  clearAuth: () =>
    set({
      user: null,
      isAuthenticated: false,
    }),
}))