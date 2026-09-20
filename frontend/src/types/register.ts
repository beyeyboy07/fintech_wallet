import type { ApiResponse } from './api'
import type { User } from './auth'

export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export type RegisterResponse = ApiResponse<User>