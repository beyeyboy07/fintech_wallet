import type { ApiResponse } from './api'
import type { Transaction } from './transaction'

export interface TransferRequest {
  receiver_id: number
  amount: number
  description?: string
}

export type TransferResponse = ApiResponse<Transaction>