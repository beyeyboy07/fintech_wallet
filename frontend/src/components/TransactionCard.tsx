
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
        <TransactionTypeBadge type={transaction.type} />

        <TransactionStatusBadge
          status={transaction.status}
        />
      </div>

      <p className="transaction-amount">
        {formatCurrency(transaction.amount)}
      </p>

      <div className="transaction-details">
        <div className="transaction-detail">
          <span className="transaction-detail-label">
            Date
          </span>

          <span className="transaction-detail-value">
            {formatDate(transaction.created_at)}
          </span>
        </div>

        <div className="transaction-detail">
          <span className="transaction-detail-label">
            Reference
          </span>

          <span className="transaction-detail-value transaction-reference">
            {transaction.reference}
          </span>
        </div>

        <div className="transaction-detail">
          <span className="transaction-detail-label">
            Description
          </span>

          <span className="transaction-detail-value">
            {transaction.description || '-'}
          </span>
        </div>
      </div>
    </Card>
  )
}

export default TransactionCard