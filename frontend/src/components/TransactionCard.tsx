import type { Transaction } from '../types/transaction'
import { formatCurrency } from '../utils/currency'
import { formatDate } from '../utils/date'
import TransactionStatusBadge from './TransactionStatusBadge'
import TransactionTypeBadge from './TransactionTypeBadge'
import Card from './Card'
import '../assets/css/TransactionCard.css'

interface TransactionCardProps {
  transaction: Transaction
}

function TransactionCard({
  transaction,
}: TransactionCardProps) {
  return (
    <Card className="transaction-card">
      <div className="transaction-card-header">
        <strong>
          <TransactionTypeBadge type={transaction.type} />
        </strong>

        <TransactionStatusBadge
          status={transaction.status}
        />
      </div>

      <p className="transaction-amount">
        {formatCurrency(transaction.amount)}
      </p>

      <div className="transaction-details">
        <p>
          <strong>Date:</strong>{' '}
          {formatDate(transaction.created_at)}
        </p>

        <p>
          <strong>Reference:</strong>{' '}
          {transaction.reference}
        </p>

        <p>
          <strong>Description:</strong>{' '}
          {transaction.description || '-'}
        </p>
      </div>

      <hr className="transaction-card-divider" />
    </Card>
  )
}

export default TransactionCard