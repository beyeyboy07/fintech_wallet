import type { TransactionStatus } from '../types/transaction'

interface TransactionStatusBadgeProps {
  status: TransactionStatus
}

const statusClassMap: Record<TransactionStatus, string> = {
  PENDING: 'status-pending',
  PROCESSING: 'status-processing',
  SUCCESS: 'status-success',
  FAILED: 'status-failed',
  CANCELLED: 'status-cancelled',
  EXPIRED: 'status-expired',
}

function TransactionStatusBadge({
  status,
}: TransactionStatusBadgeProps) {
  return (
    <span className={statusClassMap[status]}>
      {status}
    </span>
  )
}

export default TransactionStatusBadge