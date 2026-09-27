import { useAuthStore } from '../stores/authStore'

export function useAuth() {
  const user = useAuthStore((state) => state.user)
  const isAuthenticated = useAuthStore(
    (state) => state.isAuthenticated,
  )
  const initialized = useAuthStore(
    (state) => state.initialized,
  )

  return {
    user,
    isAuthenticated,
    loading: !initialized,
  }
}