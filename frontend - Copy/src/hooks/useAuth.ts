import { useEffect, useState } from 'react'
import { getMe, refreshAccessToken } from '../services/authService'
import {
  getAccessToken,
  getRefreshToken,
  setTokens,
  clearTokens,
} from '../utils/tokenStorage'
import { useAuthStore } from '../stores/authStore'

export function useAuth() {
  const user = useAuthStore((state) => state.user)
  const setUser = useAuthStore((state) => state.setUser)
  const clearAuth = useAuthStore((state) => state.clearAuth)

  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function restoreSession() {
      const accessToken = getAccessToken()
      const refreshToken = getRefreshToken()

      if (!accessToken && !refreshToken) {
        setLoading(false)
        return
      }

      try {
        if (!accessToken && refreshToken) {
          const tokenResponse = await refreshAccessToken(refreshToken)

          setTokens(
            tokenResponse.access_token,
            tokenResponse.refresh_token,
          )
        }

        const currentUser = await getMe()

        setUser(currentUser)
      } catch {
        clearTokens()
        clearAuth()
      } finally {
        setLoading(false)
      }
    }

    restoreSession()
  }, [setUser, clearAuth])

  return {
    user,
    isAuthenticated: !!user,
    loading,
  }
}