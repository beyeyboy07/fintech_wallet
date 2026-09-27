import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getWallet } from '../services/walletService'
import type { Wallet } from '../services/walletService'
import { getTransactions } from '../services/transactionService'
import type { Transaction } from '../types/transaction'
import { useAuthStore } from '../stores/authStore'
import { formatCurrency } from '../utils/currency'
import Loading from '../components/Loading'
import ErrorState from '../components/ErrorState'
import PageHeader from '../components/PageHeader'
import Card from '../components/Card'
import TransactionCard from '../components/TransactionCard'
import EmptyState from '../components/EmptyState'
import '../assets/css/DashboardPage.css'

function DashboardPage() {
  const [wallet, setWallet] = useState<Wallet | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const user = useAuthStore((state) => state.user)

  useEffect(() => {
    async function loadDashboard() {
      try {
        const [walletData, transactionData] =
          await Promise.all([
            getWallet(),
            getTransactions(1, 5),
          ])

        setWallet(walletData)
        setTransactions(transactionData.items)
      } catch {
        setError('Gagal mengambil data dashboard.')
      } finally {
        setLoading(false)
      }
    }

    loadDashboard()
  }, [])

  if (loading) {
    return <Loading message="Loading dashboard..." />
  }

  if (error) {
    return <ErrorState message={error} />
  }

  return (
    <section className="dashboard-page">
      <PageHeader
        title={`Welcome, ${user?.name ?? ''}`}
        description="Manage your wallet and transactions."
      />

      <Card className="balance-card">
        <div className="balance-card-header">
          <p className="balance-label">Wallet Balance</p>
          <span className="balance-brand">DANQIVA</span>
        </div>

        <h3 className="balance-amount">
          {wallet
            ? formatCurrency(wallet.balance)
            : formatCurrency(0)}
        </h3>

        <p className="balance-caption">
          Available balance
        </p>
      </Card>



      <div className="quick-actions">
        <Link to="/top-up">
          <span className="quick-action-icon" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none">
              <path
                d="M12 5v14M5 12h14"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </svg>
          </span>

          <span className="quick-action-content">
            <strong>Top Up</strong>
            <small>Add balance</small>
          </span>
        </Link>

        <Link to="/transfer">
          <span className="quick-action-icon" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none">
              <path
                d="M4 7h13M13 4l3 3-3 3M20 17H7M11 14l-3 3 3 3"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </span>

          <span className="quick-action-content">
            <strong>Transfer</strong>
            <small>Send money</small>
          </span>
        </Link>

        <Link to="/transactions">
          <span className="quick-action-icon" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none">
              <path
                d="M6 4h12v16H6zM9 8h6M9 12h6M9 16h4"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </span>

          <span className="quick-action-content">
            <strong>Transactions</strong>
            <small>View activity</small>
          </span>
        </Link>
      </div>

      <section className="dashboard-transactions">
        <div className="dashboard-section-header">
          <div>
            <h3>Recent Transactions</h3>
            <p>Latest activity on your wallet.</p>
          </div>

          <Link to="/transactions">View all</Link>
        </div>

        {transactions.length === 0 ? (
          <EmptyState message="No transactions yet." />
        ) : (
          transactions.map((transaction) => (
            <TransactionCard
              key={transaction.id}
              transaction={transaction}
            />
          ))
        )}
      </section>
    </section>
  )
}

export default DashboardPage