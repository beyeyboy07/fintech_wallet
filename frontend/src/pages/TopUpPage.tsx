import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { topUp } from '../services/walletService'
import { getApiErrorMessage } from '../utils/apiError'
import { formatCurrency } from '../utils/currency'
import Loading from '../components/Loading'
import ErrorState from '../components/ErrorState'
import PageHeader from '../components/PageHeader'
import FormField from '../components/FormField'
import FormActions from '../components/FormActions'
import Button from '../components/Button'
import Card from '../components/Card'
import Input from '../components/Input'
import Textarea from '../components/Textarea'

const MAX_TOP_UP_AMOUNT = 100000000

function TopUpPage() {
  const navigate = useNavigate()

  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [amountError, setAmountError] = useState('')

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()

    setAmountError('')

    if (loading) {
      return
    }

    if (!amount.trim()) {
      setAmountError('Amount wajib diisi.')
      return
    }

    const numericAmount = Number(amount)

    if (!Number.isInteger(numericAmount) || numericAmount <= 0) {
      setAmountError(
        `Amount harus berupa angka bulat lebih dari 0. Contoh: ${formatCurrency(10000)}.`,
      )
      return
    }

    if (numericAmount > MAX_TOP_UP_AMOUNT) {
      setAmountError(
        `Amount maksimal ${formatCurrency(MAX_TOP_UP_AMOUNT)}.`,
      )
      return
    }

    const trimmedDescription = description.trim()

    if (trimmedDescription.length > 255) {
      setError('Description maksimal 255 karakter.')
      return
    }

    setError('')
    setLoading(true)

    try {
      await topUp({
        amount: numericAmount,
        description: trimmedDescription || undefined,
      })

      navigate('/')
    } catch (error) {
      setError(
        getApiErrorMessage(
          error,
          'Top up gagal. Periksa nominal yang dimasukkan.',
        ),
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="top-up-page">
      <PageHeader
        title="Top Up"
        description="Add balance to your wallet."
      />

      <Card>
        <form onSubmit={handleSubmit}>
          <FormField
            label="Amount"
            htmlFor="amount"
            error={amountError}
            required
          >
            <Input
              id="amount"
              type="number"
              min="1"
              step="1"
              value={amount}
              onChange={(event) => {
                setAmount(event.target.value)
                setAmountError('')
                setError('')
              }}
              placeholder="Enter amount"
              autoComplete="off"
              aria-invalid={Boolean(amountError)}
              aria-describedby="amount-error"
              required
            />
          </FormField>

          <FormField
            label="Description"
            htmlFor="description"
          >
            <Textarea
              id="description"
              value={description}
              onChange={(event) => {
                setDescription(event.target.value)
                setError('')
              }}
              maxLength={255}
              placeholder="Enter description (optional)"
              autoComplete="off"
              rows={4}
            />
          </FormField>

          {error && <ErrorState message={error} />}

          <FormActions>
            <Button
              type="submit"
              variant="primary"
              loading={loading}
            >
              {loading ? 'Processing...' : 'Top Up'}
            </Button>
          </FormActions>

          {loading && (
            <Loading message="Processing top up..." />
          )}
        </form>
      </Card>
    </section>
  )
}

export default TopUpPage
