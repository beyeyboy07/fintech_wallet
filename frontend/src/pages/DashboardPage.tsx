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

function DashboardPage() {
  const [wallet, setWallet] = useState<Wallet | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const user = useAuthStore((state) => state.user)

  useEffect(() => {
    async function loadDashboard() {
      try {
        const [walletData, transactionData] = await Promise.all([
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
    <section>
      <PageHeader
        title={`Welcome, ${user?.name ?? ''}`}
        description="Manage your wallet and transactions."
      />

      <Card className="balance-card">
        <p className="balance-label">Wallet Balance</p>

        <h3 className="balance-amount">
          {wallet
            ? formatCurrency(wallet.balance)
            : formatCurrency(0)}
        </h3>
      </Card>

      <div className="quick-actions">
        <Link to="/top-up">Top Up</Link>
        <Link to="/transfer">Transfer</Link>
        <Link to="/transactions">Transactions</Link>
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