import api from './api'
import type {
  TransferRequest,
  TransferResponse,
} from '../types/transfer'

export async function transfer(
  data: TransferRequest,
  idempotencyKey: string,
): Promise<TransferResponse> {
  const response = await api.post<TransferResponse>(
    '/api/v1/transfers',
    data,
    {
      headers: {
        'Idempotency-Key': idempotencyKey,
      },
    },
  )

  return response.data
}