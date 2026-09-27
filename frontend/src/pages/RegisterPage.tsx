import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { register } from '../services/authService'
import { getApiErrorMessage } from '../utils/apiError'
import ErrorState from '../components/ErrorState'
import FormField from '../components/FormField'
import AuthContainer from '../components/AuthContainer'
import FormActions from '../components/FormActions'
import Button from '../components/Button'
import Input from '../components/Input'
import '../assets/css/RegisterPage.css'

function RegisterPage() {
  const navigate = useNavigate()

  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [nameError, setNameError] = useState('')
  const [emailError, setEmailError] = useState('')
  const [passwordError, setPasswordError] = useState('')

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()
    setNameError('')
    setEmailError('')
    setPasswordError('')

    if (loading) {
      return
    }

    if (name.trim().length < 3) {
      setNameError('Name minimal 3 karakter.')
    }

    if (!email.trim()) {
      setEmailError('Email wajib diisi.')
    }

    if (password.length < 8) {
      setPasswordError('Password minimal 8 karakter.')
    }

    if (
      name.trim().length < 3 ||
      !email.trim() ||
      password.length < 8
    ) {
      return
    }

    setError('')
    setLoading(true)

    try {
      await register({
        name,
        email,
        password,
      })

      navigate('/login')
    } catch (error) {
      setError(
        getApiErrorMessage(
          error,
          'Register gagal. Periksa data yang dimasukkan.',
        ),
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthContainer>
      <section className="register-page">
        <h2>Register</h2>

        <p className="auth-description">
          Create your Fintech Wallet account.
        </p>

        <form onSubmit={handleSubmit}>
          <FormField
            label="Name"
            htmlFor="name"
            error={nameError}
            required
          >
            <Input
              id="name"
              type="text"
              value={name}
              onChange={(event) => {
                setName(event.target.value)
                setNameError('')
              }}
              maxLength={100}
              placeholder="Enter your name"
              autoComplete="name"
              aria-invalid={Boolean(nameError)}
              aria-describedby="name-error"
              required
            />
          </FormField>

          <FormField
            label="Email"
            htmlFor="email"
            error={emailError}
            required
          >
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(event) => {
                setEmail(event.target.value)
                setEmailError('')
              }}
              maxLength={150}
              placeholder="Enter your email"
              autoComplete="email"
              aria-invalid={Boolean(emailError)}
              aria-describedby="email-error"
              required
            />
          </FormField>

          <FormField
            label="Password"
            htmlFor="password"
            error={passwordError}
            required
          >
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(event) => {
                setPassword(event.target.value)
                setPasswordError('')
              }}
              maxLength={72}
              placeholder="Enter your password"
              autoComplete="new-password"
              aria-invalid={Boolean(passwordError)}
              aria-describedby="password-error"
              required
            />
          </FormField>

          {error && <ErrorState message={error} />}

          <FormActions>
            <Button
              type="submit"
              variant="primary"
              loading={loading}
            >
              {loading ? 'Registering...' : 'Register'}
            </Button>
          </FormActions>
        </form>

        <p className="auth-link">
          Already have an account?{' '}
          <Link to="/login">Login</Link>
        </p>
      </section>
    </AuthContainer>
  )
}

export default RegisterPage