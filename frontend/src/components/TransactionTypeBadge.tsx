import type { TransactionType } from '../types/transaction'
import '../assets/css/TransactionTypeBadge.css'

interface TransactionTypeBadgeProps {
  type: TransactionType
}

const typeLabelMap: Record<TransactionType, string> = {
  TOPUP: 'Top Up',
  TRANSFER: 'Transfer',
}

function TransactionTypeBadge({
  type,
}: TransactionTypeBadgeProps) {
  return (
    <span className="transaction-type-badge">
      {typeLabelMap[type]}
    </span>
  )
}

export default TransactionTypeBadge