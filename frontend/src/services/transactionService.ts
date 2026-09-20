import api from './api'
import type { ApiResponse } from '../types/api'
import type { Transaction } from '../types/transaction'

export interface TransactionPagination {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface TransactionListData {
  items: Transaction[]
  pagination: TransactionPagination
}

export async function getTransactions(
  page = 1,
  limit = 10,
): Promise<TransactionListData> {
  const response = await api.get<ApiResponse<TransactionListData>>(
    '/api/v1/transactions',
    {
      params: {
        page,
        limit,
      },
    },
  )

  return response.data.data
}