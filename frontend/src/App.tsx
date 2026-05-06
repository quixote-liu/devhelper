import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthGuard } from './components/auth/AuthGuard'
import { AppLayout } from './components/layout/AppLayout'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'
import { JsonPage } from './pages/JsonPage'
import { SchemasPage } from './pages/SchemasPage'
import { ProfilePage } from './pages/ProfilePage'
import { AdminPage } from './pages/AdminPage'

const qc = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={qc}>
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route element={<AuthGuard />}>
            <Route element={<AppLayout />}>
              <Route index element={<Navigate to="/json" replace />} />
              <Route path="/json" element={<JsonPage />} />
              <Route path="/schemas" element={<SchemasPage />} />
              <Route path="/profile" element={<ProfilePage />} />
            </Route>
          </Route>
          <Route element={<AuthGuard adminOnly />}>
            <Route element={<AppLayout />}>
              <Route path="/admin" element={<AdminPage />} />
            </Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
