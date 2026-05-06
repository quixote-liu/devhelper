import { Navigate, Outlet } from 'react-router-dom'
import { useAuthStore } from '../../store/auth'

export function AuthGuard({ adminOnly = false }: { adminOnly?: boolean }) {
  const { accessToken, user } = useAuthStore()
  if (!accessToken) return <Navigate to="/login" replace />
  if (adminOnly && user?.role !== 'admin') return <Navigate to="/json" replace />
  return <Outlet />
}
