import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import Loading from '../components/Loading'

interface ProtectedRouteProps {
  children: ReactNode
}

function ProtectedRoute({ children }: ProtectedRouteProps) {
  const initialized = useAuthStore((state) => state.initialized)
  const isAuthenticated = useAuthStore(
    (state) => state.isAuthenticated,
  )

  if (!initialized) {
    return <Loading message="Checking authentication..." />
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return children
}

export default ProtectedRoute