import { useEffect, useState } from 'react'
import { getTransactions } from '../services/transactionService'
import type { Transaction } from '../types/transaction'
import { getApiErrorMessage } from '../utils/apiError'
import TransactionCard from '../components/TransactionCard'
import Loading from '../components/Loading'
import EmptyState from '../components/EmptyState'
import ErrorState from '../components/ErrorState'
import PageHeader from '../components/PageHeader'
import Button from '../components/Button'
import BackButton from '../components/BackButton'


function TransactionsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [loading, setLoading] = useState(true)
  const [paginationLoading, setPaginationLoading] = useState(false)
  const [error, setError] = useState('')
  const [totalTransactions, setTotalTransactions] = useState(0)

  useEffect(() => {
    async function loadTransactions() {
      setError('')

      if (page === 1 && transactions.length === 0) {
        setLoading(true)
      } else {
        setPaginationLoading(true)
      }

      try {
        const data = await getTransactions(page, 10)

        setTransactions(data.items)
        setTotalPages(data.pagination.total_pages)
        setTotalTransactions(data.pagination.total)
      } catch (error) {
        setError(
          getApiErrorMessage(
            error,
            'Gagal mengambil riwayat transaksi.',
          ),
        )
      } finally {
        setLoading(false)
        setPaginationLoading(false)
      }
    }

    loadTransactions()
  }, [page])

  useEffect(() => {
    if (totalPages < 1) {
      return
    }

    if (page > totalPages) {
      setPage(totalPages)
    }
  }, [page, totalPages])

  if (loading) {
    return <Loading message="Loading transactions..." />
  }

  return (
    <section className="transactions-page">
      <BackButton />

      <PageHeader
        title="Transaction History"
        description="View your wallet transaction history."
      />
      {error && <ErrorState message={error} />}
      {transactions.length === 0 ? (
        <EmptyState message="No transactions found." />
      ) : (
        <div className="transactions-list">
          {transactions.map((transaction) => (
            <TransactionCard
              key={transaction.id}
              transaction={transaction}
            />
          ))}
        </div>
      )}

      {transactions.length > 0 && totalPages > 1 && (
        <div className="pagination">
          <Button
            type="button"
            variant="secondary"
            disabled={page === 1 || paginationLoading}
            onClick={() => setPage((currentPage) => currentPage - 1)}
          >
            Previous
          </Button>

          <span>
            Page {page} of {totalPages} · {totalTransactions} transactions
          </span>

          <Button
            type="button"
            variant="secondary"
            disabled={page === totalPages || paginationLoading}
            onClick={() => setPage((currentPage) => currentPage + 1)}
          >
            Next
          </Button>
        </div>
      )}
    </section>
  )
}

export default TransactionsPage