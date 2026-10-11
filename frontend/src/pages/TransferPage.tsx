
import { useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { transfer } from '../services/transferService'
import { getApiErrorMessage } from '../utils/apiError'
import ErrorState from '../components/ErrorState'
import PageHeader from '../components/PageHeader'
import FormField from '../components/FormField'
import FormActions from '../components/FormActions'
import Button from '../components/Button'
import Card from '../components/Card'
import Input from '../components/Input'
import Textarea from '../components/Textarea'
import BackButton from '../components/BackButton'

interface PendingTransfer {
  signature: string
  idempotencyKey: string
}

function createIdempotencyKey(): string {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID()
  }

  return `${Date.now()}-${Math.random().toString(16).slice(2)}`
}

function TransferPage() {
  const navigate = useNavigate()

  const pendingTransfer = useRef<PendingTransfer | null>(null)

  const [accountNumber, setAccountNumber] = useState('')
  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [accountNumberError, setAccountNumberError] = useState('')
  const [amountError, setAmountError] = useState('')

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()

    setAccountNumberError('')
    setAmountError('')

    if (loading) {
      return
    }

    const trimmedAccountNumber = accountNumber.trim()
    const trimmedAmount = amount.trim()
    const trimmedDescription = description.trim()

    let isValid = true

    if (!trimmedAccountNumber) {
      setAccountNumberError('Account number wajib diisi.')
      isValid = false
    } else if (!/^\d{16}$/.test(trimmedAccountNumber)) {
      setAccountNumberError(
        'Account number harus terdiri dari tepat 16 digit angka.',
      )
      isValid = false
    }

    const numericAmount = Number(trimmedAmount)

    if (!trimmedAmount) {
      setAmountError('Amount wajib diisi.')
      isValid = false
    } else if (
      !Number.isSafeInteger(numericAmount) ||
      numericAmount <= 0
    ) {
      setAmountError(
        'Amount harus berupa angka bulat lebih dari 0.',
      )
      isValid = false
    }

    if (!isValid) {
      return
    }

    if (trimmedDescription.length > 255) {
      setError('Description maksimal 255 karakter.')
      return
    }

    setError('')
    setLoading(true)

    const requestData = {
      receiver_account_number: trimmedAccountNumber,
      amount: numericAmount,
      description: trimmedDescription || undefined,
    }

    const signature = JSON.stringify(requestData)

    if (
      !pendingTransfer.current ||
      pendingTransfer.current.signature !== signature
    ) {
      pendingTransfer.current = {
        signature,
        idempotencyKey: createIdempotencyKey(),
      }
    }

    try {
      await transfer(
        requestData,
        pendingTransfer.current.idempotencyKey,
      )

      pendingTransfer.current = null
      navigate('/')
    } catch (error) {
      setError(
        getApiErrorMessage(
          error,
          'Transfer gagal. Periksa data yang dimasukkan.',
        ),
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="transfer-page">
      <BackButton />

      <PageHeader
        title="Transfer"
        description="Transfer balance using the recipient's account number."
      />

      <Card>
        <form onSubmit={handleSubmit}>
          <FormField
            label="Account Number"
            htmlFor="accountNumber"
            error={accountNumberError}
            required
          >
            <Input
              id="accountNumber"
              type="text"
              inputMode="numeric"
              maxLength={16}
              value={accountNumber}
              onChange={(event) => {
                const value = event.target.value

                if (/^\d*$/.test(value) && value.length <= 16) {
                  setAccountNumber(value)
                  setAccountNumberError('')
                  setError('')
                }
              }}
              placeholder="Enter 16-digit account number"
              autoComplete="off"
              aria-invalid={Boolean(accountNumberError)}
              aria-describedby="accountNumber-error"
              required
            />
          </FormField>

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
              {loading ? 'Processing...' : 'Transfer'}
            </Button>
          </FormActions>
        </form>
      </Card>
    </section>
  )
}

export default TransferPage
