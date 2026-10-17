export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  success: boolean
  message: string
  data: {
    access_token: string
    refresh_token: string
    account_number: string
  }
}

export interface User {
  id: number
  name: string
  email: string
  role: string
  account_number: string
}