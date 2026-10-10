import { useState } from 'react'
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


function TransferPage() {
  const navigate = useNavigate()

  const [receiverId, setReceiverId] = useState('')
  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [receiverIdError, setReceiverIdError] = useState('')
  const [amountError, setAmountError] = useState('')

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()

    setReceiverIdError('')
    setAmountError('')

    if (loading) {
      return
    }

    const numericReceiverId = Number(receiverId)
    const numericAmount = Number(amount)

    if (!receiverId.trim()) {
      setReceiverIdError('Receiver ID wajib diisi.')
    }

    if (!amount.trim()) {
      setAmountError('Amount wajib diisi.')
    }

    if (
      !receiverId.trim() ||
      !amount.trim()
    ) {
      return
    }

    if (
      !Number.isInteger(numericReceiverId) ||
      numericReceiverId <= 0
    ) {
      setReceiverIdError(
        'Receiver ID harus berupa angka bulat lebih dari 0.',
      )
    }

    if (
      !Number.isInteger(numericAmount) ||
      numericAmount <= 0
    ) {
      setAmountError(
        'Amount harus berupa angka bulat lebih dari 0.',
      )
    }

    if (
      !Number.isInteger(numericReceiverId) ||
      numericReceiverId <= 0 ||
      !Number.isInteger(numericAmount) ||
      numericAmount <= 0
    ) {
      return
    }

    setError('')

    const trimmedDescription = description.trim()

    if (trimmedDescription.length > 255) {
      setError('Description maksimal 255 karakter.')
      return
    }

    setLoading(true)

    try {
      await transfer(
        {
          receiver_id: numericReceiverId,
          amount: numericAmount,
          description: trimmedDescription || undefined,
        },
        `${Date.now()}-${Math.random().toString(16).slice(2)}-${Math.random().toString(16).slice(2)}`,
      )

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
        description="Send balance to another wallet."
      />

      <Card>
        <form onSubmit={handleSubmit}>
          <FormField
            label="Receiver ID"
            htmlFor="receiverId"
            error={receiverIdError}
            required
          >
            <Input
              id="receiverId"
              type="number"
              min="1"
              step="1"
              value={receiverId}
              onChange={(event) => {
                setReceiverId(event.target.value)
                setReceiverIdError('')
                setError('')
              }}
              placeholder="Enter receiver ID"
              autoComplete="off"
              aria-invalid={Boolean(receiverIdError)}
              aria-describedby="receiverId-error"
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
