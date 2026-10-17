import api from './api'
import type { ApiResponse } from '../types/api'
import type { Transaction } from '../types/transaction'

export interface Wallet {
  id: number
  user_id: number
  balance: number
  account_number: number
}

export interface TopUpRequest {
  amount: number
  description?: string
}

export type TopUpResponse = ApiResponse<Transaction>

export async function getWallet(): Promise<Wallet> {
  const response = await api.get<ApiResponse<Wallet>>(
    '/api/v1/wallet',
  )

  return response.data.data
}

export async function topUp(
  data: TopUpRequest,
): Promise<TopUpResponse> {
  const response = await api.post<TopUpResponse>(
    '/api/v1/wallet/topup',
    data,
  )

  return response.data
}