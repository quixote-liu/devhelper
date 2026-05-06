import axios from 'axios'
import { useAuthStore } from '../store/auth'

export const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
})

client.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

let refreshing: Promise<void> | null = null

client.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config
    if (error.response?.status !== 401 || original._retry) {
      return Promise.reject(error)
    }
    original._retry = true
    if (!refreshing) {
      refreshing = (async () => {
        const { refreshToken, setTokens, logout } = useAuthStore.getState()
        try {
          const res = await axios.post(
            `${client.defaults.baseURL}/auth/refresh`,
            { refresh_token: refreshToken }
          )
          setTokens(res.data.data)
        } catch {
          logout()
          window.location.href = '/login'
        } finally {
          refreshing = null
        }
      })()
    }
    await refreshing
    return client(original)
  }
)
