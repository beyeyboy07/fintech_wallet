export type TransactionType =
  | 'TOPUP'
  | 'TRANSFER'

export type TransactionStatus =
  | 'PENDING'
  | 'PROCESSING'
  | 'SUCCESS'
  | 'FAILED'
  | 'CANCELLED'
  | 'EXPIRED'

export interface Transaction {
  id: number
  user_id: number
  sender_id: number | null
  receiver_id: number | null
  type: TransactionType
  amount: number
  status: TransactionStatus
  reference: string
  description: string
  created_at: string
  updated_at: string
}