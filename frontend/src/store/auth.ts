import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface User {
  id: number
  username: string
  email: string
  role: string
}

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  setTokens: (data: { access_token: string; refresh_token: string; user?: User }) => void
  setUser: (user: User) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      setTokens: (data) =>
        set({
          accessToken: data.access_token,
          refreshToken: data.refresh_token,
          ...(data.user ? { user: data.user } : {}),
        }),
      setUser: (user) => set({ user }),
      logout: () => set({ user: null, accessToken: null, refreshToken: null }),
    }),
    { name: 'auth' }
  )
)
