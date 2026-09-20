import type { TransactionType } from '../types/transaction'

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
    <span>
      {typeLabelMap[type]}
    </span>
  )
}

export default TransactionTypeBadge