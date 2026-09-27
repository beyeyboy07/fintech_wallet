import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { login, getMe } from '../services/authService'
import { setTokens } from '../utils/tokenStorage'
import { useAuthStore } from '../stores/authStore'
import { getApiErrorMessage } from '../utils/apiError'
import ErrorState from '../components/ErrorState'
import FormField from '../components/FormField'
import AuthContainer from '../components/AuthContainer'
import FormActions from '../components/FormActions'
import Button from '../components/Button'
import Input from '../components/Input'
import '../assets/css/LoginPage.css'

function LoginPage() {
  const navigate = useNavigate()
  const setUser = useAuthStore((state) => state.setUser)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [emailError, setEmailError] = useState('')
  const [passwordError, setPasswordError] = useState('')

  async function handleSubmit(
    event: React.FormEvent<HTMLFormElement>,
  ) {
    event.preventDefault()
    setEmailError('')
    setPasswordError('')

    if (loading) {
      return
    }

    if (!email.trim()) {
      setEmailError('Email wajib diisi.')
    }

    if (!password) {
      setPasswordError('Password wajib diisi.')
    }

    if (!email.trim() || !password) {
      return
    }

    setError('')
    setLoading(true)

    try {
      const tokenResponse = await login({
        email,
        password,
      })

      setTokens(
        tokenResponse.data.access_token,
        tokenResponse.data.refresh_token,
      )

      const user = await getMe()

      setUser(user)
      navigate('/')
    } catch (error) {
      setError(
        getApiErrorMessage(
          error,
          'Login gagal. Periksa email dan password.',
        ),
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthContainer>
      <section className="login-page">
        <h2>Login</h2>

        <p className="auth-description">
          Sign in to manage your wallet.
        </p>

        <form onSubmit={handleSubmit}>
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
              autoComplete="current-password"
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
              {loading ? 'Logging in...' : 'Login'}
            </Button>
          </FormActions>
        </form>

        <p className="auth-link">
          Don't have an account?{' '}
          <Link to="/register">Register</Link>
        </p>
      </section>
    </AuthContainer>
  )
}

export default LoginPage