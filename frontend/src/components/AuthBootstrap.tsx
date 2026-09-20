import { useEffect, type ReactNode } from 'react'
import { getMe, refreshAccessToken } from '../services/authService'
import {
  getAccessToken,
  getRefreshToken,
  setTokens,
  clearTokens,
} from '../utils/tokenStorage'
import { useAuthStore } from '../stores/authStore'
import Loading from './Loading'



interface AuthBootstrapProps {
  children: ReactNode
}

function AuthBootstrap({ children }: AuthBootstrapProps) {
  const initialized = useAuthStore((state) => state.initialized)
  const setUser = useAuthStore((state) => state.setUser)
  const setInitialized = useAuthStore(
    (state) => state.setInitialized,
  )
  const clearAuth = useAuthStore((state) => state.clearAuth)

  useEffect(() => {
    async function restoreSession() {
      const accessToken = getAccessToken()
      const refreshToken = getRefreshToken()

      if (!accessToken && !refreshToken) {
        setInitialized()
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
        setInitialized()
      }
    }

    if (!initialized) {
      restoreSession()
    }
  }, [
    initialized,
    setUser,
    setInitialized,
    clearAuth,
  ])

  if (!initialized) {
    return <Loading message="Loading application..." />
  }

  return children
}

export default AuthBootstrap