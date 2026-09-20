import axios from 'axios'
import { env } from '../config/env'
import {
  getAccessToken,
  getRefreshToken,
  setTokens,
  clearTokens,
} from '../utils/tokenStorage'
import { useAuthStore } from '../stores/authStore'
import type { ApiResponse } from '../types/api'

const api = axios.create({
  baseURL: env.apiBaseUrl,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const accessToken = getAccessToken()

  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`
  }

  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    if (
      error.response?.status !== 401 ||
      originalRequest._retry ||
      originalRequest.url?.includes('/api/v1/auth/refresh')
    ) {
      return Promise.reject(error)
    }

    const refreshToken = getRefreshToken()

    if (!refreshToken) {
      clearTokens()
      return Promise.reject(error)
    }

    originalRequest._retry = true

    try {
      const response = await axios.post<
        ApiResponse<{
          access_token: string
          refresh_token: string
        }>
      >(
        `${env.apiBaseUrl}/api/v1/auth/refresh`,
        {
          refresh_token: refreshToken,
        },
      )

      const newAccessToken = response.data.data.access_token
      const newRefreshToken = response.data.data.refresh_token

      setTokens(newAccessToken, newRefreshToken)

      originalRequest.headers.Authorization = `Bearer ${newAccessToken}`

      return api(originalRequest)
    } catch {
        clearTokens()
        useAuthStore.getState().clearAuth()

        return Promise.reject(error)
    }
  },
)

export default api