import axios from 'axios'
import api from './api'
import { env } from '../config/env'
import type { LoginRequest, LoginResponse, User } from '../types/auth'
import type { RegisterRequest, RegisterResponse } from '../types/register'
import type { ApiResponse } from '../types/api'

export async function login(
  data: LoginRequest,
): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>(
    '/api/v1/auth/login',
    data,
  )

  return response.data
}

export async function getMe(): Promise<User> {
  const response = await api.get<ApiResponse<User>>(
    '/api/v1/auth/me',
  )

  return response.data.data
}

export async function logout(): Promise<void> {
  await api.post('/api/v1/auth/logout')
}

export async function register(
  data: RegisterRequest,
): Promise<RegisterResponse> {
  const response = await api.post<RegisterResponse>(
    '/api/v1/auth/register',
    data,
  )

  return response.data
}

export async function refreshAccessToken(
  refreshToken: string,
): Promise<{ access_token: string; refresh_token: string }> {
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

  return response.data.data
}